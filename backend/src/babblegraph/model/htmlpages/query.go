package htmlpages

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type InsertHTMLPageInput struct {
	URL      string
	Language *string
	Metadata map[string]string
}

const insertHTMLPageQuery = "INSERT INTO html_pages (url, language, metadata, og_type) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING _id"

func InsertHTMLPage(tx *sqlx.Tx, input InsertHTMLPageInput) (*HTMLPageID, error) {
	var docID HTMLPageID
	var opengraphType *string
	if ogType, ok := input.Metadata["og:type"]; ok {
		opengraphType = &ogType
	}
	rows, err := tx.Query(insertHTMLPageQuery, input.URL, input.Language, dbMetadata(input.Metadata), opengraphType)
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
