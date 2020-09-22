package queuedefs

import (
	"babblegraph/worker/normalizetext"
	"babblegraph/worker/storage"
	"babblegraph/worker/wordsmith"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameNormalizeTextQueue queueTopicName = "normalize-text-topic"

type normalizeTextQueue struct{}

func (n normalizeTextQueue) GetTopicName() string {
	return queueTopicNameNormalizeTextQueue.Str()
}

type normalizeTextQueueMessage struct {
	Filename     storage.FileIdentifier `json:"filename"`
	URL          string                 `json:"url"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
}

func (n normalizeTextQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m normalizeTextQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	id, err := normalizetext.NormalizeAndStoreTextForFilename(m.Filename)
	if err != nil {
		return err
	}
	return nil
}

func publishMessageToNormalizeTextQueue(url string, languageCode wordsmith.LanguageCode, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameNormalizeTextQueue.Str(), normalizeTextQueueMessage{
		Filename:     filename,
		URL:          url,
		LanguageCode: languageCode,
	})
}
