package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/api"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/ptr"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return api.RegisterRouteGroup(api.RouteGroup{
		Prefix: "user",
		Routes: []api.Route{
			{
				Path:             "signup_user_1",
				Handler:          handleSignupUser,
				TrackEventWithID: ptr.String("signup"),
			}, {
				Path:    "get_user_content_topics_for_token_1",
				Handler: handleGetUserContentTopicsForToken,
			}, {
				Path:    "update_user_content_topics_for_token_1",
				Handler: handleUpdateUserContentTopicsForToken,
			}, {
				Path:    "handle_request_password_reset_link_1",
				Handler: requestPasswordResetLink,
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
			formattedEmailAddress := strings.ToLower(strings.Trim(*emailAddress, " "))
			user, err := users.LookupUserForIDAndEmail(tx, userID, formattedEmailAddress)
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
