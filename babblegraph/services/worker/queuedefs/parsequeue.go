package queuedefs

import (
	"encoding/json"
	"fmt"
	"log"

	"babblegraph/model/htmlpages"
	"babblegraph/services/worker/htmlparse"
	"babblegraph/util/queue"
	"babblegraph/util/storage"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const queueTopicNameParseQueue queueTopicName = "parse-queue-topic"

type parseQueue struct{}

func (p parseQueue) GetTopicName() string {
	return queueTopicNameParseQueue.Str()
}

type parseQueueMessage struct {
	Filename storage.FileIdentifier `json:"filename"`
	URL      string                 `json:"url"`
}

func (p parseQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m parseQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	parsedDoc, err := htmlparse.ParseAndStoreFileText(m.Filename)
	if err != nil {
		return err
	}
	docID, err := htmlpages.InsertHTMLPage(tx, htmlpages.InsertHTMLPageInput{
		URL:      m.URL,
		Language: parsedDoc.LanguageValue,
		Metadata: parsedDoc.Metadata,
	})
	if err != nil {
		return err
	}
	if docID == nil {
		log.Println("Did not insert document id... probably duplicate. skipping")
		return nil
	}
	if parsedDoc.LanguageValue == nil {
		log.Println("HTML document had no language code, marking complete...")
		return nil
	}
	languageCode := wordsmith.LookupLanguageCodeForLanguageLabel(*parsedDoc.LanguageValue)
	if languageCode == nil {
		log.Println(fmt.Sprintf("Unsupported language label: %s, marking complete", *parsedDoc.LanguageValue))
		return nil
	}
	if err := publishMessageToLinkHandlerQueue(parsedDoc.Links); err != nil {
		return err
	}
	return publishMessageToNormalizeTextQueue(m.URL, *languageCode, parsedDoc.BodyTextFilename)
}

func publishMessageToParseQueue(url string, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameParseQueue.Str(), parseQueueMessage{
		Filename: filename,
		URL:      url,
	})
}
