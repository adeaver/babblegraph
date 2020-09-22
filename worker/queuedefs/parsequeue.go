package queuedefs

import (
	"babblegraph/worker/htmlparse"
	"babblegraph/worker/storage"
	"babblegraph/worker/wordsmith"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
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
	// TODO: this is probably a good place to store a reference to a document
	// In the future, I want to be able to pull documents that have other languages
	// to train classifiers - for instance
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
