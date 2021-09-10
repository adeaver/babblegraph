package email

import (
	"babblegraph/model/email"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"
	"log"
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
	if err := email.InsertEmailRecord(tx, emailRecordID, recipient.UserID, email.EmailTypeUserFeedback); err != nil {
		return nil, err
	}
	if err := email.SetEmailRecordSentAtTime(tx, emailRecordID); err != nil {
		return nil, err
	}
	baseTemplate, err := createBaseTemplate(tx, emailRecordID, recipient)
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
