package unsubscribereason

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	insertUnsubscribeReasonForUserQuery = "INSERT INTO unsubscribe_reasons (user_id, language_code, reason) VALUES ($1, $2, $3)"
)

func InsertUnsubscribeReason(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, reason string) error {
	if _, err := tx.Exec(insertUnsubscribeReasonForUserQuery, userID, languageCode, reason); err != nil {
		return err
	}
	return nil
}
