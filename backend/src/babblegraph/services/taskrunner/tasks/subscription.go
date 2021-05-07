package tasks

import (
	email_actions "babblegraph/actions/email"
	email_model "babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/ses"

	"github.com/jmoiron/sqlx"
)

func CreateUserWithBetaPremiumSubscription(emailClient *ses.Client, emailAddress string) error {
	formattedEmailAddress := email.FormatEmailAddress(emailAddress)
	return database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		if err != nil {
			return err
		}
		if err := useraccounts.AddSubscriptionLevelForUser(tx, user.ID, useraccounts.SubscriptionLevelBetaPremium); err != nil {
			return err
		}
		if _, err := email_actions.SendUserCreationEmailForRecipient(tx, emailClient, email_model.Recipient{
			UserID:       user.ID,
			EmailAddress: formattedEmailAddress,
		}); err != nil {
			return err
		}
		return nil
	})
}
