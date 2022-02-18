package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/setupintent"
)

// This file is for any stripe specific methods

func GetSetupIntentClientSecretForUser(tx *sqlx.Tx, userID users.UserID) (*string, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		return nil, fmt.Errorf("Expected billing information for user %s, but got none", userID)
	default:
		externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
		if err != nil {
			return nil, err
		}
		if externalID.IDType != externalIDTypeStripe {
			return nil, fmt.Errorf("User %s is not a stripe user, has type %s", userID, externalID.IDType)
		}
		params := &stripe.SetupIntentParams{
			Customer: ptr.String(externalID.ExternalID),
			PaymentMethodTypes: []*string{
				stripe.String("card"),
			},
			Usage: ptr.String("off_session"),
		}
		si, err := setupintent.New(params)
		if err != nil {
			return nil, err
		}
		return ptr.String(si.ClientSecret), nil
	}
}

func convertStripeSubscriptionToPremiumNewsletterSubscription(tx *sqlx.Tx, stripeSubscription *stripe.Subscription, dbNewsletterSubscription *dbPremiumNewsletterSubscription) (*PremiumNewsletterSubscription, error) {
	var paymentIntentID *string
	if stripeSubscription.LatestInvoice != nil && stripeSubscription.LatestInvoice.PaymentIntent != nil {
		paymentIntentID = ptr.String(stripeSubscription.LatestInvoice.PaymentIntent.ClientSecret)
	}
	premiumNewsletterSubscription := PremiumNewsletterSubscription{
		StripePaymentIntentID: paymentIntentID,
		CurrentPeriodEnd:      time.Unix(stripeSubscription.CurrentPeriodEnd, 0),
	}
	if dbNewsletterSubscription != nil {
		billingInformation, err := getBillingInformation(tx, dbNewsletterSubscription.BillingInformationID)
		if err != nil {
			return nil, err
		}
		premiumNewsletterSubscription.userID = billingInformation.UserID
	}
	switch stripeSubscription.Status {
	case stripe.SubscriptionStatusTrialing:
		premiumNewsletterSubscription.PaymentState = PaymentStateTrialNoPaymentMethod
		if stripeSubscription.DefaultPaymentMethod != nil {
			premiumNewsletterSubscription.PaymentState = PaymentStateTrialPaymentMethodAdded
		}
	case stripe.SubscriptionStatusIncomplete:
		premiumNewsletterSubscription.PaymentState = PaymentStateCreatedUnpaid
	case stripe.SubscriptionStatusActive:
		premiumNewsletterSubscription.PaymentState = PaymentStateActive
	case stripe.SubscriptionStatusPastDue,
		stripe.SubscriptionStatusUnpaid:
		premiumNewsletterSubscription.PaymentState = PaymentStateErrored
	case stripe.SubscriptionStatusIncompleteExpired,
		stripe.SubscriptionStatusCanceled:
		premiumNewsletterSubscription.PaymentState = PaymentStateTerminated
	case stripe.SubscriptionStatusAll:
		return nil, fmt.Errorf("Unsupported payment status: all")
	default:
		return nil, fmt.Errorf("Unsupported payment status: %s", stripeSubscription.Status)
	}
	return &premiumNewsletterSubscription, nil
}
