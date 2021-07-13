package email

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlemma"
	"babblegraph/model/users"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"babblegraph/util/text"
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
	LemmaReinforcementSpotlight *dailyEmailLemmaReinforcementSpotlight
	Categories                  []dailyEmailCategory
	SetTopicsLink               *string
	ReinforcementLink           string
}

type dailyEmailLemmaReinforcementSpotlight struct {
	LemmaText       string
	Document        dailyEmailLink
	PreferencesLink string
}

type dailyEmailCategory struct {
	CategoryName *string
	Links        []dailyEmailLink
}

type dailyEmailLink struct {
	ImageURL         *string
	Title            *string
	Description      *string
	URL              string
	PaywallReportURL string
}

type CategorizedDocuments struct {
	Topic     *contenttopics.ContentTopic
	Documents []documents.Document
}

type LemmaReinforcementSpotlight struct {
	Lemma    wordsmith.Lemma
	Document documents.Document
}

type DailyEmailInput struct {
	LemmaReinforcementSpotlight *LemmaReinforcementSpotlight
	CategorizedDocuments        []CategorizedDocuments
	HasSetTopics                bool
}

func SendDailyEmailForDocuments(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient, input DailyEmailInput) error {
	/*
	   This function looks backwards but it is not.
	   This all happens in a transaction, so it is either all succesful or all failed.
	   However, sending an email has a side effect - i.e. if it is successful, it does something. It cannot
	   be reversed like all other parts of the transaction, so we want it to be the last thing we do so that way
	   if anything else fails, we can revert without side effects. So the first thing we do is insert the email record
	   and insert all the user documents.
	*/
	emailRecordID := email.NewEmailRecordID()
	if err := email.InsertEmailRecord(tx, emailRecordID, recipient.UserID, email.EmailTypeDaily); err != nil {
		return err
	}
	template, err := createDailyEmailTemplate(tx, emailRecordID, recipient, input)
	if err != nil {
		return err
	}
	emailBody, err := createEmailBody(*template)
	if err != nil {
		return err
	}
	today := time.Now()
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient: recipient.EmailAddress,
		HTMLBody:  *emailBody,
		Subject:   fmt.Sprintf("Babblegraph Daily Links - %s %d, %d", today.Month().String(), today.Day(), today.Year()),
	})
	if err != nil {
		return err
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
	return nil
}

func createDailyEmailTemplate(tx *sqlx.Tx, emailRecordID email.ID, recipient email.Recipient, input DailyEmailInput) (*dailyEmailTemplate, error) {
	baseTemplate, err := createBaseTemplate(tx, emailRecordID, recipient)
	if err != nil {
		return nil, err
	}
	alreadyHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, recipient.UserID)
	if err != nil {
		return nil, err
	}
	categories := createEmailCategories(tx, recipient.UserID, emailRecordID, input.CategorizedDocuments)
	var setTopicsLink *string
	if !input.HasSetTopics {
		if alreadyHasAccount {
			setTopicsLink = ptr.String(routes.MakeLoginLinkWithContentTopicsRedirect())
		} else {
			setTopicsLink, err = routes.MakeSetTopicsLink(recipient.UserID)
			if err != nil {
				return nil, err
			}
		}
	}
	var reinforcementLink *string
	if alreadyHasAccount {
		reinforcementLink = ptr.String(routes.MakeLoginLinkWithReinforcementRedirect())
	} else {
		reinforcementLink, err = routes.MakeWordReinforcementLink(recipient.UserID)
		if err != nil {
			return nil, err
		}
	}
	var reinforcementSpotlight *dailyEmailLemmaReinforcementSpotlight
	if input.LemmaReinforcementSpotlight != nil {
		reinforcementSpotlight = &dailyEmailLemmaReinforcementSpotlight{
			LemmaText: input.LemmaReinforcementSpotlight.Lemma.LemmaText,
			Document: createLinksForDocuments(tx, recipient.UserID, emailRecordID, []documents.Document{
				input.LemmaReinforcementSpotlight.Document,
			})[0],
			PreferencesLink: routes.MakeLoginLinkWithNewsletterPreferencesRedirect(),
		}
		if err := userlemma.UpsertLemmaReinforcementSpotlightRecord(tx, userlemma.UpsertLemmaReinforcementSpotlightRecordInput{
			UserID:       recipient.UserID,
			LemmaID:      input.LemmaReinforcementSpotlight.Lemma.ID,
			LanguageCode: input.LemmaReinforcementSpotlight.Lemma.Language,
		}); err != nil {
			return nil, err
		}
	}
	return &dailyEmailTemplate{
		BaseEmailTemplate:           *baseTemplate,
		LemmaReinforcementSpotlight: reinforcementSpotlight,
		Categories:                  categories,
		SetTopicsLink:               setTopicsLink,
		ReinforcementLink:           *reinforcementLink,
	}, nil
}

func createEmailCategories(tx *sqlx.Tx, userID users.UserID, emailRecordID email.ID, categorizedDocuments []CategorizedDocuments) []dailyEmailCategory {
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
			Links:        createLinksForDocuments(tx, userID, emailRecordID, categorized.Documents),
		})
	}
	return out
}

func createLinksForDocuments(tx *sqlx.Tx, userID users.UserID, emailRecordID email.ID, documents []documents.Document) []dailyEmailLink {
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
		userDocumentID, err := userdocuments.InsertDocumentForUserAndReturnID(tx, userID, emailRecordID, doc)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		link, err := routes.MakeArticleLink(*userDocumentID)
		if err != nil {
			log.Println(fmt.Sprintf("Got error making article link: %s", err.Error()))
			continue
		}
		paywallReportLink, err := routes.MakePaywallReportLink(*userDocumentID)
		if err != nil {
			log.Println(fmt.Sprintf("Got error making paywall report link: %s", err.Error()))
			continue
		}
		links = append(links, dailyEmailLink{
			ImageURL:         imageURL,
			Title:            title,
			Description:      description,
			URL:              *link,
			PaywallReportURL: *paywallReportLink,
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
