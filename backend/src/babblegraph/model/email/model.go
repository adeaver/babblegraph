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

	EmailTypePasswordReset EmailType = "password-reset"

	EmailTypeTrialEndingSoonActionRequired EmailType = "trial-ending-soon-action-required"
	EmailTypeTrialEndingSoon               EmailType = "trial-ending-soon"
	EmailTypePremiumSubscriptionCanceled   EmailType = "premium-subscription-canceled"
	EmailTypePaymentFailureNotification    EmailType = "payment-failure-notification"

	EmailTypeAdminTwoFactorAuthenticationCode EmailType = "admin-two-factor-authentication-code"

	// Deprecated types
	EmailTypeUserFeedbackDEPRECATED                EmailType = "user-feedback"
	EmailTypePrivacyPolicyUpdateJune2021DEPRECATED           = "privacy-policy-update-june-2021"
	EmailTypeUserCreationDEPRECATED                EmailType = "user-creation"
	EmailTypeUserReactivationDEPRECATED            EmailType = "user-reactivation"
	EmailTypeUserExpirationDEPRECATED              EmailType = "user-expiration"
	EmailTypeAccountCreationNotificationDEPRECATED EmailType = "account-creation-notification"
	EmailTypeInitialPremiumAdvertisementDEPRECATED EmailType = "initial-premium-advertisement"
	EmailTypePremiumAnnouncementDEPRECATED         EmailType = "premium-announcement-september-2021"
)

func (e EmailType) Ptr() *EmailType {
	return &e
}

type Recipient struct {
	EmailAddress string
	UserID       users.UserID
}

type bounceRecordID string

type dbBounceRecord struct {
	ID             bounceRecordID `db:"_id"`
	CreatedAt      time.Time      `db:"created_at"`
	LastModifiedAt time.Time      `db:"last_modified_at"`
	UserID         users.UserID   `db:"user_id"`
	LastBounceAt   time.Time      `db:"last_bounce_at"`
	AttemptNumber  int            `db:"attempt_number"`
}
