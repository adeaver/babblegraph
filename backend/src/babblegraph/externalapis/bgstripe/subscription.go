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

	getStripeSubscriptionsForUserQuery        = "SELECT * FROM bgstripe_subscription WHERE babblegraph_user_id = $1"
	insertStripeSubscriptionForUserQuery      = "INSERT INTO bgstripe_subscription (babblegraph_user_id, stripe_subscription_id, payment_state, stripe_client_secret, stripe_product_id) VALUES ($1, $2, $3, $4, $5)"
	updateStripeSubscriptionPaymentStateQuery = "UPDATE bgstripe_subscription SET payment_state = $1 WHERE babblegraph_user_id = $2 AND stripe_subscription_id = $3"
	updateStripeSubscriptionProductID         = "UPDATE bgstripe_subscription SET stripe_product_id = $1 WHERE babblegraph_user_id = $2 AND stripe_subscription_id = $3"
)

type StripeCustomerSubscriptionOutput struct {
	ClientSecret         string
	SubscriptionID       SubscriptionID
	PaymentState         PaymentState
	IsYearlySubscription bool
}

func CreateUnpaidStripeCustomerSubscriptionForUser(tx *sqlx.Tx, userID users.UserID, isYearlySubscription bool) (*StripeCustomerSubscriptionOutput, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	stripeCustomer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	stripeProductID := getPriceIDForEnvironmentAndPaymentType(isYearlySubscription)
	// Trials are automatically considered active
	newPaymentState := PaymentStateTrialNoPaymentMethod
	trialPeriodDays, err := getNumberOfDaysOfTrial(tx, userID)
	if err != nil {
		return nil, err
	}
	if *trialPeriodDays <= 0 {
		newPaymentState = PaymentStateCreatedUnpaid
		trialPeriodDays = ptr.Int64(0)
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
	case PaymentStateTrialNoPaymentMethod:
		// Create a trial
		if err := useraccounts.AddSubscriptionLevelForUser(tx, userID, useraccounts.SubscriptionLevelPremium); err != nil {
			return nil, err
		}
		subscriptionParams.TrialPeriodDays = stripe.Int64(*trialPeriodDays)
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
	case PaymentStateTrialNoPaymentMethod:
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
		SubscriptionID:       SubscriptionID(stripeSubscription.ID),
		ClientSecret:         *clientSecret,
		PaymentState:         newPaymentState,
		IsYearlySubscription: isYearlySubscription,
	}, nil
}

func LookupNonterminatedStripeSubscriptionForUser(tx *sqlx.Tx, userID users.UserID) (*StripeCustomerSubscriptionOutput, bool, error) {
	stripeSubscriptionsForUser, err := lookupStripeSubscriptionsForUser(tx, userID)
	if err != nil {
		return nil, false, err
	}
	numberOfTrialDays, err := getNumberOfDaysOfTrial(tx, userID)
	if err != nil {
		return nil, false, err
	}
	isEligibleForTrial := *numberOfTrialDays > 0
	for _, subscription := range stripeSubscriptionsForUser {
		if subscription.PaymentState != PaymentStateTerminated {
			isYearlySubscription, err := isStripeProductIDYearly(subscription.StripeProductID)
			if err != nil {
				return nil, false, err
			}
			return &StripeCustomerSubscriptionOutput{
				SubscriptionID:       subscription.StripeSubscriptionID,
				ClientSecret:         subscription.StripeClientSecret,
				PaymentState:         subscription.PaymentState,
				IsYearlySubscription: *isYearlySubscription,
			}, isEligibleForTrial, nil
		}
	}
	return nil, isEligibleForTrial, nil
}

func CancelStripeSubscription(tx *sqlx.Tx, userID users.UserID, stripeSubscriptionID SubscriptionID) (bool, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	res, err := tx.Exec(updateStripeSubscriptionPaymentStateQuery, PaymentStateTerminated, userID, stripeSubscriptionID)
	if err != nil {
		return false, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if numRows <= 0 {
		return false, nil
	}
	if err := useraccounts.ExpireSubscriptionForUser(tx, userID); err != nil {
		return false, err
	}
	if _, err := sub.Cancel(string(stripeSubscriptionID), &stripe.SubscriptionCancelParams{}); err != nil {
		return false, err
	}
	return true, nil
}

func UpdateStripeSubscriptionChargeFrequency(tx *sqlx.Tx, userID users.UserID, stripeSubscriptionID SubscriptionID, isYearlySubscription bool) (bool, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	stripeProductID := getPriceIDForEnvironmentAndPaymentType(isYearlySubscription)
	res, err := tx.Exec(updateStripeSubscriptionProductID, stripeProductID, userID, stripeSubscriptionID)
	if err != nil {
		return false, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if numRows <= 0 {
		return false, nil
	}
	subscription, err := sub.Get(string(stripeSubscriptionID), nil)
	if err != nil {
		return false, nil
	}
	subscriptionParams := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(subscription.Items.Data[0].ID),
				Price: stripe.String(stripeProductID.Str()),
			},
		},
	}
	if _, err := sub.Update(string(stripeSubscriptionID), subscriptionParams); err != nil {
		return false, err
	}
	return true, nil
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

func isStripeProductIDYearly(stripeProductID StripeProductID) (*bool, error) {
	switch stripeProductID {
	case StripeProductIDYearlySubscriptionTest,
		StripeProductIDYearlySubscriptionProd:
		return ptr.Bool(true), nil
	case StripeProductIDMonthlySubscriptionTest,
		StripeProductIDMonthlySubscriptionProd:
		return ptr.Bool(false), nil
	default:
		return nil, fmt.Errorf("Unrecognized product ID %s", stripeProductID)
	}
}

func getNumberOfDaysOfTrial(tx *sqlx.Tx, userID users.UserID) (*int64, error) {
	stripeSubscriptions, err := lookupStripeSubscriptionsForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	var daysSinceOldestTrialPeriod int64 = 0
	now := time.Now()
	for _, subscription := range stripeSubscriptions {
		daysSinceTrialForSubscription := int64(math.Abs(now.Sub(subscription.CreatedAt).Hours() / 24.0))
		if daysSinceTrialForSubscription > daysSinceOldestTrialPeriod {
			daysSinceOldestTrialPeriod = daysSinceTrialForSubscription
		}
	}
	return ptr.Int64(defaultSubscriptionTrialLength - daysSinceOldestTrialPeriod), nil
}
