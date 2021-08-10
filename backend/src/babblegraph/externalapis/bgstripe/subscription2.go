package bgstripe

import (
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

	insertStripeSubscriptionForUserQuery = "INSERT INTO bgstripe_subscription (babblegraph_user_id, stripe_subscription_id, payment_state, stripe_product_id) VALUES ($1, $2, $3, $4)"

	lookupActiveSubscriptionForUserQuery = "SELECT * FROM bgstripe_subscription WHERE babblegraph_user_id = $1 AND payment_state != $2"

	getTrialEligibilityLengthForUserIDQuery = "SELECT * FROM bgstripe_subscription WHERE babblegraph_user_id = $1 ORDER BY created_at ASC LIMIT 1"
)

type Subscription struct {
	StripeSubscriptionID      SubscriptionID        `json:"stripe_subscription_id"`
	PaymentState              PaymentState          `json:"payment_state"`
	CurrentPeriodEnd          time.Time             `json:"current_period_end"`
	CancelAtPeriodEnd         bool                  `json:"cancel_at_period_end"`
	PaymentIntentClientSecret *string               `json:"payment_intent_client_secret,omitempty"`
	SubscriptionType          SubscriptionType      `json:"subscription_type"`
	TrialInfo                 SubscriptionTrialInfo `json:"trial_info"`
}

type SubscriptionType string

const (
	SubscriptionTypeYearly  SubscriptionType = "yearly"
	SubscriptionTypeMonthly SubscriptionType = "monthly"
)

func (s SubscriptionType) Ptr() *SubscriptionType {
	return &s
}

type SubscriptionTrialInfo struct {
	IsCurrentlyTrialing  bool  `json:"is_currently_trialing"`
	TrialEligibilityDays int64 `json:"trial_eligibility_days"`
}

