package index

import (
	"github.com/adeaver/babblegraph/lib/model/documents"
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
