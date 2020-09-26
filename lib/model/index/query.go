package index

import (
	"github.com/adeaver/babblegraph/lib/model/documents"
	"github.com/adeaver/babblegraph/lib/wordsmith"
	"github.com/jmoiron/sqlx"
)

const insertTermEntry = "INSERT INTO document_term_entries (document_id, term_id, language_code, count) VALUES ($1, $2, $3, $4)"

func InsertTermEntry(tx *sqlx.Tx, documentID documents.DocumentID, term string, documentLanguage wordsmith.LanguageCode, count int64) error {
	_, err := tx.Exec(insertTermEntry, documentID, term, documentLanguage, count)
	return err
}
