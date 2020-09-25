package indexer

import (
	"babblegraph/worker/documents"
	"babblegraph/worker/index"
	"babblegraph/worker/storage"
	"strings"

	"github.com/jmoiron/sqlx"
)

func IndexTermsForFile(tx *sqlx.Tx, documentID documents.DocumentID, filename storage.FileIdentifier) error {
	documentBytes, err := storage.ReadFile(filename)
	if err != nil {
		return err
	}
	return insertTermsForDocument(tx, documentID, string(documentBytes))
}

func insertTermsForDocument(tx *sqlx.Tx, documentID documents.DocumentID, documentBody string) error {
	termCounts := make(map[string]int64)
	tokens := strings.Split(documentBody, " ")
	for _, token := range tokens {
		count, ok := termCounts[token]
		if !ok {
			count = 0
		}
		termCounts[token] = count + 1
	}
	for term, count := range termCounts {
		if err := index.InsertTermEntry(tx, documentID, term, count); err != nil {
			return err
		}
	}
	return nil
}
