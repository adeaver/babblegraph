package verification

import (
	"babblegraph/model/routes"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func VerifyUserByToken(token string) (*users.UserID, error) {
	var userID users.UserID
	if err := encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != routes.UserVerificationKey.Str() {
			return fmt.Errorf("Token has wrong key: %s", t.Key)
		}
		userIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("Token has wrong value type")
		}
		userID = users.UserID(userIDStr)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		if err := users.SetUserStatusToVerified(tx, userID); err != nil {
			return err
		}
		return userreadability.InitializeReadingLevelClassification(tx, userID, wordsmith.LanguageCodeSpanish)
	}); err != nil {
		return nil, err
	}
	return &userID, nil
}
