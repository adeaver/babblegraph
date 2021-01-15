package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "user",
		Routes: []router.Route{
			{
				Path:    "unsubscribe_user_1",
				Handler: handleUnsubscribeUser,
			}, {
				Path:    "get_user_preferences_for_token_1",
				Handler: handleGetUserPreferencesForToken,
			}, {
				Path:    "update_user_preferences_for_token_1",
				Handler: handleUpdateUserPreferencesForToken,
			}, {
				Path:    "get_user_content_topics_for_token_1",
				Handler: handleGetUserContentTopicsForToken,
			}, {
				Path:    "update_user_content_topics_for_token_1",
				Handler: handleUpdateUserContentTopicsForToken,
			},
		},
	})
}

func parseSubscriptionManagementToken(token string, emailAddress *string) (*users.UserID, error) {
	var userID users.UserID
	if err := encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != routes.SubscriptionManagementRouteEncryptionKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		userIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		userID = users.UserID(userIDStr)
		return nil
	}); err != nil {
		return nil, err
	}
	if emailAddress != nil {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			user, err := users.LookupUserForIDAndEmail(tx, userID, strings.ToLower(*emailAddress))
			if err != nil {
				return err
			}
			if user == nil {
				return fmt.Errorf("Invalid email address for token")
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return &userID, nil
}
