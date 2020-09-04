package main

import (
	"log"

	"babblegraph/worker/database"
	"babblegraph/worker/queue"

	"github.com/jmoiron/sqlx"
)

type TestMessage struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

type TestQueue struct{}

func (t *TestQueue) GetTopicName() string {
	return "test-topic"
}

func (t *TestQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	log.Println("Message Body: %s", msg.MessageBody)
	return nil
}

func main() {
	err := database.GetDatabaseForEnvironmentRetrying()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := queue.RegisterQueue(&TestQueue{}); err != nil {
		log.Fatal(err.Error())
	}
	queue.PublishMessageToQueueByName("test-topic", TestMessage{
		Field1: "Andrew",
		Field2: 5,
	})
	queue.StartQueue()
	for {
	}
}
