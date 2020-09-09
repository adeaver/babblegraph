package main

import (
	"babblegraph/worker/htmlfetch"
	"babblegraph/worker/htmlparse"
	"babblegraph/worker/normalizetext"
	"log"

	"github.com/adeaver/babblegraph/lib/database"
	"github.com/adeaver/babblegraph/lib/queue"
)

func main() {
	err := database.GetDatabaseForEnvironmentRetrying()
	if err != nil {
		log.Fatal(err.Error())
	}
	registerQueues()
	errs := make(chan error, 1)
	if err := queue.PublishMessageToQueueByName(htmlfetch.FetchQueueTopicName, htmlfetch.QueueMessage{
		URL: "https://cnnespanol.cnn.com/2020/08/29/la-lucha-de-europa-contra-el-covid-19-pasa-de-los-hospitales-a-las-calles/",
	}); err != nil {
		log.Fatalf(err.Error())
	}
	queue.StartQueue(errs)
	<-errs
}

func registerQueues() error {
	queues := []queue.Queue{
		htmlfetch.FetchQueueImpl,
		htmlparse.ParseQueueImpl,
		normalizetext.NormalizeTextQueueImpl,
	}
	for _, q := range queues {
		if err := queue.RegisterQueue(q); err != nil {
			return err
		}
	}
	return nil
}
