package email

import (
	"babblegraph/model/email"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
)

const (
	genericEmailTemplateFilename string = "generic_email_with_optional_action_link_template.html"
)

type genericEmailWithOptionalActionTemplate struct {
	email.BaseEmailTemplate
	EmailTitle       string
	PreheaderText    string
	BeforeParagraphs []string
	Action           *GenericEmailAction
	AfterParagraphs  []string
}

type GenericEmailAction struct {
	Link       string
	ButtonText string
}

type SendGenericEmailWithOptionalActionForRecipientInput struct {
	EmailType          email.EmailType
	Recipient          email.Recipient
	PreheaderText      string
	BeforeParagraphs   []string
	AfterParagraphs    []string
	GenericEmailAction *GenericEmailAction
	FromEmailName      *string
	Subject            string
	EmailTitle         string
}

func SendGenericEmailWithOptionalActionForRecipient(tx *sqlx.Tx, cl *ses.Client, input SendGenericEmailWithOptionalActionForRecipientInput) (*email.ID, error) {
	switch {
	case len(input.BeforeParagraphs) == 0:
		return nil, fmt.Errorf("Must include before paragraphs")
	case input.GenericEmailAction != nil && len(input.GenericEmailAction.Link) == 0:
		return nil, fmt.Errorf("Cannot have empty link in email")
	case input.GenericEmailAction != nil && len(input.GenericEmailAction.ButtonText) == 0:
		return nil, fmt.Errorf("Cannot have empty button text in email")
	case len(input.Subject) == 0:
		return nil, fmt.Errorf("Cannot have empty subject for email")
	case len(input.EmailTitle) == 0:
		return nil, fmt.Errorf("Cannot have empty title for email")
	}
	emailRecordID := email.NewEmailRecordID()
	template, err := createGenericEmailWithOptionalActionTemplate(emailRecordID, input)
	if err != nil {
		return nil, err
	}
	body, err := createGenericEmailBody(*template)
	if err != nil {
		return nil, err
	}
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient:       input.Recipient.EmailAddress,
		HTMLBody:        *body,
		Subject:         input.Subject,
		EmailSenderName: input.FromEmailName,
	})
	if err != nil {
		return nil, err
	}
	if err := email.InsertEmailRecord(tx, emailRecordID, *sesMessageID, input.Recipient.UserID, input.EmailType); err != nil {
		return nil, err
	}
	return &emailRecordID, nil
}

func createGenericEmailWithOptionalActionTemplate(emailRecordID email.ID, input SendGenericEmailWithOptionalActionForRecipientInput) (*genericEmailWithOptionalActionTemplate, error) {
	baseTemplate, err := createBaseTemplate(emailRecordID, input.Recipient)
	if err != nil {
		return nil, err
	}
	return &genericEmailWithOptionalActionTemplate{
		BaseEmailTemplate: *baseTemplate,
		EmailTitle:        input.EmailTitle,
		PreheaderText:     input.PreheaderText,
		BeforeParagraphs:  input.BeforeParagraphs,
		AfterParagraphs:   input.AfterParagraphs,
		Action:            input.GenericEmailAction,
	}, nil
}

func createGenericEmailBody(templateData genericEmailWithOptionalActionTemplate) (*string, error) {
	templateFile, err := getPathForTemplateFile(genericEmailTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(genericEmailTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, templateData); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
