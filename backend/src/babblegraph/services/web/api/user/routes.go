package user

import (
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"encoding/json"
	"fmt"

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
	var didUpdate bool
	if err := encrypt.WithDecodedToken(r.Token, func(t encrypt.TokenPair) error {
		if t.Key != "unsubscribe" {
			return fmt.Errorf("incorrect key")
		}
		userID, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			didUpdate, err = users.UnsubscribeUserForIDAndEmail(tx, users.UserID(userID), r.EmailAddress)
			return err
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return unsubscribeUserResponse{
		Success: didUpdate,
	}, nil
}
