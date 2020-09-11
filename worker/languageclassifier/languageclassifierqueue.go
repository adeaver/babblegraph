package languageclassifier

import (
	"github.com/adeaver/babblegraph/lib/queue"

	"github.com/jmoiron/sqlx"
)

type queueMessage struct {
	Filename storage.FileIdentifer `json:"filename"`
}

var LanguageClassifierQueueImpl languageClassifierQueue = languageClassifierQueue{}

type languageClassifierQueue struct{}

const languageClassifierTopicName string = "language-classifier-topic"

func (l languageClassifierQueue) GetTopicName() string {
	return languageClassifierTopicName
}

func (l languageClassifierQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	return nil
}
