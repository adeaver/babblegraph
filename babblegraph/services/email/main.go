package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/email/documentrank"
	"babblegraph/services/email/labels"
	"babblegraph/services/email/wordrank"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func main() {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully connected to main db")
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully connected to wordsmith")
	rankedWordsForSpanish, err := wordrank.GetRankedWords(wordsmith.LanguageCodeSpanish)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully ranked words")
	labelSearchTerms, err := labels.GetLemmaIDsForLabelNames()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("successfully got label search terms")
	var rankedDocumentByLabel map[labels.LabelName][]documentrank.RankedDocument
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		documentCount, err := documents.GetDocumentCountForLanguage(tx, wordsmith.LanguageCodeSpanish)
		if err != nil {
			return err
		}
		log.Println("successfully got document count")
		rankedDocumentByLabel, err = documentrank.GetDocumentsRankedByLabel(tx, documentrank.GetDocumentsRankedByLabelInput{
			RankedWords:      rankedWordsForSpanish,
			LabelSearchTerms: labelSearchTerms,
			DocumentCount:    *documentCount,
		})
		return err
	}); err != nil {
		log.Fatal(err.Error())
	}
	for labelName, documents := range rankedDocumentByLabel {
		for _, doc := range documents {
			log.Println(fmt.Sprintf("Label %s. Doc ID: %s: %f", labelName, doc.DocumentID, doc.Score.ToFloat64()))
		}
	}
}
