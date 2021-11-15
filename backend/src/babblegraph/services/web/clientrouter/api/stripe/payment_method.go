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

type getPaymentMethodsForUserRequest struct{}

type getPaymentMethodsForUserResponse struct {
	PaymentMethods []bgstripe.PaymentMethod `json:"payment_methods"`
}

func getPaymentMethodsForUser(userID users.UserID, body []byte) (interface{}, error) {
	var paymentMethods []bgstripe.PaymentMethod
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		paymentMethods, err = bgstripe.GetPaymentMethodsForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getPaymentMethodsForUserResponse{
		PaymentMethods: paymentMethods,
	}, nil
}

type getPaymentMethodByIDRequest struct {
	StripePaymentMethodID bgstripe.PaymentMethodID `json:"stripe_payment_method_id"`
}

type getPaymentMethodByIDResponse struct {
	PaymentMethod *bgstripe.PaymentMethod `json:"payment_method,omitempty"`
}

func getPaymentMethodByID(userID users.UserID, body []byte) (interface{}, error) {
	var req getPaymentMethodByIDRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var paymentMethod *bgstripe.PaymentMethod
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		paymentMethod, err = bgstripe.LookupPaymentMethod(tx, userID, req.StripePaymentMethodID)
		return err
	}); err != nil {
		return nil, err
	}
	return getPaymentMethodByIDResponse{
		PaymentMethod: paymentMethod,
	}, nil
}

type deletePaymentMethodForUserError string

const (
	deletePaymentMethodForUserErrorDefault  deletePaymentMethodForUserError = "no-delete-default"
	deletePaymentMethodForUserErrorOnlyCard deletePaymentMethodForUserError = "only-card"
)

func (d deletePaymentMethodForUserError) Ptr() *deletePaymentMethodForUserError {
	return &d
}

type deletePaymentMethodForUserRequest struct {
	StripePaymentMethodID bgstripe.PaymentMethodID `json:"stripe_payment_method_id"`
}

type deletePaymentMethodForUserResponse struct {
	Error *deletePaymentMethodForUserError `json:"error,omitempty"`
}

func deletePaymentMethodForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req deletePaymentMethodForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var deleteErr *deletePaymentMethodForUserError
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		paymentMethods, err := bgstripe.GetPaymentMethodsForUser(tx, userID)
		switch {
		case err != nil:
			return err
		case len(paymentMethods) == 1:
			deleteErr = deletePaymentMethodForUserErrorOnlyCard.Ptr()
			return nil
		}
		for _, p := range paymentMethods {
			if p.StripePaymentMethodID == req.StripePaymentMethodID {
				if p.IsDefault {
					deleteErr = deletePaymentMethodForUserErrorDefault.Ptr()
					return nil
				} else {
					break
				}
			}
		}
		return bgstripe.CancelPaymentMethodAndRemoveForUser(tx, userID, req.StripePaymentMethodID)
	}); err != nil {
		return nil, err
	}
	return deletePaymentMethodForUserResponse{
		Error: deleteErr,
	}, nil
}
