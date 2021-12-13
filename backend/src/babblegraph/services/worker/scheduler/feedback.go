package scheduler

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ses"

	"github.com/jmoiron/sqlx"
)

const (
	numberOfEmailsToSendAt = 7
)

func sendUserFeedbackEmails(c async.Context) {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
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
		c.Errorf("Error getting users and usage: %s", err.Error())
		return
	}
	verifiedUsersByIDHash := make(map[users.UserID]string)
	for _, user := range activeUsers {
		verifiedUsersByIDHash[user.ID] = user.EmailAddress
	}
	for _, usage := range dailyEmailUsages {
		emailAddress, ok := verifiedUsersByIDHash[usage.UserID]
		if !ok {
			c.Infof("User %s is no longer verified. Continuing...", usage.UserID)
			continue
		}
		switch {
		case usage.NumberOfSentEmails != numberOfEmailsToSendAt:
			c.Infof("User %s has been sent the correct number of emails. Continuing...", usage.UserID)
		case !usage.HasOpenedOneEmail:
			c.Infof("User %s has not opened an email. Continuing...", usage.UserID)
		default:
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				_, err := email_actions.SendUserFeedbackEmailForRecipient(tx, emailClient, email.Recipient{
					EmailAddress: emailAddress,
					UserID:       usage.UserID,
				})
				return err
			}); err != nil {
				c.Errorf("Error fulfilling user feedback attempt for user %s: %s. Continuing...", usage.UserID, err.Error())
			}
		}
	}
}
