package scheduler

import (
	"babblegraph/model/billing"
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/routes"
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
			if err := userverificationattempt.MarkVerificationAttemptAsFulfilledByUserID(tx, userID); err != nil {
				return err
			}
			verificationLink, err := routes.MakeUserVerificationLink(userID)
			if err != nil {
				return err
			}
			emailRecordID := email.NewEmailRecordID()
			if err := email.InsertEmailRecord(tx, emailRecordID, user.ID, email.EmailTypeUserVerification); err != nil {
				return err
			}
			emailUserAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, user.ID)
			if err != nil {
				return err
			}
			emailInput := emailtemplates.MakeGenericUserEmailHTMLInput{
				EmailRecordID: emailRecordID,
				UserAccessor:  emailUserAccessor,
				EmailTitle:    "Verify your Babblegraph Subscription",
				PreheaderText: "Before you can begin receiving Babblegraph emails, you’ll need to verify your subscription to Babblegraph",
				ExcludeFooter: true,
			}
			premiumTrialEligibility, err := billing.GetPremiumNewsletterSubscriptionTrialEligibilityForUser(tx, userID)
			switch {
			case err != nil:
				return err
			case *premiumTrialEligibility > 0:
				emailInput.BeforeParagraphs = []string{
					"Hola!",
					"Thanks for signing up for Babblegraph. Verifying your subscription is the last step to complete before you’ll officially start receiving emails with Spanish-language content from real media sources from Spain and Latin America.",
					"Your 30 day free trial will begin once your email address is verified.",
				}
				emailInput.GenericEmailAction = &emailtemplates.GenericEmailAction{
					Link:       *verificationLink,
					ButtonText: "Click here to verify your email address",
				}
				emailInput.AfterParagraphs = []string{
					"Your subscription will be verified once you click the button. You’ll begin receiving emails within the next day.",
					"If you didn't sign up for Babblegraph or changed your mind about receiving daily Spanish-language emails, then you don’t have to do anything with this email.",
					"Thanks again for signing up for Babblegraph!",
				}
			default:
				emailInput.BeforeParagraphs = []string{
					"Hola!",
					"Thanks for signing up for Babblegraph. Verifying your subscription is the last step to complete before you’ll officially start receiving emails with Spanish-language content from real media sources from Spain and Latin America.",
					"Our records indicate that you’re not eligible for a free 30 day trial of Babblegraph, because you’ve previously used Babblegraph. If you believe that is incorrect, then just respond to this email and we’ll fix that for you. Otherwise, the link below will take you to a checkout form where you can complete the transaction for your subscription.",
				}
				emailInput.GenericEmailAction = &emailtemplates.GenericEmailAction{
					Link:       *verificationLink,
					ButtonText: "Click here to verify your email address and proceed to checkout",
				}
				emailInput.AfterParagraphs = []string{
					"Your subscription will be verified once you click the button and you’ve provided payment information. You’ll begin receiving emails within the next day.",
					"If you didn't sign up for Babblegraph or changed your mind about receiving daily Spanish-language emails, then you don’t have to do anything with this email.",
					"Thanks again for signing up for Babblegraph!",
				}
			}
			emailHTML, err := emailtemplates.MakeGenericUserEmailHTML(emailInput)
			if err != nil {
				return err
			}
			return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
				ID:           emailRecordID,
				EmailAddress: user.EmailAddress,
				Subject:      "Verify your Babblegraph Subscription",
				Body:         *emailHTML,
			})
		}); err != nil {
			c.Errorf("Error fulfilling verification attempt for user %s: %s. Continuing...", userID, err.Error())
		}
	}
}
