package email

import (
	"babblegraph/model/users"
	"time"
)

type ID string

type dbEmail struct {
	ID            ID           `db:"_id"`
	SESMessageID  *string      `db:"ses_message_id"`
	UserID        users.UserID `db:"user_id"`
	SentAt        time.Time    `db:"sent_at"`
	FirstOpenedAt *time.Time   `db:"first_opened_at"`
	Type          EmailType    `db:"type"`
}

type EmailType string

const (
	EmailTypeDaily            EmailType = "daily-email"
	EmailTypeUserVerification EmailType = "user-verification"
	EmailTypeUserFeedback     EmailType = "user-feedback"
	EmailTypeUserCreation     EmailType = "user-creation"

	EmailTypePrivacyPolicyUpdateJune2021 = "privacy-policy-update-june-2021"
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

type EmailUsage struct {
	UserID             users.UserID `db:"user_id"`
	NumberOfSentEmails int          `db:"number_emails_sent"`
	HasOpenedOneEmail  bool         `db:"has_opened_one_email"`
}
