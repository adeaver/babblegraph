package scheduler

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func handlePendingForgotPasswordAttempts(localSentryHub *sentry.Hub, emailClient *ses.Client) error {
	var forgotPasswordAttempts []useraccounts.ForgotPasswordAttempt
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		forgotPasswordAttempts, err = useraccounts.GetAllUnfulfilledForgotPasswordAttempts(tx)
		return err
	}); err != nil {
		localSentryHub.CaptureException(err)
		return err
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
				log.Println(fmt.Sprintf("User with ID %s does not have verified status. Skipping forgot password attempt", attempt.UserID))
				return nil
			}
			return sendForgotPasswordEmailForUserAndAttemptID(tx, emailClient, *user, attempt.ID)
		}); err != nil {
			localSentryHub.CaptureException(fmt.Errorf("Error sending forgot password attempt for user %s: %s", attempt.UserID, err.Error()))
			continue
		}
	}
	return nil
}

func sendForgotPasswordEmailForUserAndAttemptID(tx *sqlx.Tx, sesClient *ses.Client, user users.User, forgotPasswordAttemptID useraccounts.ForgotPasswordAttemptID) error {
	passwordResetLink, err := routes.MakeForgotPasswordLink(forgotPasswordAttemptID)
	if err != nil {
		return err
	}
	if _, err := email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, sesClient, email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType: email.EmailTypePasswordReset,
		Recipient: email.Recipient{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
		},
		Subject:       "Reset your Babblegraph Account Password",
		EmailTitle:    "Password Reset",
		PreheaderText: "Reset your password for your Babblegraph account",
		BeforeParagraphs: []string{
			"Hola!",
			"There was recently a request made to reset your password.",
			"If you did not make this request, you do not need to do anything. No one has access to your account.",
			"If you did make this request, you can reset your password using the button below.",
		},
		GenericEmailAction: &email_actions.GenericEmailAction{
			Link:       *passwordResetLink,
			ButtonText: "Reset your password",
		},
	}); err != nil {
		return err
	}
	return nil
}

func handleArchiveForgotPasswordAttempts() error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		return useraccounts.ArchiveAllForgotPasswordAttemptsOlderThan20Minutes(tx)
	})
}
