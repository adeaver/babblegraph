package queuedefs

import (
	"babblegraph/model/documents"
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
	// A null document version corresponds to version 1
	DocumentVersion  *documents.Version `json:"document_version,omitempty"`
	DocumentType     *documents.Type    `json:"document_type,omitempty"`
	DocumentMetadata map[string]string  `json:"document_metadata,omitempty"`
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
	return publishMessageToLemmatizeQueue(lemmatizeQueueMessage{
		Filename:         m.Filename,
		ReadabilityScore: *score,
		LanguageCode:     m.LanguageCode,
		URL:              m.URL,
		DocumentVersion:  m.DocumentVersion,
		DocumentType:     m.DocumentType,
		DocumentMetadata: m.DocumentMetadata,
	})
}

func publishMessageToReadabilityQueue(msg readabilityQueueMessage) error {
	return queue.PublishMessageToQueueByName(queueTopicNameReadabilityQueue.Str(), msg)
}
