package scheduler

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

const (
	numberOfEmailsToSendAt = 7
)

func sendUserFeedbackEmails(localSentryHub *sentry.Hub, emailClient *ses.Client) error {
	var activeUsers []users.User
	var dailyEmailUsages []email.EmailUsage
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		activeUsers, err = users.GetAllActiveUsers(tx)
		if err != nil {
			return err
		}
		dailyEmailUsages, err = email.GetEmailUsageForType(tx, email.EmailTypeDaily)
		return err
	}); err != nil {
		return err
	}
	verifiedUsersByIDHash := make(map[users.UserID]string)
	for _, user := range activeUsers {
		verifiedUsersByIDHash[user.ID] = user.EmailAddress
	}
	for _, usage := range dailyEmailUsages {
		emailAddress, ok := verifiedUsersByIDHash[usage.UserID]
		if !ok {
			log.Println(fmt.Sprintf("User %s is no longer verified. Continuing...", usage.UserID))
			continue
		}
		switch {
		case usage.NumberOfSentEmails != numberOfEmailsToSendAt:
			log.Println(fmt.Sprintf("User %s has been sent the correct number of emails. Continuing...", usage.UserID))
		case !usage.HasOpenedOneEmail:
			log.Println(fmt.Sprintf("User %s has not opened an email. Continuing...", usage.UserID))
		default:
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				_, err := email_actions.SendUserFeedbackEmailForRecipient(tx, emailClient, email.Recipient{
					EmailAddress: emailAddress,
					UserID:       usage.UserID,
				})
				return err
			}); err != nil {
				log.Println(fmt.Sprintf("Error fulfilling user feedback attempt for user %s: %s. Continuing...", usage.UserID, err.Error()))
				localSentryHub.CaptureException(err)
			}
		}
	}
	return nil
}
