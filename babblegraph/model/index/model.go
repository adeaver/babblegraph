package index

import (
	"babblegraph/model/documents"
	"babblegraph/wordsmith"
)

type DocumentTermEntryID string

type TermID string

type dbDocumentTermEntry struct {
	ID           DocumentTermEntryID    `db:"_id"`
	DocumentID   documents.DocumentID   `db:"document_id"`
	TermID       wordsmith.LemmaID      `db:"term_id"`
	LanguageCode wordsmith.LanguageCode `db:"language_code"`
	Count        int64                  `db:"count"`
}

func (d dbDocumentTermEntry) ToNonDB() DocumentTermEntry {
	return DocumentTermEntry{
		DocumentID: d.DocumentID,
		TermID:     d.TermID,
		Count:      d.Count,
	}
}

type DocumentTermEntry struct {
	DocumentID documents.DocumentID
	TermID     wordsmith.LemmaID
	Count      int64
}

type TermWithStats struct {
	TermID        wordsmith.LemmaID `json:"term_id"`
	TotalCount    int64             `json:"total_count"`
	DocumentCount int64             `json:"document_count"`
}

type dbTermWithStats struct {
	TermID        wordsmith.LemmaID `db:"term_id"`
	TotalCount    int64             `db:"total_count"`
	DocumentCount int64             `db:"document_count"`
}

func (d dbTermWithStats) ToNonDB() TermWithStats {
	return TermWithStats{
		TermID:        d.TermID,
		TotalCount:    d.TotalCount,
		DocumentCount: d.DocumentCount,
	}
}