func CreateSubscriptionForUser(tx *sqlx.Tx, userID users.UserID, subscriptionType SubscriptionType) (*Subscription, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	existingSubscription, err := LookupActiveSubscriptionForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	if existingSubscription != nil {
		return nil, fmt.Errorf("User already has an active subscription")
	}
	stripeCustomer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	trialEligibilityDays, err := getTrialEligibilityLengthForUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	stripeProductID, err := getProductIDForSubscriptionType(subscriptionType)
	if err != nil {
		return nil, err
	}
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(string(stripeCustomer.StripeCustomerID)),
		Items: []*stripe.SubscriptionItemsParams{
			&stripe.SubscriptionItemsParams{
				Price: ptr.String(stripeProductID.Str()),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	paymentState := PaymentStateCreatedUnpaid
	if *trialEligibilityDays > 0 {
		paymentState = PaymentStateTrialNoPaymentMethod
		subscriptionParams.TrialPeriodDays = trialEligibilityDays
	}
	stripeSubscription, err := sub.New(subscriptionParams)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(insertStripeSubscriptionForUserQuery, userID, stripeSubscription.ID, paymentState, stripeProductID); err != nil {
		log.Println("Attempting to rollback stripe subscription")
		if _, sErr := sub.Cancel(stripeSubscription.ID, &stripe.SubscriptionCancelParams{}); sErr != nil {
			formattedSErr := fmt.Errorf("Error rolling back stripe subscription %s for user %s because of %s. Original error: %s", stripeSubscription.ID, userID, sErr.Error(), err.Error())
			log.Println(formattedSErr.Error())
			sentry.CaptureException(formattedSErr)
		}
		return nil, err
	}
	dbSubscription, err := lookupActiveDBSubscriptionForUser(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case dbSubscription == nil:
		return nil, fmt.Errorf("Could not retrieve subscription from database")
	}
	return mergeSubscriptionObjects(tx, *dbSubscription, stripeSubscription)
}

func LookupActiveSubscriptionForUser(tx *sqlx.Tx, userID users.UserID) (*Subscription, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	var matches []dbStripeSubscription
	if err := tx.Select(&matches, lookupActiveSubscriptionForUserQuery, userID, PaymentStateTerminated); err != nil {
		return nil, err
	}
	dbSubscription, err := lookupActiveDBSubscriptionForUser(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case dbSubscription == nil:
		return nil, nil
	}
	subscriptionParams := &stripe.SubscriptionParams{}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	stripeSubscription, err := sub.Get(string(dbSubscription.StripeSubscriptionID), subscriptionParams)
	if err != nil {
		return nil, err
	}
	return mergeSubscriptionObjects(tx, *dbSubscription, stripeSubscription)
}

func mergeSubscriptionObjects(tx *sqlx.Tx, bgSub dbStripeSubscription, stripeSub *stripe.Subscription) (*Subscription, error) {
	trialEligibilityDays, err := getTrialEligibilityLengthForUserID(tx, bgSub.BabblegraphUserID)
	if err != nil {
		return nil, err
	}
	var paymentIntentClientSecret *string
	if bgSub.PaymentState == PaymentStateCreatedUnpaid {
		if stripeSub.LatestInvoice == nil || stripeSub.LatestInvoice.PaymentIntent == nil {
			return nil, fmt.Errorf("Expected latest invoice and payment intent to be nonnil")
		}
		paymentIntentClientSecret = ptr.String(stripeSub.LatestInvoice.PaymentIntent.ClientSecret)
	}
	subscriptionType, err := getSubscriptionTypeForProductID(bgSub.StripeProductID)
	if err != nil {
		return nil, err
	}
	return &Subscription{
		StripeSubscriptionID:      bgSub.StripeSubscriptionID,
		PaymentState:              bgSub.PaymentState,
		CurrentPeriodEnd:          time.Unix(stripeSub.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:         stripeSub.CancelAtPeriodEnd,
		PaymentIntentClientSecret: paymentIntentClientSecret,
		SubscriptionType:          *subscriptionType,
		TrialInfo: SubscriptionTrialInfo{
			IsCurrentlyTrialing:  bgSub.PaymentState == PaymentStateTrialPaymentMethodAdded || bgSub.PaymentState == PaymentStateTrialNoPaymentMethod,
			TrialEligibilityDays: *trialEligibilityDays,
		},
	}, nil
}

func lookupActiveDBSubscriptionForUser(tx *sqlx.Tx, userID users.UserID) (*dbStripeSubscription, error) {
	var matches []dbStripeSubscription
	if err := tx.Select(&matches, lookupActiveSubscriptionForUserQuery, userID, PaymentStateTerminated); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		m := matches[0]
		return &m, nil
	default:
		return nil, fmt.Errorf("Expected at most one subscription, but got %d", len(matches))
	}
	return nil, fmt.Errorf("Unreachable")
}

func getTrialEligibilityLengthForUserID(tx *sqlx.Tx, userID users.UserID) (*int64, error) {
	var matches []dbStripeSubscription
	if err := tx.Select(&matches, getTrialEligibilityLengthForUserIDQuery, userID); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return ptr.Int64(defaultSubscriptionTrialLength), nil
	case len(matches) == 1:
		now := time.Now()
		daysSinceFirstSubscription := int64(math.Abs(now.Sub(matches[0].CreatedAt).Hours() / 24.0))
		if daysSinceFirstSubscription >= defaultSubscriptionTrialLength {
			return ptr.Int64(0), nil
		}
		return ptr.Int64(defaultSubscriptionTrialLength - daysSinceFirstSubscription), nil
	default:
		return nil, fmt.Errorf("Expected at most one subscription, but got %d", len(matches))
	}
}

func getSubscriptionTypeForProductID(productID StripeProductID) (*SubscriptionType, error) {
	switch productID {
	case StripeProductIDYearlySubscriptionTest,
		StripeProductIDYearlySubscriptionProd:
		return SubscriptionTypeYearly.Ptr(), nil
	case StripeProductIDMonthlySubscriptionTest,
		StripeProductIDMonthlySubscriptionProd:
		return SubscriptionTypeMonthly.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unrecognized product ID %s", productID)
	}
}

func getProductIDForSubscriptionType(subscriptionType SubscriptionType) (*StripeProductID, error) {
	currentEnv := env.MustEnvironmentName()
	switch currentEnv {
	case env.EnvironmentProd:
		if subscriptionType == SubscriptionTypeYearly {
			return StripeProductIDYearlySubscriptionProd.Ptr(), nil
		}
		return StripeProductIDMonthlySubscriptionProd.Ptr(), nil
	case env.EnvironmentStage,
		env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		if subscriptionType == SubscriptionTypeYearly {
			return StripeProductIDYearlySubscriptionTest.Ptr(), nil
		}
		return StripeProductIDMonthlySubscriptionTest.Ptr(), nil
	default:
		return nil, fmt.Errorf("unsupported environment: %s", currentEnv)
	}
}
