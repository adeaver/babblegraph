package scheduler

import (
	"babblegraph/admin/model/auth"
	"babblegraph/admin/model/user"
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func handleSendAdminTwoFactorAuthenticationCode(localSentryHub *sentry.Hub, emailClient *ses.Client) error {
	var unfulfilledCodes []auth.Admin2FACode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		unfulfilledCodes, err = auth.GetUnfulfilledTwoFactorAuthenticationAttempts(tx)
		return err
	}); err != nil {
		localSentryHub.CaptureException(err)
		return err
	}
	for _, code := range unfulfilledCodes {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			adminUser, err := user.GetAdminUser(tx, code.AdminUserID)
			if err != nil {
				return err
			}
			if err := auth.FulfillTwoFactorAuthenticationAttempt(tx, adminUser.AdminID, code.Code); err != nil {
				return err
			}
			twoFactorEmailHTML, err := emailtemplates.MakeGenericEmailHTML(emailtemplates.MakeGenericEmailHTMLInput{
				EmailTitle:    "Requested Two Factor Code",
				PreheaderText: "This is a requested two factor authentication code",
				BeforeParagraphs: []string{
					fmt.Sprintf("Your code is %s", code.Code),
				},
			})
			if err != nil {
				return err
			}
			return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
				// This is a bit of a hack, technically, this will update an email that doesn't exist
				ID:              email.NewEmailRecordID(),
				EmailAddress:    adminUser.EmailAddress,
				Subject:         "Two Factor Authentication Code",
				EmailSenderName: ptr.String("Babblegraph Admin"),
				Body:            *twoFactorEmailHTML,
			})
		}); err != nil {
			localSentryHub.CaptureException(err)
		}
	}
	return nil
}

func handleCleanUpAdminTwoFactorCodesAndAccessTokens() error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		if err := auth.RemoveExpiredAccessTokens(tx); err != nil {
			return err
		}
		return auth.RemoveExpiredTwoFactorCodes(tx)
	})
}
