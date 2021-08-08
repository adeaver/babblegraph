package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/users"
	"babblegraph/util/database"

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
