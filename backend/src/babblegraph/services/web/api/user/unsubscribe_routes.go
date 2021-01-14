package user

import (
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
		return err
	}); err != nil {
		return nil, err
	}
	return unsubscribeUserResponse{
		Success: didUpdate,
	}, nil
}
