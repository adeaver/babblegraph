package email

import (
	"babblegraph/model/email"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
)

const (
	userFeedbackEmailTemplateFilename string = "feedback_template.html"

	feedbackEmailSenderName = "Andrew from Babblegraph"
)

type userFeedbackEmailTemplate struct {
	email.BaseEmailTemplate
}

func SendUserFeedbackEmailForRecipient(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient) (*email.ID, error) {
	emailRecordID := email.NewEmailRecordID()
	baseTemplate, err := createBaseTemplate(emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	body, err := createUserFeedbackEmailBody(userFeedbackEmailTemplate{
		BaseEmailTemplate: *baseTemplate,
	})
	if err != nil {
		return nil, err
	}
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient:       recipient.EmailAddress,
		HTMLBody:        *body,
		Subject:         "What can Babblegraph do better?",
		EmailSenderName: ptr.String(feedbackEmailSenderName),
	})
	if err != nil {
		return nil, err
	}
	if err := email.InsertEmailRecord(tx, emailRecordID, *sesMessageID, recipient.UserID, email.EmailTypeUserFeedback); err != nil {
		return nil, err
	}
	return &emailRecordID, nil
}

func createUserFeedbackEmailBody(templateData userFeedbackEmailTemplate) (*string, error) {
	templateFile, err := getPathForTemplateFile(userFeedbackEmailTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(userFeedbackEmailTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, templateData); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
