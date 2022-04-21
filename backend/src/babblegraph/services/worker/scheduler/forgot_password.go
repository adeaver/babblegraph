package scheduler

import (
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func handlePendingForgotPasswordAttempts(c async.Context) {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	var forgotPasswordAttempts []useraccounts.ForgotPasswordAttempt
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		forgotPasswordAttempts, err = useraccounts.GetAllUnfulfilledForgotPasswordAttempts(tx)
		return err
	}); err != nil {
		c.Errorf("Error getting forgot password attempts: %s", err.Error())
		return
	}
	for _, attempt := range forgotPasswordAttempts {
		// This is intentionally two different transactions
		// so that we can abort on a single forgot password without
		// having to abort all subsequent transactions
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			err := useraccounts.FulfillForgotPasswordAttempt(tx, attempt.ID)
			if err != nil {
				return err
			}
			user, err := users.GetUser(tx, attempt.UserID)
			switch {
			case err != nil:
				return err
			case user == nil:
				return fmt.Errorf("Expected user from ID %s, but got none", attempt.UserID)
			case user.Status != users.UserStatusVerified:
				c.Infof("User with ID %s does not have verified status. Skipping forgot password attempt", attempt.UserID)
				return nil
			}
			return sendForgotPasswordEmailForUserAndAttemptID(tx, emailClient, *user, attempt.ID)
		}); err != nil {
			c.Errorf("Error sending forgot password attempt for user %s: %s", attempt.UserID, err.Error())
			continue
		}
	}
}

func sendForgotPasswordEmailForUserAndAttemptID(tx *sqlx.Tx, sesClient *ses.Client, user users.User, forgotPasswordAttemptID useraccounts.ForgotPasswordAttemptID) error {
	passwordResetLink, err := routes.MakeForgotPasswordLink(forgotPasswordAttemptID)
	if err != nil {
		return err
	}
	emailRecordID := email.NewEmailRecordID()
	if err := email.InsertEmailRecord(tx, emailRecordID, user.ID, email.EmailTypePasswordReset); err != nil {
		return err
	}
	userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, user.ID)
	if err != nil {
		return err
	}
	emailHTML, err := emailtemplates.MakeGenericUserEmailHTML(emailtemplates.MakeGenericUserEmailHTMLInput{
		EmailRecordID: emailRecordID,
		UserAccessor:  userAccessor,
		EmailTitle:    "Password Reset",
		PreheaderText: "Reset your password for your Babblegraph account",
		BeforeParagraphs: []string{
			"Hola!",
			"There was recently a request made to reset your password.",
			"If you did not make this request, you do not need to do anything. No one has access to your account.",
			"If you did make this request, you can reset your password using the button below.",
		},
		GenericEmailAction: &emailtemplates.GenericEmailAction{
			Link:       *passwordResetLink,
			ButtonText: "Reset your password",
		},
		AfterParagraphs: []string{
			"This link is only valid for 15 minutes after the delivery time.",
		},
	})
	if err != nil {
		return err
	}
	return email.SendEmailWithHTMLBody(tx, sesClient, email.SendEmailWithHTMLBodyInput{
		ID:           emailRecordID,
		EmailAddress: user.EmailAddress,
		Subject:      "Reset your Babblegraph Account Password",
		Body:         *emailHTML,
	})
}

func handleArchiveForgotPasswordAttempts(c async.Context) {
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return useraccounts.ArchiveAllForgotPasswordAttemptsOlderThan20Minutes(tx)
	}); err != nil {
		c.Errorf("Error archiving forgot password attempts: %s", err.Error())
	}
}
