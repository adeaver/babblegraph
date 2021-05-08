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
	userCreationEmailTemplateFilename string = "create_user_template.html"
)

type userCreationEmailTemplate struct {
	email.BaseEmailTemplate
	SignupLink string
}

func SendUserCreationEmailForRecipient(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient) (*email.ID, error) {
	emailRecordID := email.NewEmailRecordID()
	template, err := createUserCreationEmailTemplate(emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	body, err := createUserCreationEmailBody(*template)
	if err != nil {
		return nil, err
	}
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient: recipient.EmailAddress,
		HTMLBody:  *body,
		Subject:   "Create your Babblegraph Account to Access Premium Features",
	})
	if err != nil {
		return nil, err
	}
	if err := email.InsertEmailRecord(tx, emailRecordID, *sesMessageID, recipient.UserID, email.EmailTypeUserCreation); err != nil {
		return nil, err
	}
	return &emailRecordID, nil
}

func createUserCreationEmailTemplate(emailRecordID email.ID, recipient email.Recipient) (*userCreationEmailTemplate, error) {
	baseTemplate, err := createBaseTemplate(emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	signupLink, err := routes.MakeUserCreationLink(recipient.UserID)
	if err != nil {
		return nil, err
	}
	return &userCreationEmailTemplate{
		BaseEmailTemplate: *baseTemplate,
		SignupLink:        *signupLink,
	}, nil
}

func createUserCreationEmailBody(templateData userCreationEmailTemplate) (*string, error) {
	templateFile, err := getPathForTemplateFile(userCreationEmailTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(userCreationEmailTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, templateData); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
