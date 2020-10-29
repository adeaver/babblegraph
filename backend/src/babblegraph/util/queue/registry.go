package queue

import (
	"babblegraph/util/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type queueRegistry struct {
	queuesByName map[string]*Queue
}

var registry queueRegistry = queueRegistry{
	queuesByName: make(map[string]*Queue),
}

func RegisterQueue(queue Queue) error {
	if _, exists := registry.queuesByName[queue.GetTopicName()]; exists {
		return fmt.Errorf("queue with name %s already exists", queue.GetTopicName())
	}
	registry.queuesByName[queue.GetTopicName()] = &queue
	return nil
}

func PublishMessageToQueueByName(topicName string, message interface{}) error {
	if _, exists := registry.queuesByName[topicName]; !exists {
		return fmt.Errorf("queue with name %s does not exist", topicName)
	}
	return database.WithTx(func(tx *sqlx.Tx) error {
		jsonBody, err := json.Marshal(message)
		if err != nil {
			return err
		}
		return saveMessage(tx, dbMessage{
			Topic: topicName,
			Body:  string(jsonBody),
		})
	})
}

func StartQueue(errs chan error) error {
	for _, queue := range registry.queuesByName {
		go readQueue(*queue, errs)
	}
	return nil
}

const (
	QueueWaitTimeMilliseconds time.Duration = 1500
)

func readQueue(queue Queue, errs chan error) {
	for {
		var msg *dbMessage
		for msg == nil {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				msg, err = getMostRecentMessageForQueue(tx, queue.GetTopicName())
				return err
			}); err != nil {
				if err != sql.ErrNoRows {
					log.Println(fmt.Sprintf("Error: %s", err.Error()))
					close(errs)
				}
				log.Println(fmt.Sprintf("no messages in queue %s... sleeping", queue.GetTopicName()))
				time.Sleep(QueueWaitTimeMilliseconds * time.Millisecond)
			}
		}
		log.Println(fmt.Sprintf("Got message with ID %s for queue %s", msg.ID, queue.GetTopicName()))
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := queue.ProcessMessage(tx, msg.ToMessage()); err != nil {
				return err
			}
			return clearMessage(tx, msg.ID)
		}); err != nil {
			log.Println(fmt.Sprintf("Error: %s", err.Error()))
			close(errs)
		}
	}
}
