package dailyemail

import (
	"babblegraph/model/documents"
	email_model "babblegraph/model/email"
	"babblegraph/model/userdocuments"
	"babblegraph/util/deref"
	"babblegraph/util/email"
	"babblegraph/util/urlparser"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func sendDailyEmailsForDocuments(tx *sqlx.Tx, cl *email.Client, recipient email.Recipient, docs []documents.Document) error {
	var links []email.DailyEmailLink
	for _, doc := range docs {
		var title, imageURL, description *string
		if isNotEmpty(doc.Metadata.Title) {
			title = doc.Metadata.Title
		}
		if isNotEmpty(doc.Metadata.Image) {
			imageURL = doc.Metadata.Image
		}
		if isNotEmpty(doc.Metadata.Description) {
			log.Println(fmt.Sprintf("Found description %s for URL %s", *doc.Metadata.Description, doc.URL))
			description = doc.Metadata.Description
		}
		if !urlparser.IsValidURL(doc.URL) {
			continue
		}
		log.Println(fmt.Sprintf("Sending URL %s to recipient %s", doc.URL, recipient.EmailAddress))
		links = append(links, email.DailyEmailLink{
			ImageURL:    imageURL,
			Title:       title,
			Description: description,
			URL:         doc.URL,
		})
	}
	sesMessageID, err := cl.SendDailyEmailForLinks(recipient, links)
	if err != nil {
		return err
	}
	return insertEmailRecordAndUpdateUserDocuments(tx, recipient, *sesMessageID, docs)
}

// TODO: remove this function when documents
// have been reindexed
func isNotEmpty(s *string) bool {
	return len(deref.String(s, "")) > 0
}

func insertEmailRecordAndUpdateUserDocuments(tx *sqlx.Tx, recipient email.Recipient, sesMessageID string, documents []documents.Document) error {
	emailRecordID, err := email_model.CreateEmailRecord(tx, *sesMessageID, recipient.UserID)
	if err != nil {
		return err
	}
	return userdocuments.InsertDocumentIDsForUser(tx, recipient.UserID, *emailRecordID, docs)
}
