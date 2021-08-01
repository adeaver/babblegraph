package bgstripe

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

const (
	defaultSubscriptionTrialLength = 14

	getStripeSubscriptionsForUserQuery   = "SELECT * FROM bgstripe_subscription WHERE babblegraph_user_id = $1"
	insertStripeSubscriptionForUserQuery = "INSERT INTO bgstripe_subscription (babblegraph_user_id, stripe_subscription_id, payment_state, stripe_client_secret, stripe_product_id) VALUES ($1, $2, $3, $4, $5)"
)

type StripeCustomerSubscriptionOutput struct {
	ClientSecret   string
	SubscriptionID SubscriptionID
	PaymentState   PaymentState
}

func GetOrCreateUnpaidStripeCustomerSubscriptionForUser(tx *sqlx.Tx, userID users.UserID, isYearlySubscription bool) (*StripeCustomerSubscriptionOutput, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	stripeCustomer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	stripeSubscriptions, err := lookupStripeSubscriptionsForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	var daysSinceOldestTrialPeriod int64 = 0
	now := time.Now()
	stripeProductID := getPriceIDForEnvironmentAndPaymentType(isYearlySubscription)
	for _, subscription := range stripeSubscriptions {
		// Heuristics here => if there's already an unpaid subscription with the
		// same product ID, then we should use that. Otherwise, we need to create a new
		// one. We also need to figure out the amount of trial days that a user has left.
		if subscription.PaymentState == PaymentStateActive {
			return nil, fmt.Errorf("User %s already has an active subscription", userID)
		}
		daysSinceTrialForSubscription := int64(math.Abs(now.Sub(subscription.CreatedAt).Hours() / 24.0))
		if daysSinceTrialForSubscription > daysSinceOldestTrialPeriod {
			daysSinceOldestTrialPeriod = daysSinceTrialForSubscription
		}
		if subscription.StripeProductID == stripeProductID {
			switch subscription.PaymentState {
			case PaymentStateCreatedUnpaid,
				PaymentStateActiveNoPaymentMethod:
				// These are unpaid subscriptions with the same product ID. We should return this.
				return &StripeCustomerSubscriptionOutput{
					SubscriptionID: subscription.StripeSubscriptionID,
					ClientSecret:   subscription.StripeClientSecret,
					PaymentState:   subscription.PaymentState,
				}, nil
			case PaymentStateActive:
				return nil, fmt.Errorf("Invalid state for creating a new subscription")
			case PaymentStateTerminated,
				PaymentStateErrored:
				// no-op
			default:
				return nil, fmt.Errorf("Unrecognized payment state: %d", subscription.PaymentState)
			}
		}
	}
	// Trials are automatically considered active
	newPaymentState := PaymentStateActiveNoPaymentMethod
	trialPeriodDays := defaultSubscriptionTrialLength - daysSinceOldestTrialPeriod
	if trialPeriodDays <= 0 {
		newPaymentState = PaymentStateCreatedUnpaid
		trialPeriodDays = 0
	}
	subscriptionPriceLineItem := stripe.SubscriptionItemsParams{
		Price: stripe.String(stripeProductID.Str()),
	}
	subscriptionParams := &stripe.SubscriptionParams{
		Customer:        stripe.String(string(stripeCustomer.StripeCustomerID)),
		Items:           []*stripe.SubscriptionItemsParams{&subscriptionPriceLineItem},
		PaymentBehavior: stripe.String("default_incomplete"),
	}
	switch newPaymentState {
	case PaymentStateActiveNoPaymentMethod:
		// Create a trial
		if err := useraccounts.AddSubscriptionLevelForUser(tx, userID, useraccounts.SubscriptionLevelPremium); err != nil {
			return nil, err
		}
		subscriptionParams.TrialPeriodDays = stripe.Int64(trialPeriodDays)
		subscriptionParams.AddExpand("pending_setup_intent")
	case PaymentStateCreatedUnpaid:
		// Create a subscription without a trial
		subscriptionParams.AddExpand("latest_invoice.payment_intent")
	default:
		return nil, fmt.Errorf("Invalid payment state for new subscription: %d", newPaymentState)
	}
	stripeSubscription, err := sub.New(subscriptionParams)
	if err != nil {
		return nil, err
	}
	var clientSecret *string
	switch newPaymentState {
	case PaymentStateActiveNoPaymentMethod:
		if stripeSubscription.PendingSetupIntent == nil {
			return nil, fmt.Errorf("Expected pending setup to be nonnil")
		}
		clientSecret = ptr.String(stripeSubscription.PendingSetupIntent.ClientSecret)
	case PaymentStateCreatedUnpaid:
		if stripeSubscription.LatestInvoice == nil || stripeSubscription.LatestInvoice.PaymentIntent == nil {
			return nil, fmt.Errorf("Expected latest invoice and payment intent to be nonnil")
		}
		clientSecret = ptr.String(stripeSubscription.LatestInvoice.PaymentIntent.ClientSecret)
	default:
		return nil, fmt.Errorf("Invalid payment state: %d", newPaymentState)
	}
	if _, err := tx.Exec(insertStripeSubscriptionForUserQuery, userID, stripeSubscription.ID, newPaymentState, *clientSecret, stripeProductID); err != nil {
		log.Println("Attempting to rollback stripe subscription")
		if _, sErr := sub.Cancel(stripeSubscription.ID, &stripe.SubscriptionCancelParams{}); sErr != nil {
			formattedSErr := fmt.Errorf("Error rolling back stripe subscription %s for user %s because of %s. Original error: %s", stripeSubscription.ID, userID, sErr.Error(), err.Error())
			log.Println(formattedSErr.Error())
			sentry.CaptureException(formattedSErr)
		}
		return nil, err
	}
	return &StripeCustomerSubscriptionOutput{
		SubscriptionID: SubscriptionID(stripeSubscription.ID),
		ClientSecret:   *clientSecret,
		PaymentState:   newPaymentState,
	}, nil
}

func lookupStripeSubscriptionsForUser(tx *sqlx.Tx, userID users.UserID) ([]dbStripeSubscription, error) {
	var matches []dbStripeSubscription
	if err := tx.Select(&matches, getStripeSubscriptionsForUserQuery, userID); err != nil {
		return nil, err
	}
	return matches, nil
}

func getPriceIDForEnvironmentAndPaymentType(isYearlySubscription bool) StripeProductID {
	currentEnv := env.MustEnvironmentName()
	switch currentEnv {
	case env.EnvironmentProd:
		if isYearlySubscription {
			return StripeProductIDYearlySubscriptionProd
		}
		return StripeProductIDMonthlySubscriptionProd
	case env.EnvironmentStage,
		env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		if isYearlySubscription {
			return StripeProductIDYearlySubscriptionTest
		}
		return StripeProductIDMonthlySubscriptionTest
	default:
		panic(fmt.Sprintf("unsupported environment: %s", currentEnv))
	}
}
