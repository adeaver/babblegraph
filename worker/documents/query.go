package documents

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type InsertDocumentInput struct {
	URL      string
	Language *string
	Metadata map[string]string
}

const insertDocumentQuery = "INSERT INTO documents (url, language, metadata) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING _id"

func InsertDocument(tx *sqlx.Tx, input InsertDocumentInput) (*DocumentID, error) {
	var docID DocumentID
	rows, err := tx.Query(insertDocumentQuery, input.URL, input.Language, dbMetadata(input.Metadata))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	callCount := 0
	for rows.Next() {
		if callCount != 0 {
			return nil, fmt.Errorf("Insert returned multiple rows")
		}
		if err := rows.Scan(&docID); err != nil {
			return nil, err
		}
	}
	if callCount == 0 {
		log.Println("Duplicate document")
		return nil, nil
	}
	return &docID, nil
}
