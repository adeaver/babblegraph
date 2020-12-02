package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/util/deref"
	"babblegraph/util/email"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
)

func sendDailyEmailsForDocuments(cl *email.Client, recipient string, docs []documents.Document) error {
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
		log.Println(fmt.Sprintf("Sending URL %s to recipient %s", doc.URL, recipient))
		links = append(links, email.DailyEmailLink{
			ImageURL:    imageURL,
			Title:       title,
			Description: description,
			URL:         doc.URL,
		})
	}
	return cl.SendDailyEmailForLinks(recipient, links)
}

// TODO: remove this function when documents
// have been reindexed
func isNotEmpty(s *string) bool {
	return len(deref.String(s, "")) > 0
}
