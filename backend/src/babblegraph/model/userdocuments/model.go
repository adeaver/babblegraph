package userdocuments

import (
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"time"
)

type UserDocumentID string

type dbUserDocument struct {
	ID          UserDocumentID       `db:"_id"`
	UserID      users.UserID         `db:"user_id"`
	DocumentID  documents.DocumentID `db:"document_id"`
	SentOn      time.Time            `db:"sent_on"`
	EmailID     *email.ID            `db:"email_id"`
	DocumentURL *string              `db:"document_url"`
}

func (d dbUserDocument) ToNonDB() UserDocument {
	return UserDocument{
		ID:          d.ID,
		UserID:      d.UserID,
		DocumentID:  d.DocumentID,
		SentOn:      d.SentOn,
		EmailID:     d.EmailID,
		DocumentURL: d.DocumentURL,
	}
}

type UserDocument struct {
	ID          UserDocumentID       `json:"id"`
	UserID      users.UserID         `json:"user_id"`
	DocumentID  documents.DocumentID `json:"document_id"`
	SentOn      time.Time            `json:"sent_on"`
	EmailID     *email.ID            `json:"email_id"`
	DocumentURL *string              `json:"document_url"`
}
