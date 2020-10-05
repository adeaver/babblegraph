package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/email/labels"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
	"fmt"
	"log"
)

func main() {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully connected to main db")
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully connected to wordsmith db")
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up elasticsearch: %s", err.Error())
	}
	log.Println("successfully connected to elasticsearch")
	labelSearchTerms, err := labels.GetLemmaIDsForLabelNames()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully got label search terms")
	for label, terms := range labelSearchTerms {
		docs, err := documents.FindDocumentsContainingTerms(terms)
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(docs) > 0 {
			log.Println("Got top doc %+v, for label %s", docs[0], label)
		}
	}
}
