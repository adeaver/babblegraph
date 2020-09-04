package queue

import (
	"babblegraph/worker/database"
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

func StartQueue() error {
	for _, queue := range registry.queuesByName {
		go readQueue(*queue)
	}
	return nil
}

const (
	QueueWaitTimeMilliseconds time.Duration = 1500
)

func readQueue(queue Queue) {
	for {
		var msg *dbMessage
		for msg == nil {
			time.Sleep(QueueWaitTimeMilliseconds * time.Millisecond)
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				msg, err = getMostRecentMessageForQueue(tx, queue.GetTopicName())
				return err
			}); err != nil {
				// TODO: figure out errors
				panic(err.Error())
			}
			log.Println("Got message with ID %s for queue %s", msg.ID, queue.GetTopicName())
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := queue.ProcessMessage(tx, msg.ToMessage()); err != nil {
				return err
			}
			return clearMessage(tx, msg.ID)
		}); err != nil {
			// TODO: figure out errors
			panic(err.Error())
		}
	}
}
