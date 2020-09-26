package queuedefs

import (
	"babblegraph/worker/documents"
	"babblegraph/worker/normalizetext"
	"babblegraph/worker/storage"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/adeaver/babblegraph/lib/wordsmith"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameNormalizeTextQueue queueTopicName = "normalize-text-topic"

type normalizeTextQueue struct{}

func (n normalizeTextQueue) GetTopicName() string {
	return queueTopicNameNormalizeTextQueue.Str()
}

type normalizeTextQueueMessage struct {
	Filename     storage.FileIdentifier `json:"filename"`
	DocumentID   documents.DocumentID   `json:"document_id"`
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
	return publishMessageToLemmatizeQueue(*id, m.DocumentID, m.LanguageCode)
}

func publishMessageToNormalizeTextQueue(docID documents.DocumentID, languageCode wordsmith.LanguageCode, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameNormalizeTextQueue.Str(), normalizeTextQueueMessage{
		Filename:     filename,
		DocumentID:   docID,
		LanguageCode: languageCode,
	})
}
