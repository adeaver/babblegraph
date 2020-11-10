package queuedefs

import (
	"babblegraph/model/documents"
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
	// A null document version corresponds to version 1
	DocumentVersion  *documents.Version `json:"document_version,omitempty"`
	DocumentType     *documents.Type    `json:"document_type,omitempty"`
	DocumentMetadata map[string]string  `json:"document_metadata,omitempty"`
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
	if err := publishMessageToReadabilityQueue(readabilityQueueMessage{
		Filename:         *id,
		LanguageCode:     m.LanguageCode,
		URL:              m.URL,
		DocumentVersion:  m.DocumentVersion,
		DocumentType:     m.DocumentType,
		DocumentMetadata: m.DocumentMetadata,
	}); err != nil {
		return err
	}
	if err := storage.DeleteFile(m.Filename); err != nil {
		log.Println(fmt.Sprintf("Error deleting file %s, marking message as done", string(m.Filename)))
	}
	return nil
}

func publishMessageToNormalizeTextQueue(msg normalizeTextQueueMessage) error {
	return queue.PublishMessageToQueueByName(queueTopicNameNormalizeTextQueue.Str(), msg)
}
