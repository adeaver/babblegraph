package documents

import (
	"babblegraph/wordsmith"
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
		callCount++
	}
	if callCount == 0 {
		log.Println("Duplicate document")
		return nil, nil
	}
	return &docID, nil
}

type dbDocumentCount struct {
	Count int64 `db:"count"`
}

func GetDocumentCountForLanguage(tx *sqlx.Tx, languageCode wordsmith.LanguageCode) (*int64, error) {
	var count []dbDocumentCount
	searchLabels := wordsmith.LookupLanguageLabelsForLanguageCode(languageCode)
	query, args, err := sqlx.In("SELECT COUNT(*) count FROM documents WHERE language IN (?)", searchLabels)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&count, sql, args...); err != nil {
		return nil, err
	}
	if len(count) != 1 {
		return nil, fmt.Errorf("Did not get the right value for count. Expected 1, but got %d", len(count))
	}
	return &count[0].Count, nil
}
