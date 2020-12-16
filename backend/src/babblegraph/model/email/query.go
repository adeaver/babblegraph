package email

import (
	"babblegraph/model/users"

	"github.com/jmoiron/sqlx"
)

const createEmailRecordQuery = "INSERT INTO email_records (ses_message_id, user_id) VALUES ($1, $2) RETURNING _id"

func CreateEmailRecord(tx *sqlx.Tx, sesMessageID string, userID users.UserID) (*ID, error) {
	rows, err := tx.Queryx(createEmailRecordQuery, sesMessageID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id ID
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
	}
	return &id, nil
}
