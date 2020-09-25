package queuedefs

import (
	"babblegraph/worker/documents"
	"babblegraph/worker/indexer"
	"babblegraph/worker/storage"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameIndexQueue queueTopicName = "index-queue-topic"

type indexQueue struct{}

func (i indexQueue) GetTopicName() string {
	return queueTopicNameIndexQueue.Str()
}

type indexQueueMessage struct {
	Filename   storage.FileIdentifier `json:'filename'`
	DocumentID documents.DocumentID   `json:"document_id"`
}

func (i indexQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m indexQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	return indexer.IndexTermsForFile(tx, m.DocumentID, m.Filename)
}

func publishMessageToIndexQueue(docID documents.DocumentID, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameIndexQueue.Str(), indexQueueMessage{
		Filename:   filename,
		DocumentID: docID,
	})
}
