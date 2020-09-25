package index

import (
	"babblegraph/worker/documents"
)

type DocumentTermEntryID string

type dbDocumentTermEntry struct {
	ID         DocumentTermEntryID  `db:"_id"`
	DocumentID documents.DocumentID `db:"document_id"`
	Term       string               `db:"term"`
	Count      int64                `db:"count"`
}

func (d dbDocumentTermEntry) ToNonDB() DocumentTermEntry {
	return DocumentTermEntry{
		DocumentID: d.DocumentID,
		Term:       d.Term,
		Count:      d.Count,
	}
}

type DocumentTermEntry struct {
	DocumentID documents.DocumentID
	Term       string
	Count      int64
}
