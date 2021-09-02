package bgstripe

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/event"
)

func ForceSyncStripeEvents() {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	params := &stripe.EventListParams{
		DeliverySuccess: ptr.Bool(false),
	}
	events := event.List(params)
	for events.Next() {
		event := events.Event()
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			log.Println(fmt.Sprintf("Handling event of type %s", event.Type))
			return HandleStripeEvent(tx, *event)
		}); err != nil {
			log.Println(fmt.Sprintf("Error processing event: %s", err.Error()))
		}
	}
}

func HandleStripeEvent(tx *sqlx.Tx, event stripe.Event) error {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	switch event.Type {
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			return err
		}
		return handlePaymentMethodAttached(tx, paymentMethod)
	case "payment_method.detached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			return err
		}
		return handlePaymentMethodDetached(tx, paymentMethod)
	case "invoice.paid",
		"invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return err
		}
		return handleSuccessfulInvoice(tx, invoice)
	case "invoice.payment_failed",
		"invoice.payment_action_required":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return err
		}
		return handleFailedInvoice(tx, invoice)
	case "customer.subscription.trial_will_end":
		// Alert user that trial subscription will be ending
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return err
		}
		return handleTrialWillEnd(tx, subscription)
	case "customer.subscription.updated":
		// Handle state change
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return err
		}
		return handleSubscriptionUpdated(tx, subscription)
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return err
		}
		return handleSubscriptionDeleted(tx, subscription)
	}
	return nil
}

func handlePaymentMethodAttached(tx *sqlx.Tx, paymentMethod stripe.PaymentMethod) error {
	if paymentMethod.Customer == nil {
		return fmt.Errorf("No customer for payment method")
	}
	customerID := CustomerID(paymentMethod.Customer.ID)
	userID, err := GetUserIDForStripeCustomerID(tx, customerID)
	if err != nil {
		return err
	}
	return InsertPaymentMethod(tx, *userID, &paymentMethod)
}

func handlePaymentMethodDetached(tx *sqlx.Tx, paymentMethod stripe.PaymentMethod) error {
	return RemovePaymentMethod(tx, PaymentMethodID(paymentMethod.ID))
}

func handleSuccessfulInvoice(tx *sqlx.Tx, invoice stripe.Invoice) error {
	userID, subscription, err := ReconcileInvoice(tx, invoice)
	if err != nil {
		return err
	}
	return useraccounts.UpdateSubscriptionExpirationTime(tx, *userID, subscription.CurrentPeriodEnd)
}

func handleFailedInvoice(tx *sqlx.Tx, invoice stripe.Invoice) error {
	userID, subscription, err := ReconcileInvoice(tx, invoice)
	if err != nil {
		return err
	}
	holdUntilTime := time.Now().Add(5 * time.Minute)
	enqueuedNotificationRequest, err := useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePaymentError, holdUntilTime)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", subscription.StripeSubscriptionID, enqueuedNotificationRequest))
	return nil
}

func handleTrialWillEnd(tx *sqlx.Tx, subscription stripe.Subscription) error {
	userID, err := LookupBabblegraphUserIDForStripeSubscriptionID(tx, subscription.ID)
	if err != nil {
		return err
	}
	holdUntilTime := time.Now().Add(30 * time.Minute)
	enqueuedNotificationRequest, err := useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypeTrialEndingSoon, holdUntilTime)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", subscription.ID, enqueuedNotificationRequest))
	return nil
}

func handleSubscriptionUpdated(tx *sqlx.Tx, subscription stripe.Subscription) error {
	userID, err := LookupBabblegraphUserIDForStripeSubscriptionID(tx, subscription.ID)
	if err != nil {
		return err
	}
	reconciliationAction, err := ReconcileSubscriptionUpdate(tx, subscription)
	if err != nil {
		return err
	}
	if reconciliationAction != nil {
		switch *reconciliationAction {
		case SubscriptionReconciliationActionCancellation:
			if err := useraccounts.ExpireSubscriptionForUser(tx, *userID); err != nil {
				return err
			}
			holdUntilTime := time.Now().Add(5 * time.Minute)
			enqueuedNotificationRequest, err := useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, holdUntilTime)
			if err != nil {
				return err
			}
			log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", subscription.ID, enqueuedNotificationRequest))
		case SubscriptionReconciliationActionFirstPaymentSuccessful:
			// Send notification that subscription has started
		default:
			// No-op
		}
	}
	return nil
}

func handleSubscriptionDeleted(tx *sqlx.Tx, subscription stripe.Subscription) error {
	userID, err := LookupBabblegraphUserIDForStripeSubscriptionID(tx, subscription.ID)
	if err != nil {
		return err
	}
	if err := useraccounts.ExpireSubscriptionForUser(tx, *userID); err != nil {
		return err
	}
	holdUntilTime := time.Now().Add(5 * time.Minute)
	_, err = useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, holdUntilTime)
	return err
}
