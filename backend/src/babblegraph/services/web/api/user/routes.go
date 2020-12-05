package user

import (
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "user",
		Routes: []router.Route{
			{
				Path:    "unsubscribe_user_1",
				Handler: handleUnsubscribeUser,
			},
		},
	})
}

type unsubscribeUserRequest struct {
	UserID       users.UserID `json:"user_id"`
	EmailAddress string       `json:"email_address"`
}

type unsubscribeUserResponse struct {
	Success bool `json:"success"`
}

func handleUnsubscribeUser(body []byte) (interface{}, error) {
	var r unsubscribeUserRequest
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		didUpdate, err = users.UnsubscribeUserForIDAndEmail(tx, r.UserID, r.EmailAddress)
		return err
	}); err != nil {
		return nil, err
	}
	return unsubscribeUserResponse{
		Success: didUpdate,
	}, nil
}
