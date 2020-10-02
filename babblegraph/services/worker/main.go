package main

import (
	"fmt"
	"log"

	"babblegraph/model/documents"
	"babblegraph/services/worker/queuedefs"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	if err := setupElasticsearchIndices(); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	if err := queuedefs.RegisterQueues(errs); err != nil {
		log.Fatal(err.Error())
	}
	<-errs
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

func setupElasticsearchIndices() error {
	if err := documents.CreateDocumentIndex(); err != nil {
		log.Println(err.Error())
	}
	return nil
}
