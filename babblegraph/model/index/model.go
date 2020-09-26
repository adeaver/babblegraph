package index

import (
	"babblegraph/model/documents"
)

type DocumentTermEntryID string

type TermID string

type dbDocumentTermEntry struct {
	ID         DocumentTermEntryID  `db:"_id"`
	DocumentID documents.DocumentID `db:"document_id"`
	TermID     TermID               `db:"term_id"`
	Count      int64                `db:"count"`
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
	TermID     TermID
	Count      int64
}

type TermWithStats struct {
	TermID        TermID `json:"term_id"`
	TotalCount    int64  `json:"total_count"`
	DocumentCount int64  `json:"document_count"`
}

type dbTermWithStats struct {
	TermID        TermID `db:"term_id"`
	TotalCount    int64  `db:"total_count"`
	DocumentCount int64  `db:"document_count"`
}

func (d dbTermWithTotalCount) ToNonDB() TermWithTotalCount {
	return TermWithTotalCount{
		TermID:     d.TermID,
		TotalCount: d.TotalCount,
	}
}
