package queuedefs

import (
	"encoding/json"
	"fmt"
	"log"

	"babblegraph/model/documents"
	"babblegraph/services/worker/lemmatize"
	"babblegraph/util/queue"
	"babblegraph/util/storage"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const queueTopicNameLemmatizeQueue queueTopicName = "lemmatize-topic"

type lemmatizeQueue struct{}

func (l lemmatizeQueue) GetTopicName() string {
	return queueTopicNameLemmatizeQueue.Str()
}

type lemmatizeQueueMessage struct {
	DocumentID   documents.DocumentID   `json:"document_id"`
	Filename     storage.FileIdentifier `json:"file_name"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
}

func (l lemmatizeQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m lemmatizeQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	id, err := lemmatize.LemmatizeWordsForFile(m.Filename, m.LanguageCode)
	if err != nil {
		return err
	}
	return publishMessageToIndexQueue(m.DocumentID, m.LanguageCode, *id)
}

func publishMessageToLemmatizeQueue(filename storage.FileIdentifier, docID documents.DocumentID, languageCode wordsmith.LanguageCode) error {
	return queue.PublishMessageToQueueByName(queueTopicNameLemmatizeQueue.Str(), lemmatizeQueueMessage{
		Filename:     filename,
		DocumentID:   docID,
		LanguageCode: languageCode,
	})
}
