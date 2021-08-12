package verification

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"time"

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
		holdUntilTime := time.Now().Add(14 * 24 * time.Hour)
		if _, err := useraccountsnotifications.EnqueueNotificationRequest(tx, userID, useraccountsnotifications.NotificationTypeInitialPremiumInformation, holdUntilTime); err != nil {
			return err
		}
		return userreadability.InitializeReadingLevelClassification(tx, userID, wordsmith.LanguageCodeSpanish)
	}); err != nil {
		log.Println(fmt.Sprintf("Error verifying: %s", err.Error()))
		return nil, err
	}
	return &userID, nil
}
