package scheduler

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/model/userverificationattempt"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func handlePendingVerifications(emailClient *ses.Client) error {
	var userIDs []users.UserID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userIDs, err = userverificationattempt.GetUserIDsWithPendingVerificationAttempts(tx)
		return err
	}); err != nil {
		return err
	}
	for _, userID := range userIDs {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			user, err := users.GetUser(tx, userID)
			if err != nil {
				return err
			}
			// While this looks incorrect, we want to fail on this before
			// sending the email, rather than sending the email and then failing.
			// And because this is a transaction, if we fail to send the email
			// then this doesn't update either.
			if err := userverificationattempt.MarkVerificationAttemptAsFulfilledByUserID(tx, userID); err != nil {
				return err
			}
			_, err = email_actions.SendVerificationEmailForRecipient(tx, emailClient, email.Recipient{
				UserID:       userID,
				EmailAddress: user.EmailAddress,
			})
			return err
		}); err != nil {
			log.Println(fmt.Sprintf("Error fulfilling verification attempt for user %s: %s. Continuing...", userID, err.Error()))
		}
	}
	return nil
}
