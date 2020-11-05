package userdocuments

import (
	"babblegraph/model/documents"
	"babblegraph/model/users"
	"time"
)

type UserDocumentID string

type dbUserDocument struct {
	ID         UserDocumentID       `db:"_id"`
	UserID     users.UserID         `db:"user_id"`
	DocumentID documents.DocumentID `db:"document_id"`
	SentOn     time.Time            `db:"sent_on"`
}

func (d dbUserDocument) ToNonDB() UserDocument {
	return UserDocument{
		ID:         d.ID,
		UserID:     d.UserID,
		DocumentID: d.DocumentID,
		SentOn:     d.SentOn,
	}
}

type UserDocument struct {
	ID         UserDocumentID       `json:"id"`
	UserID     users.UserID         `json:"user_id"`
	DocumentID documents.DocumentID `json:"document_id"`
	SentOn     time.Time            `json:"sent_on"`
}
