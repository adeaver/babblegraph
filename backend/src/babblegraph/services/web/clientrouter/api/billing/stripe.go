package billing

import (
	"babblegraph/model/billing"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

type stripeBeginPaymentMethodSetupRequest struct{}

type stripeBeginPaymentMethodSetupResponse struct {
	SetupIntentClientSecret string `json:"client_secret"`
}

func stripeBeginPaymentMethodSetup(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var setupIntentClientSecret *string
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		setupIntentClientSecret, err = billing.GetSetupIntentClientSecretForUser(tx, userAuth.UserID)
		return err
	}); err != nil {
		return nil, err
	}
	return stripeBeginPaymentMethodSetupResponse{
		SetupIntentClientSecret: *setupIntentClientSecret,
	}, nil
}

type stripeHandleWebhookEventResponse struct{}

func stripeHandleWebhookEvent(r *router.Request) (interface{}, error) {
	bodyBytes, err := r.GetBodyAsBytes()
	if err != nil {
		return nil, err
	}
	stripeSignature := r.GetHeader("Stripe-Signature")
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return billing.HandleStripeEvent(r, tx, stripeSignature, bodyBytes)
	}); err != nil {
		return nil, err
	}
	return stripeHandleWebhookEventResponse{}, nil
}
