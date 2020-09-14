package queuedefs

import (
	"babblegraph/worker/languageclassifier"
	"babblegraph/worker/storage"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameLanguageClassifierQueue queueTopicName = "language-classifier-topic"

type languageClassifierQueue struct{}

func (l languageClassifierQueue) GetTopicName() string {
	return queueTopicNameLanguageClassifierQueue.Str()
}

type languageClassifierQueueMessage struct {
	Filename storage.FileIdentifier `json:"filename"`
	URL      string                 `json:"url"`
	Links    []string               `json:"links"`
}

func (l languageClassifierQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m languageClassifierQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	language, err := languageclassifier.ClassifyLanguageForFile(m.Filename)
	if err != nil {
		return err
	}
	if language == nil {
		log.Println("do not recognize language: skipping")
		return nil
	}
	log.Println(fmt.Sprintf("Got language %s for %s", language.Str(), m.URL))
	if err := publishMessageToLinkHandlerQueue(m.Links); err != nil {
		return err
	}
	return nil
}

func publishMessageToLanguageClassifierQueue(url string, links []string, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameLanguageClassifierQueue.Str(), languageClassifierQueueMessage{
		Filename: filename,
		URL:      url,
		Links:    links,
	})
}
