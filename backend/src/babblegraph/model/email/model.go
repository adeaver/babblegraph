package email

import (
	"babblegraph/model/users"
	"time"
)

type ID string

type dbEmail struct {
	ID            ID           `db:"_id"`
	SESMessageID  string       `db:"ses_message_id"`
	UserID        users.UserID `db:"user_id"`
	SentAt        time.Time    `db:"sent_at"`
	FirstOpenedAt *time.Time   `db:"first_opened_at"`
	Type          EmailType    `db:"type"`
}

type EmailType string

const (
	EmailTypeDaily EmailType = "daily-email"
)

type Recipient struct {
	EmailAddress string
	UserID       users.UserID
}

// All email templates should use this
type BaseEmailTemplate struct {
	SubscriptionManagementLink string
	HeroImageURL               string
	HomePageURL                string
}
