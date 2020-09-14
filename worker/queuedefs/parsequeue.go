package queuedefs

import (
	"babblegraph/worker/htmlparse"
	"babblegraph/worker/storage"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameParseQueue queueTopicName = "parse-queue-topic"

type parseQueue struct{}

func (p parseQueue) GetTopicName() string {
	return queueTopicNameParseQueue.Str()
}

type parseQueueMessage struct {
	Filename storage.FileIdentifier `json:"filename"`
	URL      string                 `json:"url"`
}

func (p parseQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m parseQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	id, links, err := htmlparse.ParseAndStoreFileText(m.Filename)
	if err != nil {
		return err
	}
	return publishMessageToNormalizeTextQueue(m.URL, links, *id)
}

func publishMessageToParseQueue(url string, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(queueTopicNameParseQueue.Str(), parseQueueMessage{
		Filename: filename,
		URL:      url,
	})
}
