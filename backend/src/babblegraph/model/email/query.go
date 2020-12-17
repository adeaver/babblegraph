package email

import (
	"babblegraph/model/users"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const createEmailRecordQuery = "INSERT INTO email_records (_id, ses_message_id, user_id, type) VALUES ($1, $2, $3, $4)"

func InsertEmailRecord(tx *sqlx.Tx, id ID, sesMessageID string, userID users.UserID, emailType EmailType) error {
	if _, err := tx.Exec(createEmailRecordQuery, id, sesMessageID, userID, emailType); err != nil {
		return err
	}
	return nil
}

const setEmailFirstOpenedQuery = "UPDATE email_records SET first_opened_at = timezone('utc', now()) WHERE _id = $1 AND first_opened_at IS NULL"

func SetEmailFirstOpened(tx *sqlx.Tx, id ID) error {
	if _, err := tx.Exec(setEmailFirstOpenedQuery, id); err != nil {
		return err
	}
	return nil
}

func NewEmailRecordID() ID {
	return ID(uuid.New().String())
}
