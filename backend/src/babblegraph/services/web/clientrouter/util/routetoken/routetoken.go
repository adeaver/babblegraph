package routetoken

import (
	"babblegraph/model/routes"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func ValidateRouteToken(token string, expectedToken routes.RouteEncryptionKey) error {
	return encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != expectedToken.Str() {
			return fmt.Errorf("incorrect key")
		}
		return nil
	})
}

func ValidateTokenAndGetUserID(token string, expectedToken routes.RouteEncryptionKey) (*users.UserID, error) {
	var userID users.UserID
	if err := encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != expectedToken.Str() {
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
	return &userID, nil
}

func ValidateTokenAndEmailAndGetUserID(token string, expectedToken routes.RouteEncryptionKey, expectedEmailAddress string) (*users.UserID, error) {
	var userID users.UserID
	if err := encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != expectedToken.Str() {
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
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		formattedEmailAddress := strings.ToLower(strings.Trim(expectedEmailAddress, " "))
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
	return &userID, nil
}
