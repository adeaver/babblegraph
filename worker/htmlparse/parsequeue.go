package htmlparse

import (
	"babblegraph/worker/normalizetext"
	"babblegraph/worker/storage"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

type queueMessage struct {
	Filename storage.FileIdentifier `json:"filename"`
	URL      string                 `json:"url"`
}

var ParseQueueImpl parseQueue = parseQueue{}

type parseQueue struct{}

const parseQueueTopicName string = "parse-queue-topic"

func (p parseQueue) GetTopicName() string {
	return parseQueueTopicName
}

func (p parseQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m queueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	htmlBytes, err := storage.ReadFile(m.Filename)
	if err != nil {
		return err
	}
	text, links, err := getTextAndLinksForHTML(string(htmlBytes))
	if err != nil {
		return err
	}
	id, err := storage.WriteFile("txt", *text)
	if err != nil {
		return err
	}
	return normalizetext.PublishMessageToNormalizeTextQueue(m.URL, links, *id)
}

func PublishFilenameToParseQueue(url string, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(parseQueueTopicName, queueMessage{
		Filename: filename,
		URL:      url,
	})
}
