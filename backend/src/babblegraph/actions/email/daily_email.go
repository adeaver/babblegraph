package email

import (
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"babblegraph/util/urlparser"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
)

const dailyEmailTemplateFilename = "daily_email_template.html"

type dailyEmailTemplate struct {
	email.BaseEmailTemplate
	Links         []dailyEmailLink
	SetTopicsLink *string
}

type dailyEmailLink struct {
	ImageURL    *string
	Title       *string
	Description *string
	URL         string
}

type DailyEmailInput struct {
	Documents    []documents.Document
	HasSetTopics bool
}

func SendDailyEmailForDocuments(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient, input DailyEmailInput) (*email.ID, error) {
	emailRecordID := email.NewEmailRecordID()
	template, err := createDailyEmailTemplate(emailRecordID, recipient, input)
	if err != nil {
		return nil, err
	}
	emailBody, err := createEmailBody(*template)
	if err != nil {
		return nil, err
	}
	today := time.Now()
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient: recipient.EmailAddress,
		HTMLBody:  *emailBody,
		Subject:   fmt.Sprintf("Babblegraph Daily Links - %s %d, %d", today.Month().String(), today.Day(), today.Year()),
	})
	if err != nil {
		return nil, err
	}
	if err := email.InsertEmailRecord(tx, emailRecordID, *sesMessageID, recipient.UserID, email.EmailTypeDaily); err != nil {
		return nil, err
	}
	return &emailRecordID, nil
}

func createDailyEmailTemplate(emailRecordID email.ID, recipient email.Recipient, input DailyEmailInput) (*dailyEmailTemplate, error) {
	baseTemplate, err := createBaseTemplate(emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	links := createLinksFromDocuments(input.Documents)
	var setTopicsLink *string
	if !input.HasSetTopics {
		setTopicsLink, err = routes.MakeSetTopicsLink(recipient.UserID)
		if err != nil {
			return nil, err
		}
	}
	return &dailyEmailTemplate{
		BaseEmailTemplate: *baseTemplate,
		Links:             links,
		SetTopicsLink:     setTopicsLink,
	}, nil
}

func createLinksFromDocuments(documents []documents.Document) []dailyEmailLink {
	var links []dailyEmailLink
	for _, doc := range documents {
		var title, imageURL, description *string
		if isNotEmpty(doc.Metadata.Title) {
			title = doc.Metadata.Title
		}
		if isNotEmpty(doc.Metadata.Image) {
			imageURL = doc.Metadata.Image
		}
		if isNotEmpty(doc.Metadata.Description) {
			description = doc.Metadata.Description
		}
		if !urlparser.IsValidURL(doc.URL) {
			continue
		}
		links = append(links, dailyEmailLink{
			ImageURL:    imageURL,
			Title:       title,
			Description: description,
			URL:         doc.URL,
		})
	}
	return links
}

// TODO: remove this function when documents
// have been reindexed
func isNotEmpty(s *string) bool {
	return len(deref.String(s, "")) > 0
}

func createEmailBody(emailData dailyEmailTemplate) (*string, error) {
	templateFile, err := getPathForTemplateFile(dailyEmailTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(dailyEmailTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, emailData); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
