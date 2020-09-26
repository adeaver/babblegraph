package index

import (
	"babblegraph/worker/documents"

	"github.com/jmoiron/sqlx"
)

const insertTermEntry = "INSERT INTO document_term_entries (document_id, term, count) VALUES ($1, $2, $3)"

func InsertTermEntry(tx *sqlx.Tx, documentID documents.DocumentID, term string, count int64) error {
	_, err := tx.Exec(insertTermEntry, documentID, term, count)
	return err
}
