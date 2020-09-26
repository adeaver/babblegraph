package queuedefs

import (
	"encoding/json"
	"fmt"
	"log"

	"babblegraph/util/queue"
	"babblegraph/worker/htmlfetch"

	"github.com/jmoiron/sqlx"
)

const queueTopicNameFetchQueue queueTopicName = "fetch-queue-topic"

type fetchQueue struct{}

func (f fetchQueue) GetTopicName() string {
	return queueTopicNameFetchQueue.Str()
}

type fetchQueueMessage struct {
	URL string `json:"url"`
}

func (f fetchQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m fetchQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	filename, err := htmlfetch.FetchAndStoreHTMLForURL(m.URL)
	if err != nil {
		log.Println(fmt.Sprintf("Error on url %s, message: %s... skipping", m.URL, err.Error()))
		return nil
	}
	return publishMessageToParseQueue(m.URL, *filename)
}

func publishMessageToFetchQueue(u string) error {
	return queue.PublishMessageToQueueByName(queueTopicNameFetchQueue.Str(), fetchQueueMessage{
		URL: u,
	})
}
