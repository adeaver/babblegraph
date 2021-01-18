package email

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
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
	template, err := createUserVerificationEmailTemplate(emailRecordID, recipient)
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
	if err := email.InsertEmailRecord(tx, emailRecordID, *sesMessageID, recipient.UserID, email.EmailTypeUserVerification); err != nil {
		return nil, err
	}
	return &emailRecordID, nil
}

func createUserVerificationEmailTemplate(emailRecordID email.ID, recipient email.Recipient) (*userVerificationEmailTemplate, error) {
	baseTemplate, err := createBaseTemplate(emailRecordID, recipient)
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
