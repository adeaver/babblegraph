package billing

import (
	"babblegraph/model/users"
	"fmt"
	"time"
)

type BillingInformationID string

type BillingInformation struct {
	UserID           *users.UserID `json:"user_id"`
	StripeCustomerID *string       `json:"stripe_customer_id"`
}

type dbBillingInformation struct {
	CreatedAt           time.Time            `db:"created_at"`
	LastModifiedAt      time.Time            `db:"last_modified_at"`
	ID                  BillingInformationID `db:"_id"`
	UserID              *users.UserID        `db:"user_id"`
	ExternalIDMappingID externalIDMappingID  `db:"external_id_mapping_id"`
}

type externalIDMappingID string

type dbExternalIDMapping struct {
	CreatedAt      time.Time           `db:"created_at"`
	LastModifiedAt time.Time           `db:"last_modified_at"`
	ID             externalIDMappingID `db:"_id"`
	IDType         externalIDType      `db:"id_type"`
	ExternalID     string              `db:"external_id"`
}

type externalIDType string

const (
	externalIDTypeStripe externalIDType = "stripe"
)

type PremiumNewsletterSubscriptionID string

type PremiumNewsletterSubscription struct {
	userID *users.UserID

	PaymentState          PaymentState `json:"payment_state"`
	CurrentPeriodEnd      time.Time    `json:"current_period_end"`
	StripePaymentIntentID *string      `json:"stripe_payment_intent_id,omitempty"`
}

func (p *PremiumNewsletterSubscription) GetUserID() (*users.UserID, error) {
	if p.userID == nil {
		return nil, fmt.Errorf("Premium Newsletter Subscription User ID query called in a context without UserID")
	}
	return p.userID, nil
}

type dbPremiumNewsletterSubscription struct {
	CreatedAt            time.Time                       `db:"created_at"`
	LastModifiedAt       time.Time                       `db:"last_modified_at"`
	ID                   PremiumNewsletterSubscriptionID `db:"_id"`
	BillingInformationID BillingInformationID            `db:"billing_information_id"`
	ExternalIDMappingID  externalIDMappingID             `db:"external_id_mapping_id"`
	IsTerminated         bool                            `db:"is_terminated"`
}

type dbPremiumNewsletterSubscriptionDebounceRecord struct {
	CreatedAt            time.Time            `db:"created_at"`
	LastModifiedAt       time.Time            `db:"last_modified_at"`
	BillingInformationID BillingInformationID `db:"billing_information_id"`
}

type PaymentState int

const (
	// This happens when a user is ineligible for a free
	// trial, and has not paid their subscription
	PaymentStateCreatedUnpaid PaymentState = 0

	// This happens when a user has started a free trial but
	// has not added a payment method - this user technically has an
	// active subscription
	PaymentStateTrialNoPaymentMethod PaymentState = 1

	// This happens when a user has started a free trial and
	// has added a payment method. However, the payment could still fail.
	// This user has an active subscription
	PaymentStateTrialPaymentMethodAdded PaymentState = 2

	// This is a normal subscription or trial with an active payment method
	PaymentStateActive PaymentState = 3

	// This is a subscription that has any error with its payment
	PaymentStateErrored PaymentState = 4

	// This subscription has ended
	PaymentStateTerminated PaymentState = 5
)

type dbPremiumNewsletterSubscriptionSyncRequest struct {
	CreatedAt                       time.Time                               `db:"created_at"`
	LastModifiedAt                  time.Time                               `db:"last_modified_at"`
	PremiumNewsletterSubscriptionID PremiumNewsletterSubscriptionID         `db:"premium_newsletter_subscription_id"`
	UpdateType                      PremiumNewsletterSubscriptionUpdateType `db:"update_type"`
	AttemptNumber                   int64                                   `db:"attempt_number"`
	HoldUntil                       *time.Time                              `db:"hold_until"`
}

type PremiumNewsletterSubscriptionUpdateType string

const (
	PremiumNewsletterSubscriptionUpdateTypeTransitionToActive PremiumNewsletterSubscriptionUpdateType = "transition-to-active"
)

type PaymentMethodID string

type dbPaymentMethod struct {
	ID                   PaymentMethodID      `db:"_id"`
	BillingInformationID BillingInformationID `db:"billing_information_id"`
	IsDefault            bool                 `db:"is_default"`
	ExternalIDMappingID  externalIDMappingID  `db:"external_id_mapping_id"`
}
