package user

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type unsubscribeUserRequest struct {
	Token        string `json:"token"`
	EmailAddress string `json:"email_address"`
}

type unsubscribeUserResponse struct {
	Success bool `json:"success"`
}

func handleUnsubscribeUser(body []byte) (interface{}, error) {
	var r unsubscribeUserRequest
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	userID, err := parseSubscriptionManagementToken(r.Token, ptr.String(r.EmailAddress))
	if err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		didUpdate, err = users.UnsubscribeUserForIDAndEmail(tx, *userID, r.EmailAddress)
		if err != nil {
			return err
		}
		subscription, err := bgstripe.LookupActiveSubscriptionForUser(tx, *userID)
		switch {
		case err != nil:
			return err
		case subscription == nil:
			return nil
		default:
			return bgstripe.CancelSubscription(tx, *userID)
		}
	}); err != nil {
		return nil, err
	}
	return unsubscribeUserResponse{
		Success: didUpdate,
	}, nil
}
