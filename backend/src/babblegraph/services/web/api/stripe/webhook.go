package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/useraccounts"
	"babblegraph/model/useraccountsnotifications"
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
	"github.com/stripe/stripe-go/v72/sub"
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
		// Add Payment Method, Update Subscription State
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			handleWebhookError(w, "payment method event", err)
			return
		}
	case "payment_method.detached":
		// Remove payment method, update subscription state.
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			handleWebhookError(w, "payment method event", err)
			return
		}
	case "invoice.paid":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			handleWebhookError(w, "invoice event", err)
			return
		}
		if invoice.Subscription != nil {
			subscription, err := sub.Get(invoice.Subscription.ID, &stripe.SubscriptionParams{})
			if err != nil {
				handleWebhookError(w, "capture invoice success", err)
				return
			}
			stripeSubscriptionID := bgstripe.SubscriptionID(invoice.Subscription.ID)
			newExpirationTime := time.Unix(subscription.CurrentPeriodEnd, 0)
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				userID, err := bgstripe.LookupBabblegraphUserIDForStripeSubscriptionID(tx, stripeSubscriptionID)
				if err != nil {
					return err
				}
				return useraccounts.UpdateSubscriptionExpirationTime(tx, *userID, newExpirationTime)
			}); err != nil {
				handleWebhookError(w, "capture invoice success", err)
				return
			}
		}
	case "invoice.payment_failed",
		"invoice.payment_action_required":
		// Alert user
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			handleWebhookError(w, "invoice event", err)
			return
		}
		stripeSubscriptionID := bgstripe.SubscriptionID(invoice.Subscription.ID)
		var enqueuedNotificationRequest bool
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, err := bgstripe.LookupBabblegraphUserIDForStripeSubscriptionID(tx, stripeSubscriptionID)
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
		log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", stripeSubscriptionID, enqueuedNotificationRequest))
	case "customer.subscription.trial_will_end":
		// Alert user that trial subscription will be ending
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			handleWebhookError(w, "subscription event", err)
			return
		}
		stripeSubscriptionID := bgstripe.SubscriptionID(subscription.ID)
		var enqueuedNotificationRequest bool
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, err := bgstripe.LookupBabblegraphUserIDForStripeSubscriptionID(tx, stripeSubscriptionID)
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
		log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", stripeSubscriptionID, enqueuedNotificationRequest))
	case "customer.subscription.updated":
		// Handle state change
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			handleWebhookError(w, "subscription event", err)
			return
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			// TODO: If moving from stage 0 to 3, then we need to send an email alerting them that their subscription is active
			return bgstripe.HandleStripeSubscriptionStatusUpdate(tx, bgstripe.SubscriptionID(subscription.ID), subscription.Status)
		}); err != nil {
			handleWebhookError(w, "subscription update", err)
			return
		}
	case "customer.subscription.deleted":
		// Mark the subscription as terminated
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			handleWebhookError(w, "subscription event", err)
			return
		}
		stripeSubscriptionID := bgstripe.SubscriptionID(subscription.ID)
		var enqueuedNotificationRequest bool
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			userID, err := bgstripe.LookupBabblegraphUserIDForStripeSubscriptionID(tx, stripeSubscriptionID)
			if err != nil {
				return err
			}
			if err := useraccounts.ExpireSubscriptionForUser(tx, *userID); err != nil {
				return err
			}
			holdUntilTime := time.Now().Add(5 * time.Minute)
			enqueuedNotificationRequest, err = useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, holdUntilTime)
			if err != nil {
				return err
			}
			return bgstripe.HandleStripeSubscriptionStatusUpdate(tx, stripeSubscriptionID, stripe.SubscriptionStatusCanceled)
		}); err != nil {
			handleWebhookError(w, "subscription deleted", err)
			return
		}
		log.Println(fmt.Sprintf("Did enqueue message for subscription %s: %t", stripeSubscriptionID, enqueuedNotificationRequest))
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
