package bgstripe

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

const (
	defaultSubscriptionTrialLength = 14

	insertStripeSubscriptionForUserQuery = "INSERT INTO bgstripe_subscription (babblegraph_user_id, stripe_subscription_id, is_active) VALUES ($1, $2, FALSE)"
)

func CreateStripeCustomerSubscriptionForUser(tx *sqlx.Tx, userID users.UserID, isYearlySubscription bool) (*SubscriptionID, *string, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	stripeCustomer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return nil, nil, err
	}
	subscriptionPriceLineItem := stripe.SubscriptionItemsParams{
		Price: stripe.String(getPriceIDForEnvironmentAndPaymentType(isYearlySubscription)),
	}
	// TODO: query the database to see if they're eligible for a free trial
	subscriptionParams := &stripe.SubscriptionParams{
		Customer:        stripe.String(string(stripeCustomer.StripeCustomerID)),
		Items:           []*stripe.SubscriptionItemsParams{&subscriptionPriceLineItem},
		TrialPeriodDays: stripe.Int64(defaultSubscriptionTrialLength),
		PaymentBehavior: stripe.String("default_incomplete"),
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	subscriptionParams.AddExpand("pending_setup_intent")
	stripeSubscription, err := sub.New(subscriptionParams)
	if err != nil {
		return nil, nil, err
	}
	if _, err := tx.Exec(insertStripeSubscriptionForUserQuery, userID, stripeSubscription.ID); err != nil {
		log.Println("Attempting to rollback stripe subscription")
		if _, sErr := sub.Cancel(stripeSubscription.ID, &stripe.SubscriptionCancelParams{}); sErr != nil {
			formattedSErr := fmt.Errorf("Error rolling back stripe subscription %s for user %s because of %s. Original error: %s", stripeSubscription.ID, userID, sErr.Error(), err.Error())
			log.Println(formattedSErr.Error())
			sentry.CaptureException(formattedSErr)
		}
		return nil, nil, err
	}
	log.Println(fmt.Sprintf("%+v", stripeSubscription))
	log.Println(fmt.Sprintf("%+v", stripeSubscription.PendingSetupIntent))
	asSubscriptionID := SubscriptionID(stripeSubscription.ID)
	return &asSubscriptionID, ptr.String(stripeSubscription.PendingSetupIntent.ClientSecret), nil
}

func getPriceIDForEnvironmentAndPaymentType(isYearlySubscription bool) string {
	currentEnv := env.MustEnvironmentName()
	switch currentEnv {
	case env.EnvironmentProd:
		if isYearlySubscription {
			return "price_1JIMqNJscBSiX47SxOGRUX1p"
		}
		return "price_1JIMqNJscBSiX47SnYtkOVv6"
	case env.EnvironmentStage,
		env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		if isYearlySubscription {
			return "price_1JIMr1JscBSiX47SEEUzRf0e"
		}
		return "price_1JIMr1JscBSiX47SReF6SdJj"
	default:
		panic(fmt.Sprintf("unsupported environment: %s", currentEnv))
	}
}
