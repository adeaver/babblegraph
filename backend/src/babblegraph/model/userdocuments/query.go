package userdocuments

import (
	"babblegraph/model/documents"
	"babblegraph/model/users"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const insertQuery = `INSERT INTO user_documents (user_id, document_id) VALUES ($1, $2)`

func InsertDocumentIDsForUser(tx *sqlx.Tx, userID users.UserID, docIDs []documents.DocumentID) error {
	for _, docID := range docIDs {
		if _, err := tx.Exec(insertQuery, userID, docID); err != nil {
			return err
		}
	}
	return nil
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
