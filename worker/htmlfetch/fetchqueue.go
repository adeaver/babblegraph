package htmlfetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"babblegraph/worker/htmlparse"

	"github.com/adeaver/babblegraph/lib/queue"
	"github.com/google/uuid"
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
	return htmlparse.PublishFilenameToParseQueue(*filename)
}

func fetchAndStoreHTMLForURL(url string) (*string, error) {
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
	// TODO: replace this with storage package
	filename := fmt.Sprintf("/tmp/%s.html", uuid.New())
	return &filename, ioutil.WriteFile(filename, data, 0644)
}
