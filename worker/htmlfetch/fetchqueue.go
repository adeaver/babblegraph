package htmlfetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/jmoiron/sqlx"
)

type QueueMessage struct {
	URL string `json:"url"`
}

var FetchQueueImpl fetchQueue = fetchQueue{}

type fetchQueue struct{}

const FetchQueueTopicName string = "fetch-queue-topic"

func (f fetchQueue) GetTopicName() string {
	return FetchQueueTopicName
}

func (f fetchQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m QueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	return fetchAndStoreHTMLForURL(m.URL)
}

func fetchAndStoreHTMLForURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got status code for website: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(data))
	return nil
}
