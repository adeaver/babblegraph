package billing

import (
	"babblegraph/util/ptr"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
)

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
