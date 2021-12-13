package scheduler

import (
	"babblegraph/model/admin"
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func handleSendAdminTwoFactorAuthenticationCode(localSentryHub *sentry.Hub, emailClient *ses.Client) error {
	var unfulfilledCodes []admin.TwoFactorAuthenticationCode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		unfulfilledCodes, err = admin.GetUnfulfilledTwoFactorAuthenticationAttempts(tx)
		return err
	}); err != nil {
		localSentryHub.CaptureException(err)
		return err
	}
	for _, code := range unfulfilledCodes {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			adminUser, err := admin.GetAdminUser(tx, code.AdminUserID)
			if err != nil {
				return err
			}
			if err := admin.FulfillTwoFactorAuthenticationAttempt(tx, adminUser.ID, code.Code); err != nil {
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

func handleCleanUpAdminTwoFactorCodesAndAccessTokens(c async.Context) {
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		if err := admin.RemoveExpiredAccessTokens(tx); err != nil {
			return err
		}
		return admin.RemoveExpiredTwoFactorCodes(tx)
	}); err != nil {
		c.Errorf("Error cleaning up admin two factor codes and access tokens: %s", err.Error())
	}
}
