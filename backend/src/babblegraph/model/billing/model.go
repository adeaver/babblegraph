package billing

import (
	"babblegraph/model/users"
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
