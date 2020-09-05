package htmlparse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

type queueMessage struct {
	Filename string `json:"filename"`
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
	htmlBytes, err := ioutil.ReadFile(m.Filename)
	if err != nil {
		return err
	}
	text, _, err := getTextAndLinksForHTML(string(htmlBytes))
	if err != nil {
		return err
	}
	log.Println(*text)
	return nil
}

func PublishFilenameToParseQueue(filename string) error {
	return queue.PublishMessageToQueueByName(parseQueueTopicName, queueMessage{
		Filename: filename,
	})
}
