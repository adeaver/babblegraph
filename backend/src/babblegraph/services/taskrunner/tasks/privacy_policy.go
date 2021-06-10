package tasks

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func SendPrivacyPolicyUpdate() {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		activeUsers, err := users.GetAllActiveUsers(tx)
		if err != nil {
			return err
		}
		for _, u := range activeUsers {
			if _, err := email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, emailClient, email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
				EmailType:     email.EmailTypePrivacyPolicyUpdateJune2021,
				FromEmailName: ptr.String("Andrew from Babblegraph"),
				Subject:       "Updated Privacy Policy",
				EmailTitle:    "Privacy Policy Update",
				Recipient: email.Recipient{
					EmailAddress: u.EmailAddress,
					UserID:       u.ID,
				},
				BeforeParagraphs: []string{
					"Hola!",
					"As a result of a recent feature release, I’ve updated the privacy policy.",
					"Starting now, Babblegraph will be collecting information when links are clicked in the daily email. This helps keep track of whether or not you’ve reached your monthly quota of free articles for news sources that enforce a limit. Likewise, Babblegraph will now attempt to completely avoid sending all content that requires a paid subscription to access!",
					"Take a minute to read the privacy policy, which is linked below.",
				},
				GenericEmailAction: &email_actions.GenericEmailAction{
					Link:       env.GetAbsoluteURLForEnvironment("privacy-policy"),
					ButtonText: "Read Babblegraph’s Privacy Policy",
				},
				AfterParagraphs: []string{
					"This does not change Babblegraph’s policy of not selling user data.",
					"As usual, if you have any questions, concerns, or comments, you can respond directly to this email!",
					"Gracias!",
				},
				PreheaderText: "Babblegraph’s privacy policy has been updated!",
			}); err != nil {
				log.Println(fmt.Sprintf("Got error sending privacy policy update to %s: %s", u.EmailAddress, err.Error()))
				sentry.CaptureException(fmt.Errorf("Got error sending privacy policy update to %s: %s", u.EmailAddress, err.Error()))
				continue
			}
			log.Println(fmt.Sprintf("Sent privacy policy update to %s", u.EmailAddress))
		}
		return nil
	}); err != nil {
		sentry.CaptureException(err)
	}
}
