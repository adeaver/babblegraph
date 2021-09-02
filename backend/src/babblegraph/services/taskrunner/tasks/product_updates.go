package tasks

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/routes"
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

func SendProductUpdates() {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	var activeUsers []users.User
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		activeUsers, err = users.GetAllActiveUsers(tx)
		return err
	}); err != nil {
		sentry.CaptureException(err)
		return
	}
	for _, u := range activeUsers {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			premiumInformationLink, err := routes.MakePremiumInformationLink(u.ID)
			if err != nil {
				return err
			}
			_, err = email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, emailClient, email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
				EmailType:     email.EmailTypePremiumAnnouncement,
				FromEmailName: ptr.String("Andrew from Babblegraph"),
				Subject:       "Introducing Babblegraph Premium",
				EmailTitle:    "Introducing Babblegraph Premium",
				PreheaderText: "Learn more about the features you can unlock with Babblegraph Premium",
				Recipient: email.Recipient{
					EmailAddress: u.EmailAddress,
					UserID:       u.ID,
				},
				BeforeParagraphs: []string{
					"Hola!",
					"Andrew here! As you may or may not know, Babblegraph is a one-person operation and is completely independent. There’s no Silicon Valley venture capital money, no team of engineers, no lavish offices! Just me!",
					"To keep Babblegraph independent, support myself, and cover the costs of running Babblegraph, I’m introducing a premium subscription tier which gives subscribers access to exclusive features to enhace their Babblegraph experience!",
					"Some of the features include:",
					"The ability to pick how many articles you receive in each newsletter, and which topics they cover. Want to make sure every email has some articles on cooking? You can do that!",
					"Pick which days you receive your newsletter and which days you don’t.",
					"Spotlight words that are on your tracking list by having Babblegraph prominently display an article guaranteed to use a word you’re practicing.",
				},
				GenericEmailAction: &email_actions.GenericEmailAction{
					Link:       *premiumInformationLink,
					ButtonText: "Learn more about Babblegraph Premium",
				},
				AfterParagraphs: []string{
					"If the premium features aren’t enough to convince you, I’ll be continuing to develop new features for both Premium and non-Premium subscribers.",
					"If you’re not interested, then you can safely ignore this email and continue using the Babblegraph daily newsletter!",
					"And if you have any questions, you can always reply to this email!",
				},
				ShouldDedupeByType: true,
			})
			return err
		}); err != nil {
			log.Println(fmt.Sprintf("Got error sending product updates to %s: %s", u.EmailAddress, err.Error()))
			sentry.CaptureException(fmt.Errorf("Got error sending product updates to %s: %s", u.EmailAddress, err.Error()))
			continue
		}
		log.Println(fmt.Sprintf("Sent product updates to %s", u.EmailAddress))
	}
}
