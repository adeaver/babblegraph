package email

import (
	"babblegraph/model/users"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const createEmailRecordQuery = "INSERT INTO email_records (_id, user_id, type) VALUES ($1, $2, $3)"

func InsertEmailRecord(tx *sqlx.Tx, id ID, userID users.UserID, emailType EmailType) error {
	if _, err := tx.Exec(createEmailRecordQuery, id, userID, emailType); err != nil {
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

const getEmailUsageForTypeQuery = "SELECT user_id, COUNT(DISTINCT _id) number_emails_sent, MAX(first_opened_at) IS NOT NULL AS has_opened_one_email FROM email_records WHERE type = $1 GROUP BY user_id"

func GetEmailUsageForType(tx *sqlx.Tx, emailType EmailType) ([]EmailUsage, error) {
	var matches []EmailUsage
	if err := tx.Select(&matches, getEmailUsageForTypeQuery, emailType); err != nil {
		return nil, err
	}
	return matches, nil
}

const updateEmailRecordIDWithSESMessageIDQuery = "UPDATE email_records SET ses_message_id = $1 WHERE _id = $2"

func UpdateEmailRecordIDWithSESMessageID(tx *sqlx.Tx, id ID, sesMessageID string) error {
	if _, err := tx.Exec(updateEmailRecordIDWithSESMessageIDQuery, id, sesMessageID); err != nil {
		return err
	}
	return nil
}

func NewEmailRecordID() ID {
	return ID(uuid.New().String())
}
