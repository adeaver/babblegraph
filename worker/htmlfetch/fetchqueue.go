package htmlfetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"babblegraph/worker/htmlparse"
	"babblegraph/worker/storage"

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
	filename, err := fetchAndStoreHTMLForURL(m.URL)
	if err != nil {
		return err
	}
	return htmlparse.PublishFilenameToParseQueue(m.URL, *filename)
}

func fetchAndStoreHTMLForURL(url string) (*storage.FileIdentifier, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got status code for website: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return storage.WriteFile("html", string(data))
}
