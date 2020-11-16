package main

import (
	"babblegraph/services/worker/domains"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
	"fmt"
	"log"
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	linkProcessor, err := linkprocessing.CreateLinkProcessor()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := linkProcessor.AddURLs(domains.GetSeedURLs()); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	for i := 0; i < 4; i++ {
		workerThread := startWorkerThread(linkProcessor, errs)
		go workerThread()
	}
	for {
		select {
		case err := <-errs:
			log.Println(fmt.Sprintf("Saw panic: %s. Starting new worker thread.", err.Error()))
			workerThread := startWorkerThread(linkProcessor, errs)
			go workerThread()
		}
	}
}

func setupDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("Error setting up main-db: %s", err.Error())
	}
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up wordsmith: %s", err.Error())
	}
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up elasticsearch: %s", err.Error())
	}
	return nil
}
