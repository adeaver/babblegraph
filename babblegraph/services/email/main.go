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
		log.Fatal(fmt.Errorf("Error setting up elasticsearch: %s", err.Error()))
	}
	log.Println("successfully connected to elasticsearch")
	labelSearchTerms, err := labels.GetLemmaIDsForLabelNames()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully got label search terms")
	docs, err := documents.FindDocumentsContainingTerms([]wordsmith.LemmaID{wordsmith.LemmaID("11b024c4-f772-464d-90a1-9893df2d2094")})
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(docs) > 0 {
		log.Println("Got top doc %+v, for label %s", docs[0], "none")
	}
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
