package queuedefs

import (
	"babblegraph/worker/lemmatize"
	"babblegraph/worker/storage"
	"babblegraph/worker/wordsmith"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameLemmatizeQueue queueTopicName = "lemmatize-topic"

type lemmatizeQueue struct{}

func (l lemmatizeQueue) GetTopicName() string {
	return queueTopicNameLemmatizeQueue.Str()
}

type lemmatizeQueueMessage struct {
	Filename     storage.FileIdentifier `json:"file_name"`
	URL          string                 `json:"url"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
}

func (l lemmatizeQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m lemmatizeQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	_, err := lemmatize.LemmatizeWordsForFile(m.Filename, m.LanguageCode)
	if err != nil {
		return err
	}
	return nil
}

func publishMessageToLemmatizeQueue(filename storage.FileIdentifier, url string, languageCode wordsmith.LanguageCode) error {
	return queue.PublishMessageToQueueByName(queueTopicNameLemmatizeQueue.Str(), lemmatizeQueueMessage{
		Filename:     filename,
		URL:          url,
		LanguageCode: languageCode,
	})
}
