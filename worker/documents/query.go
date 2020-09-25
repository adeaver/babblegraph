package documents

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type InsertDocumentInput struct {
	URL      string
	Language *string
	Metadata map[string]string
}

const insertDocumentQuery = "INSERT INTO documents (url, language, metadata) VALUES ($1, $2, $3) RETURNING _id"

func InsertDocument(tx *sqlx.Tx, input InsertDocumentInput) (*DocumentID, error) {
	metadataPairBytes, err := json.Marshal(makeMetadataPairs(input.Metadata))
	if err != nil {
		return nil, err
	}
	metadataPairJson := types.JsonText(string(metadataPairBytes))
	var docID DocumentID
	if err := tx.Select(&docID, insertDocumentQuery, input.URL, input.Language, metadataPairJson); err != nil {
		return nil, err
	}
	return &docID, nil
}

func makeMetadataPairs(metadata map[string]string) []MetadataPair {
	var out []MetadataPair
	for key, value := range metadata {
		out = append(out, MetadataPair{
			Key:   key,
			Value: value,
		})
	}
	return out
}
