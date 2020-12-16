package email

import (
	"babblegraph/model/users"
	"time"
)

type ID string

type Email struct {
	ID            ID
	SESMessageID  string
	UserID        users.UserID
	SentAt        time.Time
	FirstOpenedAt *time.Time
}

type dbEmail struct {
	ID            ID           `db:"_id"`
	SESMessageID  string       `db:"ses_message_id"`
	UserID        users.UserID `db:"user_id"`
	SentAt        time.Time    `db:"sent_at"`
	FirstOpenedAt *time.Time   `db:"first_opened_at"`
}
