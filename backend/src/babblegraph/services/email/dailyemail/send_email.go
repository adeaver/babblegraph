package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/util/email"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
)

func sendDailyEmailsForDocuments(cl *email.Client, recipient string, docs []documents.Document) error {
	var links []email.DailyEmailLink
	for _, doc := range docs {
		title := doc.Metadata.Title
		imageURL := doc.Metadata.Image
		description := doc.Metadata.Description
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
