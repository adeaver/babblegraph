package scheduler

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/model/userverificationattempt"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ses"

	"github.com/jmoiron/sqlx"
)

func handlePendingVerifications(c async.Context) {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	var userIDs []users.UserID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userIDs, err = userverificationattempt.GetUserIDsWithPendingVerificationAttempts(tx)
		return err
	}); err != nil {
		c.Errorf("Error getting pending verification attempts: %s", err.Error())
		return
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
			c.Errorf("Error fulfilling verification attempt for user %s: %s. Continuing...", userID, err.Error())
		}
	}
}
