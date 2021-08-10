package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/useraccounts"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

func handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	webhookSecret := env.MustEnvironmentVariable("STRIPE_WEBHOOK_SECRET")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleWebhookError(w, "reading body", err)
		return
	}
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		handleWebhookError(w, "constructing event", err)
		return
	}
	switch event.Type {
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			handleWebhookError(w, "payment method event", err)
			return
		}
		if paymentMethod.Customer == nil {
			handleWebhookError(w, "inserting payment method", fmt.Errorf("No customer for payment method"))
			return
		}
		customerID := bgstripe.CustomerID(paymentMethod.Customer.ID)
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, err := bgstripe.GetUserIDForStripeCustomerID(tx, customerID)
			if err != nil {
				return err
			}
			if err := bgstripe.InsertPaymentMethod(tx, *userID, &paymentMethod); err != nil {
				return err
			}
			stripeSubscription, err := bgstripe.LookupActiveSubscriptionForUser(tx, *userID)
			if err != nil {
				return err
			}
			if stripeSubscription != nil && stripeSubscription.PaymentState == bgstripe.PaymentStateTrialNoPaymentMethod {
				if err := bgstripe.UpdatePaymentStateForSubscription(tx, *userID, stripeSubscription.StripeSubscriptionID, bgstripe.PaymentStateTrialPaymentMethodAdded); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			handleWebhookError(w, "capturing payment method event", err)
			return
		}
	case "payment_method.detached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			handleWebhookError(w, "payment method event", err)
			return
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := bgstripe.RemovePaymentMethod(tx, bgstripe.PaymentMethodID(paymentMethod.ID)); err != nil {
				return err
			}
			return nil
		}); err != nil {
			handleWebhookError(w, "capturing payment method event", err)
			return
		}
	case "invoice.paid":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			handleWebhookError(w, "invoice event", err)
			return
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, subscription, err := bgstripe.ReconcileInvoice(tx, invoice)
			if err != nil {
				return err
			}
			return useraccounts.UpdateSubscriptionExpirationTime(tx, *userID, subscription.CurrentPeriodEnd)
		}); err != nil {
			handleWebhookError(w, "capturing payment method event", err)
			return
		}
	case "invoice.payment_failed",
		"invoice.payment_action_required":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			handleWebhookError(w, "invoice event", err)
			return
		}
		var subscription *bgstripe.Subscription
		var enqueuedNotificationRequest bool
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			var userID *users.UserID
			userID, subscription, err = bgstripe.ReconcileInvoice(tx, invoice)
			if err != nil {
				return err
			}
			holdUntilTime := time.Now().Add(5 * time.Minute)
			enqueuedNotificationRequest, err = useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePaymentError, holdUntilTime)
			return err
		}); err != nil {
			handleWebhookError(w, "invoice failure notification", err)
			return
		}
		log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", subscription.StripeSubscriptionID, enqueuedNotificationRequest))
	case "customer.subscription.trial_will_end":
		// Alert user that trial subscription will be ending
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			handleWebhookError(w, "subscription event", err)
			return
		}
		var enqueuedNotificationRequest bool
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, err := bgstripe.LookupBabblegraphUserIDForStripeSubscriptionID(tx, subscription.ID)
			if err != nil {
				return err
			}
			holdUntilTime := time.Now().Add(30 * time.Minute)
			enqueuedNotificationRequest, err = useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypeTrialEndingSoon, holdUntilTime)
			return err
		}); err != nil {
			handleWebhookError(w, "invoice failure notification", err)
			return
		}
		log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", subscription.ID, enqueuedNotificationRequest))
	case "customer.subscription.updated",
		"customer.subscription.deleted":
		// Handle state change
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			handleWebhookError(w, "subscription event", err)
			return
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, err := bgstripe.LookupBabblegraphUserIDForStripeSubscriptionID(tx, subscription.ID)
			if err != nil {
				return err
			}
			reconciliationAction, err := bgstripe.ReconcileSubscriptionUpdate(tx, subscription)
			if err != nil {
				return err
			}
			if reconciliationAction != nil {
				switch *reconciliationAction {
				case bgstripe.SubscriptionReconciliationActionCancellation:
					if err := useraccounts.ExpireSubscriptionForUser(tx, *userID); err != nil {
						return err
					}
					holdUntilTime := time.Now().Add(5 * time.Minute)
					enqueuedNotificationRequest, err := useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, holdUntilTime)
					if err != nil {
						return err
					}
					log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", subscription.ID, enqueuedNotificationRequest))
				case bgstripe.SubscriptionReconciliationActionFirstPaymentSuccessful:
					// Send notification that subscription has started
				default:
					// No-op
				}
			}
			return nil
		}); err != nil {
			handleWebhookError(w, "subscription update", err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func handleWebhookError(w http.ResponseWriter, webhookEventType string, err error) {
	fErr := fmt.Errorf("Error processing %s for stripe webhook: %s", webhookEventType, err.Error())
	envName := env.MustEnvironmentName()
	switch envName {
	case env.EnvironmentProd,
		env.EnvironmentStage:
		sentry.CaptureException(fErr)
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		log.Println(fErr.Error())
	default:
		log.Println(fmt.Sprintf("Unknown environment: %s", envName))
	}
	w.WriteHeader(http.StatusBadRequest)
}
