package languageclassifier

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"

	"github.com/jmoiron/sqlx"
)

type queueMessage struct {
	Filename storage.FileIdentifer `json:"filename"`
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
	log.Println(fmt.Sprintf("Got language: %s", *language.Str()))
	return nil
}

func PublishMessageToLanguageClassifierQueue(filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(languageClassifierTopicName, queueMessage{
		Filename: filename,
	})
}
