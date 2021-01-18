package scheduler

import (
	"babblegraph/model/users"
	"babblegraph/model/userverificationattempt"
	"babblegraph/util/database"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func handlePendingVerifications() error {
	var userIDs []users.UserID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userIDs, err = userverificationattempt.GetUserIDsWithPendingVerificationAttempts(tx)
		return err
	}); err != nil {
		return err
	}
	for _, userID := range userIDs {
		// TODO: create email template and send it here
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			return userverificationattempt.MarkVerificationAttemptAsFulfilledByUserID(tx, userID)
		}); err != nil {
			log.Println(fmt.Sprintf("Error fulfilling verification attempt for user %s: %s. Continuing...", userID, err.Error()))
		}
	}
	return nil
}
