package scheduler

import (
	"babblegraph/model/admin"
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func handleSendAdminTwoFactorAuthenticationCode(c async.Context) {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	var unfulfilledCodes []admin.TwoFactorAuthenticationCode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		unfulfilledCodes, err = admin.GetUnfulfilledTwoFactorAuthenticationAttempts(tx)
		return err
	}); err != nil {
		c.Errorf("Error getting unfulfilled 2FA codes: %s", err.Error())
		return
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
			c.Errorf("Error fulfilling attempt for admin ID %s: %s", code.AdminUserID, err.Error())
		}
	}
	return
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
