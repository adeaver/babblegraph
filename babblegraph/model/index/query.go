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

const getOrderedTermsQuery = "SELECT term_id, SUM(count) total_count, COUNT(*) document_count FROM document_term_entries WHERE term_id IN (?) GROUP BY term_id ORDER BY total_count DESC"

func GetOrderedTermsForLemmaIDs(tx *sqlx.Tx, lemmaIDs []wordsmith.LemmaID) ([]TermWithStats, error) {
	var matches []dbTermWithStats
	query, args, err := sqlx.In(getOrderedTermsQuery, lemmaIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	var out []TermWithStats
	for _, dbMatch := range matches {
		out = append(out, dbMatch.ToNonDB())
	}
	return out, nil
}

func GetTermEntriesForTerms(tx *sqlx.Tx, searchTerms []wordsmith.LemmaID) ([]DocumentTermEntry, error) {
	var matches []dbDocumentTermEntry
	query, args, err := sqlx.In("SELECT * FROM document_term_entries WHERE term_id IN (?)", searchTerms)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	var out []DocumentTermEntry
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}

func GetTermEntriesForDocuments(tx *sqlx.Tx, documentIDs []documents.DocumentID) ([]DocumentTermEntry, error) {
	var matches []dbDocumentTermEntry
	query, args, err := sqlx.In("SELECT * FROM document_term_entries WHERE document_id IN (?)", documentIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	var out []DocumentTermEntry
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
