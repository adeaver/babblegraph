package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/util/email"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
)

func sendDailyEmailsForDocuments(cl *email.Client, recipient string, docs []documents.Document) error {
	var links []email.DailyEmailLink
	for _, doc := range docs {
		var title, imageURL *string
		if doc.Metadata != nil {
			if len(doc.Metadata.Title) > 0 {
				title = ptr.String(doc.Metadata.Title)
			}
			if len(doc.Metadata.Image) > 0 && urlparser.IsValidURL(doc.Metadata.Image) {
				imageURL = ptr.String(doc.Metadata.Image)
			}
		}
		if !urlparser.IsValidURL(doc.URL) {
			continue
		}
		log.Println(fmt.Sprintf("Sending URL %s to recipient %s", doc.URL, recipient))
		links = append(links, email.DailyEmailLink{
			ImageURL: imageURL,
			Title:    title,
			URL:      doc.URL,
		})
	}
	return cl.SendDailyEmailForLinks(recipient, links)
}
