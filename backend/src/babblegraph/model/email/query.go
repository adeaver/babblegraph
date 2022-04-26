package email

import (
	"babblegraph/model/users"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	createEmailRecordQuery                   = "INSERT INTO email_records (_id, user_id, type) VALUES ($1, $2, $3)"
	setEmailFirstOpenedQuery                 = "UPDATE email_records SET first_opened_at = timezone('utc', now()) WHERE _id = $1 AND first_opened_at IS NULL"
	updateEmailRecordIDWithSESMessageIDQuery = "UPDATE email_records SET ses_message_id = $1 WHERE _id = $2"
	updateEmailRecordSentAtQuery             = "UPDATE email_records SET sent_at = timezone('utc', now()) WHERE _id = $1"
)

func InsertEmailRecord(tx *sqlx.Tx, id ID, userID users.UserID, emailType EmailType) error {
	if _, err := tx.Exec(createEmailRecordQuery, id, userID, emailType); err != nil {
		return err
	}
	return nil
}

func SetEmailFirstOpened(tx *sqlx.Tx, id ID) error {
	if _, err := tx.Exec(setEmailFirstOpenedQuery, id); err != nil {
		return err
	}
	return nil
}

func UpdateEmailRecordIDWithSESMessageID(tx *sqlx.Tx, id ID, sesMessageID string) error {
	if _, err := tx.Exec(updateEmailRecordIDWithSESMessageIDQuery, sesMessageID, id); err != nil {
		return err
	}
	return nil
}

func NewEmailRecordID() ID {
	return ID(uuid.New().String())
}

func SetEmailRecordSentAtTime(tx *sqlx.Tx, id ID) error {
	if _, err := tx.Exec(updateEmailRecordSentAtQuery, id); err != nil {
		return err
	}
	return nil
}
