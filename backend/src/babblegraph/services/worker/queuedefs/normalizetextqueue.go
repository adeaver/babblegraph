package queuedefs

import (
	"babblegraph/services/worker/normalizetext"
	"babblegraph/util/queue"
	"babblegraph/util/storage"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const queueTopicNameNormalizeTextQueue queueTopicName = "normalize-text-topic"

type normalizeTextQueue struct{}

func (n normalizeTextQueue) GetTopicName() string {
	return queueTopicNameNormalizeTextQueue.Str()
}

type normalizeTextQueueMessage struct {
	Filename     storage.FileIdentifier `json:"filename"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	URL          string                 `json:"url"`
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
	if err := publishMessageToReadabilityQueue(*id, m.URL, m.LanguageCode); err != nil {
		return err
	}
	if err := storage.DeleteFile(m.Filename); err != nil {
		log.Println(fmt.Sprintf("Error deleting file %s, marking message as done", string(m.Filename)))
	}
	return nil
}

func publishMessageToNormalizeTextQueue(url string, languageCode wordsmith.LanguageCode, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameNormalizeTextQueue.Str(), normalizeTextQueueMessage{
		Filename:     filename,
		LanguageCode: languageCode,
		URL:          url,
	})
}
