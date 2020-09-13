package main

import (
	"babblegraph/worker/htmlfetch"
	"babblegraph/worker/htmlparse"
	"babblegraph/worker/languageclassifier"
	"babblegraph/worker/linkhandler"
	"babblegraph/worker/normalizetext"
	"babblegraph/worker/wordsmith"
	"log"

	"github.com/adeaver/babblegraph/lib/database"
	"github.com/adeaver/babblegraph/lib/queue"
)

func main() {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		log.Fatal(err.Error())
	}
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	if err := registerQueues(errs); err != nil {
		log.Fatal(err.Error())
	}
	if err := htmlfetch.PublishToFetchQueue("https://cnnespanol.cnn.com/2020/08/29/la-lucha-de-europa-contra-el-covid-19-pasa-de-los-hospitales-a-las-calles/"); err != nil {
		log.Fatalf(err.Error())
	}
	queue.StartQueue(errs)
	<-errs
}

func registerQueues(errs chan error) error {
	linkHandlerQueue, err := linkhandler.InitializeLinkHandlerQueue(errs)
	if err != nil {
		return err
	}
	queues := []queue.Queue{
		htmlfetch.FetchQueueImpl,
		htmlparse.ParseQueueImpl,
		normalizetext.NormalizeTextQueueImpl,
		languageclassifier.LanguageClassifierQueueImpl,
		*linkHandlerQueue,
	}
	for _, q := range queues {
		if err := queue.RegisterQueue(q); err != nil {
			return err
		}
	}
	return nil
}
