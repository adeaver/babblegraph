package index

import (
	"babblegraph/model/documents"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const insertTermEntry = "INSERT INTO document_term_entries (document_id, term_id, language_code, count) VALUES ($1, $2, $3, $4)"

func InsertTermEntry(tx *sqlx.Tx, documentID documents.DocumentID, term string, documentLanguage wordsmith.LanguageCode, count int64) error {
	_, err := tx.Exec(insertTermEntry, documentID, term, documentLanguage, count)
	return err
}

const getOrderedTermsQuery = "SELECT term_id, SUM(count) total_count, COUNT(*) document_count FROM document_term_entries WHERE language_code=? GROUP BY term_id ORDER BY total DESC"

func GetOrderedTermsForLanguage(tx *sqlx.Tx, languageCode wordsmith.LanguageCode) ([]TermWithStats, error) {
	var matches []dbTermWithStats
	if err := tx.Select(&out, getOrderedTermsQuery, languageCode); err != nil {
		return nil, err
	}
	var out []TermWithStats
	for _, dbMatch := range matches {
		out = append(out, dbMatch.ToNonDB())
	}
	return out, nil
}

func GetTermEntriesForDocument(tx *sqlx.Tx, documentID documents.DocumentID) ([]DocumentTermEntry, error) {
	return nil, nil
}

func GetTermEntriesForTerm(tx *sqlx.Tx, termID wordsmith.LemmaID) ([]DocumentTermEntry, error) {
	return nil, nil
}
