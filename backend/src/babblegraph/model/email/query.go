package email

import (
	"babblegraph/model/users"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const createEmailRecordQuery = "INSERT INTO email_records (ses_message_id, user_id, type) VALUES ($1, $2, $3, $4)"

func insertEmailRecord(tx *sqlx.Tx, id ID, sesMessageID string, userID users.UserID, emailType emailType) error {
	if _, err := tx.Exec(createEmailRecordQuery, id, sesMessageID, userID, emailType); err != nil {
		return err
	}
	return nil
}

func newEmailRecordID() ID {
	return ID(uuid.New().String())
}
