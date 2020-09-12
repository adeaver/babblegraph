package normalizetext

import (
	"babblegraph/worker/languageclassifier"
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
	Links    []string               `json:"links"`
}

var NormalizeTextQueueImpl normalizeTextQueue = normalizeTextQueue{}

type normalizeTextQueue struct{}

const normalizeTextTopicName string = "normalize-text-topic"

func (n normalizeTextQueue) GetTopicName() string {
	return normalizeTextTopicName
}

func (n normalizeTextQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m queueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	textBytes, err := storage.ReadFile(m.Filename)
	if err != nil {
		return err
	}
	normalizedTextLines := normalizeText(textBytes)
	id, err := storage.WriteFile("txt", normalizedTextLines)
	if err != nil {
		return err
	}
	return languageclassifier.PublishMessageToLanguageClassifierQueue(m.URL, m.Links, *id)
}

func PublishMessageToNormalizeTextQueue(url string, links []string, filename storage.FileIdentifier) error {
	return queue.PublishMessageToQueueByName(normalizeTextTopicName, queueMessage{
		Filename: filename,
		URL:      url,
		Links:    links,
	})
}
