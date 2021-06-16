package email

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
)

const (
	verificationEmailTemplateFilename string = "user_verification_template.html"
)

type userVerificationEmailTemplate struct {
	email.BaseEmailTemplate
	VerificationLink string
}

func SendVerificationEmailForRecipient(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient) (*email.ID, error) {
	emailRecordID := email.NewEmailRecordID()
	if err := email.InsertEmailRecord(tx, emailRecordID, recipient.UserID, email.EmailTypeUserVerification); err != nil {
		return nil, err
	}
	template, err := createUserVerificationEmailTemplate(tx, emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	body, err := createVerificationEmailBody(*template)
	if err != nil {
		return nil, err
	}
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient: recipient.EmailAddress,
		HTMLBody:  *body,
		Subject:   "Verify your Babblegraph Subscription",
	})
	if err != nil {
		return nil, err
	}
	if err := email.UpdateEmailRecordIDWithSESMessageID(tx, emailRecordID, *sesMessageID); err != nil {
		// **** VERY IMPORTANT HERE ****
		// This *cannot* return an error if it fails
		// since SES has already successfully sent the email and returning an error
		// causes the transaction to abort and rollback. However, SES has a side effect.
		// The only thing that this does is update the email record
		// to have the SES message ID. It is not super important - but it is useful.
		log.Println(fmt.Sprintf("Error updating email record %s with SES Message ID %s: %s", emailRecordID, *sesMessageID, err.Error()))
	}
	return &emailRecordID, nil
}

func createUserVerificationEmailTemplate(tx *sqlx.Tx, emailRecordID email.ID, recipient email.Recipient) (*userVerificationEmailTemplate, error) {
	baseTemplate, err := createBaseTemplate(tx, emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	verificationLink, err := routes.MakeUserVerificationLink(recipient.UserID)
	if err != nil {
		return nil, err
	}
	return &userVerificationEmailTemplate{
		BaseEmailTemplate: *baseTemplate,
		VerificationLink:  *verificationLink,
	}, nil
}

func createVerificationEmailBody(templateData userVerificationEmailTemplate) (*string, error) {
	templateFile, err := getPathForTemplateFile(verificationEmailTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(verificationEmailTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, templateData); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
