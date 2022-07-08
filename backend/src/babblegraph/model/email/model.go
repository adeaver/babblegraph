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

	EmailTypeGoodbye EmailType = "goodbye"

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
