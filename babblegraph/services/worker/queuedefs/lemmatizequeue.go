package queuedefs

import (
	"encoding/json"
	"fmt"
	"log"

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
}

func (l lemmatizeQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m lemmatizeQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	// TODO: insert Elasticsearch here
	_, err = lemmatize.LemmatizeWordsForFile(m.Filename, m.LanguageCode)
	return err
}

func publishMessageToLemmatizeQueue(filename storage.FileIdentifier, url string, languageCode wordsmith.LanguageCode, readabilityScore decimal.Number) error {
	return queue.PublishMessageToQueueByName(queueTopicNameLemmatizeQueue.Str(), lemmatizeQueueMessage{
		URL:              url,
		Filename:         filename,
		LanguageCode:     languageCode,
		ReadabilityScore: readabilityScore,
	})
}
