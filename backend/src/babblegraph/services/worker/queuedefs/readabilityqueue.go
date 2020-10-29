package queuedefs

import (
	"babblegraph/services/worker/readability"
	"babblegraph/util/queue"
	"babblegraph/util/storage"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const queueTopicNameReadabilityQueue queueTopicName = "readability-queue-topic"

type readabilityQueue struct{}

func (r readabilityQueue) GetTopicName() string {
	return queueTopicNameReadabilityQueue.Str()
}

type readabilityQueueMessage struct {
	Filename     storage.FileIdentifier `json:"filename"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	URL          string                 `json:"url"`
}

func (r readabilityQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m readabilityQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	score, err := readability.CalculateReadability(readability.CalculateReadabilityInput{
		Filename:     m.Filename,
		LanguageCode: m.LanguageCode,
	})
	if err != nil {
		return err
	}
	return publishMessageToLemmatizeQueue(m.Filename, m.URL, m.LanguageCode, *score)
}

func publishMessageToReadabilityQueue(filename storage.FileIdentifier, url string, languageCode wordsmith.LanguageCode) error {
	return queue.PublishMessageToQueueByName(queueTopicNameReadabilityQueue.Str(), readabilityQueueMessage{
		URL:          url,
		Filename:     filename,
		LanguageCode: languageCode,
	})
}
