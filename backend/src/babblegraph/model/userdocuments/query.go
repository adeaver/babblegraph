package userdocuments

import (
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const insertQuery = `INSERT INTO user_documents (_id, user_id, document_id, email_id, document_url) VALUES ($1, $2, $3, $4, $5)`

func InsertDocumentForUserAndReturnID(tx *sqlx.Tx, userID users.UserID, emailRecordID email.ID, doc documents.Document) (*UserDocumentID, error) {
	userDocumentID := UserDocumentID(uuid.New().String())
	docURL := ptr.String(doc.URL)
	if doc.Metadata.URL != nil {
		docURL = doc.Metadata.URL
	}
	urlWithProtocol, err := urlparser.EnsureProtocol(*docURL)
	if err != nil {
		return nil, fmt.Errorf("Got error ensuring protocol for URL %s: %s", *docURL, err.Error())
	}
	if _, err := tx.Exec(insertQuery, userDocumentID, userID, doc.ID, emailRecordID, urlWithProtocol); err != nil {
		return nil, err
	}
	return &userDocumentID, nil
}

const selectQuery = `SELECT * FROM user_documents WHERE user_id = '%s'`

func GetDocumentIDsSentToUser(tx *sqlx.Tx, userID users.UserID) ([]documents.DocumentID, error) {
	var matches []dbUserDocument
	if err := tx.Select(&matches, fmt.Sprintf(selectQuery, string(userID))); err != nil {
		return nil, err
	}
	var out []documents.DocumentID
	for _, match := range matches {
		out = append(out, match.DocumentID)
	}
	return out, nil
}

const selectByIDQuery = `SELECT * FROM user_documents WHERE _id = $1`

func GetUserDocumentID(tx *sqlx.Tx, id UserDocumentID) (*UserDocumentID, error) {
	var matches []dbUserDocument
	if err := tx.Select(&matches, selectByIDQuery, id); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("Expected exactly one record, but got %d", len(matches))
	}
	return &matches[0].ToNonDB(), nil
}
