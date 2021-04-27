package email

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"babblegraph/util/text"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
)

const dailyEmailTemplateFilename = "daily_email_template.html"

type dailyEmailTemplate struct {
	email.BaseEmailTemplate
	Categories        []dailyEmailCategory
	SetTopicsLink     *string
	ReinforcementLink string
}

type dailyEmailCategory struct {
	CategoryName *string
	Links        []dailyEmailLink
}

type dailyEmailLink struct {
	ImageURL    *string
	Title       *string
	Description *string
	URL         string
}

type CategorizedDocuments struct {
	Topic     *contenttopics.ContentTopic
	Documents []documents.Document
}

type DailyEmailInput struct {
	CategorizedDocuments []CategorizedDocuments
	HasSetTopics         bool
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
	categories := createEmailCategories(input.CategorizedDocuments)
	var setTopicsLink *string
	if !input.HasSetTopics {
		setTopicsLink, err = routes.MakeSetTopicsLink(recipient.UserID)
		if err != nil {
			return nil, err
		}
	}
	reinforcementLink, err := routes.MakeWordReinforcementLink(recipient.UserID)
	if err != nil {
		return nil, err
	}
	return &dailyEmailTemplate{
		BaseEmailTemplate: *baseTemplate,
		Categories:        categories,
		SetTopicsLink:     setTopicsLink,
		ReinforcementLink: *reinforcementLink,
	}, nil
}

func createEmailCategories(categorizedDocuments []CategorizedDocuments) []dailyEmailCategory {
	var out []dailyEmailCategory
	// TODO(other-languages): don't hardcode this
	languageCode := wordsmith.LanguageCodeSpanish
	for _, categorized := range categorizedDocuments {
		var contentTopicCategory *string
		switch {
		case categorized.Topic != nil:
			displayName, err := contenttopics.ContentTopicNameToDisplayName(*categorized.Topic)
			if err != nil {
				log.Println(fmt.Sprintf("Got error converting content topic %s: %s", categorized.Topic.Str(), err.Error()))
			} else {
				// TODO: don't hardcode this
				contentTopicCategory = ptr.String(text.ToTitleCaseForLanguage(displayName.Str(), languageCode))
			}
		case categorized.Topic == nil && len(categorizedDocuments) > 1:
			displayName := contenttopics.GenericCategoryNameForLanguage(languageCode)
			contentTopicCategory = ptr.String(text.ToTitleCaseForLanguage(displayName.Str(), wordsmith.LanguageCodeSpanish))
		default:
			// no-op
		}
		out = append(out, dailyEmailCategory{
			CategoryName: contentTopicCategory,
			Links:        createLinksForDocuments(categorized.Documents),
		})
	}
	return out
}

func createLinksForDocuments(documents []documents.Document) []dailyEmailLink {
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
		var url *string
		switch {
		case isNotEmpty(doc.Metadata.URL) && urlparser.IsValidURL(*doc.Metadata.URL):
			url = doc.Metadata.URL
		case urlparser.IsValidURL(doc.URL):
			url = ptr.String(doc.URL)
		default:
			continue
		}
		urlWithProtocol, err := urlparser.EnsureProtocol(*url)
		if err != nil {
			log.Println(fmt.Sprintf("Got error ensuring protocol for URL %s: %s", *url, err.Error()))
			continue
		}
		links = append(links, dailyEmailLink{
			ImageURL:    imageURL,
			Title:       title,
			Description: description,
			URL:         *urlWithProtocol,
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
