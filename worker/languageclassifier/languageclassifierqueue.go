package languageclassifier

import (
	"encoding/json"
	"fmt"
	"log"

	"babblegraph/worker/storage"

	"github.com/adeaver/babblegraph/lib/queue"

	"github.com/jmoiron/sqlx"
)

type queueMessage struct {
	Filename storage.FileIdentifier `json:"filename"`
	URL      string                 `json:"url"`
	Links    []string               `json:"links"`
}

var LanguageClassifierQueueImpl languageClassifierQueue = languageClassifierQueue{}

type languageClassifierQueue struct{}

const languageClassifierTopicName string = "language-classifier-topic"

func (l languageClassifierQueue) GetTopicName() string {
	return languageClassifierTopicName
}

func (l languageClassifierQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m queueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	textBytes, err := storage.ReadFile(m.Filename)
	if err != nil {
		return err
	}
	language, err := classify(string(textBytes))
	if err != nil {
		return err
	}
	if language == nil {
		log.Println("do not recognize language: skipping")
		return nil
	}
	log.Println(fmt.Sprintf("Got language: %s", (*language).Str()))
	return nil
}

func PublishMessageToLanguageClassifierQueue(url string, links []string, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(languageClassifierTopicName, queueMessage{
		Filename: filename,
		URL:      url,
		Links:    links,
	})
}
