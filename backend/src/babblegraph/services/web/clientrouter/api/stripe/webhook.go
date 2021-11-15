package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return bgstripe.HandleStripeEvent(tx, event)
	}); err != nil {
		handleWebhookError(w, event.Type, err)
		return
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
