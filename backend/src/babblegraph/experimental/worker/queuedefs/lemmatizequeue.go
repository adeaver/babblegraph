package queuedefs

import (
	"encoding/json"
	"fmt"
	"log"

	"babblegraph/model/documents"
	"babblegraph/services/worker/lemmatize"
	"babblegraph/util/math/decimal"
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
	URL              string                 `json:"url"`
	Filename         storage.FileIdentifier `json:"file_name"`
	LanguageCode     wordsmith.LanguageCode `json:"language_code"`
	ReadabilityScore decimal.Number         `json:"readability_score"`
	// A null document version corresponds to version 1
	DocumentVersion  *documents.Version `json:"document_version,omitempty"`
	DocumentType     *documents.Type    `json:"document_type,omitempty"`
	DocumentMetadata map[string]string  `json:"document_metadata,omitempty"`
}

func (l lemmatizeQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m lemmatizeQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	lemmaText, err := lemmatize.LemmatizeWordsForFile(m.Filename, m.LanguageCode)
	if err != nil {
		return err
	}
	documentVersion := documents.Version1
	if m.DocumentVersion != nil {
		documentVersion = *m.DocumentVersion
	}
	documentMetadata := documents.ExtractMetadataFromMap(m.DocumentMetadata)
	docID, err := documents.AssignIDAndIndexDocument(&documents.Document{
		URL:              m.URL,
		Version:          documentVersion,
		ReadabilityScore: m.ReadabilityScore.ToInt64Rounded(),
		LanguageCode:     m.LanguageCode,
		LemmatizedBody:   *lemmaText,
		DocumentType:     m.DocumentType,
		Metadata:         &documentMetadata,
	})
	if err != nil {
		return err
	}
	log.Println("Indexed doc with ID %s", *docID)
	if err := storage.DeleteFile(m.Filename); err != nil {
		log.Println(fmt.Sprintf("Error deleting file %s, marking message as done", string(m.Filename)))
	}
	return nil
}

func publishMessageToLemmatizeQueue(msg lemmatizeQueueMessage) error {
	return queue.PublishMessageToQueueByName(queueTopicNameLemmatizeQueue.Str(), msg)
}
