package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type getSetupIntentForUserRequest struct{}

type getSetupIntentForUserResponse struct {
	SetupIntentID string `json:"setup_intent_id"`
	ClientSecret  string `json:"client_secret"`
}

func getSetupIntentForUser(userID users.UserID, body []byte) (interface{}, error) {
	var addPaymentMethodCreds *bgstripe.AddPaymentMethodCredentials
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		addPaymentMethodCreds, err = bgstripe.GetAddPaymentMethodCredentialsForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getSetupIntentForUserResponse{
		SetupIntentID: addPaymentMethodCreds.SetupIntentID,
		ClientSecret:  addPaymentMethodCreds.ClientSecret,
	}, nil
}

type insertNewPaymentMethodForUserRequest struct {
	StripePaymentMethodID bgstripe.PaymentMethodID `json:"stripe_payment_method_id"`
}

type insertNewPaymentMethodForUserResponse struct {
	PaymentMethod bgstripe.PaymentMethod `json:"payment_method"`
}

func insertNewPaymentMethodForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req insertNewPaymentMethodForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var paymentMethod *bgstripe.PaymentMethod
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		paymentMethod, err = bgstripe.FindStripePaymentMethodAndInsert(tx, userID, req.StripePaymentMethodID)
		return err
	}); err != nil {
		return nil, err
	}
	return insertNewPaymentMethodForUserResponse{
		PaymentMethod: *paymentMethod,
	}, nil
}

type setDefaultPaymentMethodForUserRequest struct {
	StripePaymentMethodID bgstripe.PaymentMethodID `json:"stripe_payment_method_id"`
}

type setDefaultPaymentMethodForUserResponse struct{}

func setDefaultPaymentMethodForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req setDefaultPaymentMethodForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		stripeCustomer, err := bgstripe.GetStripeCustomerForUserID(tx, userID)
		if err != nil {
			return err
		}
		return bgstripe.SetDefaultPaymentMethodForCustomer(tx, stripeCustomer.CustomerID, req.StripePaymentMethodID)
	}); err != nil {
		return nil, err
	}
	return setDefaultPaymentMethodForUserResponse{}, nil
}
