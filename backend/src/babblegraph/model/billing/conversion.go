package billing

import (
	"fmt"

	"github.com/stripe/stripe-go/v72"
)

func convertStripeSubscriptionToPremiumNewsletterSubscription(stripeSubscription *stripe.Subscription) (*PremiumNewsletterSubscription, error) {
	premiumNewsletterSubscription := PremiumNewsletterSubscription{}
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
