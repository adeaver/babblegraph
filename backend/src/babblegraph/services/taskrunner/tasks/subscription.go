package tasks

import (
	email_actions "babblegraph/actions/email"
	email_model "babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/ses"
	"log"

	"github.com/jmoiron/sqlx"
)

func CreateUserWithBetaPremiumSubscription(emailClient *ses.Client, emailAddress string) error {
	formattedEmailAddress := email.FormatEmailAddress(emailAddress)
	return database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		if err != nil {
			return err
		}
		if user == nil {
			log.Println("No user found for email subscription")
			return nil
		}
		didUserHaveSubscriptionBeforeTask, wasSubscriptionAlreadyActive, err := useraccounts.DoesUserHaveSubscription(tx, user.ID)
		switch {
		case err != nil:
			return err
		case !didUserHaveSubscriptionBeforeTask:
			// User had no account previously
			log.Println("User had no account previously, handling...")
			if err := useraccounts.AddSubscriptionLevelForUser(tx, user.ID, useraccounts.SubscriptionLevelBetaPremium); err != nil {
				return err
			}
			if _, err := email_actions.SendUserCreationEmailForRecipient(tx, emailClient, email_model.Recipient{
				UserID:       user.ID,
				EmailAddress: formattedEmailAddress,
			}); err != nil {
				return err
			}
		case didUserHaveSubscriptionBeforeTask && !wasSubscriptionAlreadyActive:
			// User reactivated account
			log.Println("User reactivated account, handling...")
			if err := useraccounts.AddSubscriptionLevelForUser(tx, user.ID, useraccounts.SubscriptionLevelBetaPremium); err != nil {
				return err
			}
			if _, err := email_actions.SendAccountReactivationEmailForRecipient(tx, emailClient, email_model.Recipient{
				UserID:       user.ID,
				EmailAddress: formattedEmailAddress,
			}); err != nil {
				return err
			}
		case didUserHaveSubscriptionBeforeTask && wasSubscriptionAlreadyActive:
			// no-op
			log.Println("User already had an active account, skipping...")

		}
		return nil
	})
}

func DeactivateUserSubscriptionForUser(emailClient *ses.Client, emailAddress string) error {
	formattedEmailAddress := email.FormatEmailAddress(emailAddress)
	return database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		if err != nil {
			return err
		}
		if user == nil {
			log.Println("No user found for email subscription")
			return nil
		}
		didUserHaveSubscriptionBeforeTask, wasSubscriptionAlreadyActive, err := useraccounts.DoesUserHaveSubscription(tx, user.ID)
		switch {
		case err != nil:
			return err
		case !didUserHaveSubscriptionBeforeTask,
			!wasSubscriptionAlreadyActive:
			log.Println("User did not have an active subscription, skipping...")
		default:
			if err := useraccounts.ExpireSubscriptionForUser(tx, user.ID); err != nil {
				return err
			}
			if _, err := email_actions.SendAccountExpirationEmail(tx, emailClient, email_model.Recipient{
				UserID:       user.ID,
				EmailAddress: formattedEmailAddress,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
