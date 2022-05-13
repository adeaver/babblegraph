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

	ID                    *PremiumNewsletterSubscriptionID `json:"id,omitempty"`
	PaymentState          PaymentState                     `json:"payment_state"`
	CurrentPeriodEnd      time.Time                        `json:"current_period_end"`
	StripePaymentIntentID *string                          `json:"stripe_payment_intent_id,omitempty"`
	IsAutoRenewEnabled    bool                             `json:"is_auto_renew_enabled"`

	PriceCents       *int64 `json:"price_cents,omitempty"`
	HasValidDiscount bool   `json:"has_valid_discount"`
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

func (p PaymentState) Ptr() *PaymentState {
	return &p
}

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
	// These update types are for when we make an update on our end
	// but are unsure if the remote payment processor has finished processing
	PremiumNewsletterSubscriptionUpdateTypeTransitionToActive PremiumNewsletterSubscriptionUpdateType = "transition-to-active"
	PremiumNewsletterSubscriptionUpdateTypeCanceled           PremiumNewsletterSubscriptionUpdateType = "canceled"

	PremiumNewsletterSubscriptionUpdateTypePaymentMethodAdded PremiumNewsletterSubscriptionUpdateType = "payment-method-added"

	// This update type is for when we receive an update pushed from
	// the remote payment processor
	PremiumNewsletterSubscriptionUpdateTypeRemoteUpdated PremiumNewsletterSubscriptionUpdateType = "remote-updated"
)

func (u PremiumNewsletterSubscriptionUpdateType) Str() string {
	return string(u)
}

func (u PremiumNewsletterSubscriptionUpdateType) Ptr() *PremiumNewsletterSubscriptionUpdateType {
	return &u
}

func GetPremiumNewsletterSubscriptionUpdateTypeFromString(u string) (*PremiumNewsletterSubscriptionUpdateType, error) {
	switch u {
	case PremiumNewsletterSubscriptionUpdateTypeTransitionToActive.Str():
		return PremiumNewsletterSubscriptionUpdateTypeTransitionToActive.Ptr(), nil
	case PremiumNewsletterSubscriptionUpdateTypePaymentMethodAdded.Str():
		return PremiumNewsletterSubscriptionUpdateTypePaymentMethodAdded.Ptr(), nil
	default:
		return nil, fmt.Errorf("unrecognized update type %s", u)
	}
}

type PaymentMethod struct {
	ExternalID     string   `json:"external_id"`
	DisplayMask    string   `json:"display_mask"`
	CardExpiration string   `json:"card_expiration"`
	CardType       CardType `json:"card_type"`
	IsDefault      bool     `json:"is_default"`
}

type CardType string

const (
	CardTypeAmex       CardType = "amex"
	CardTypeVisa       CardType = "visa"
	CardTypeMastercard CardType = "mc"
	CardTypeDiscover   CardType = "discover"
	CardTypeOther      CardType = "other"
)

type newsletterSubscriptionTrialID string

// This is intentionally keyed off of email address
// The idea is that eventually we'll delete unsubscribed users.
type dbNewsletterSubscriptionTrial struct {
	ID           newsletterSubscriptionTrialID `db:"_id"`
	EmailAddress string                        `db:"email_address"`
	CreatedAt    time.Time                     `db:"created_at"`
}

type PromotionCode struct {
	ExternalID string        `json:"external_id"`
	Code       string        `json:"code"`
	Discount   Discount      `json:"discount"`
	Type       PromotionType `json:"promotion_type"`
	IsActive   bool          `json:"is_active"`
}

type PromotionCodeID string

type dbPromotionCode struct {
	CreatedAt           time.Time           `db:"created_at"`
	LastModifiedAt      time.Time           `db:"last_modified_at"`
	ID                  PromotionCodeID     `db:"_id"`
	Type                PromotionType       `db:"type"`
	ExternalIDMappingID externalIDMappingID `db:"external_id_mapping_id"`
}

type dbUserPromotion struct {
	CreatedAt            time.Time            `db:"created_at"`
	LastModifiedAt       time.Time            `db:"last_modified_at"`
	PromotionID          PromotionCodeID      `db:"promotion_id"`
	BillingInformationID BillingInformationID `db:"billing_information_id"`
}

type Discount struct {
	PercentOffBPS  *int64 `json:"percent_off_bps,omitempty"`
	AmountOffCents *int64 `json:"amount_off_cents,omitempty"`
}

type PromotionType string

const (
	PromotionTypeCheckout PromotionType = "checkout"
	PromotionTypeURL      PromotionType = "url"
)

func (p PromotionType) Str() string {
	return string(p)
}

func (p PromotionType) Ptr() *PromotionType {
	return &p
}

func GetPromotionTypeForString(s string) (*PromotionType, error) {
	switch s {
	case PromotionTypeCheckout.Str():
		return PromotionTypeCheckout.Ptr(), nil
	case PromotionTypeURL.Str():
		return PromotionTypeURL.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unrecognized promotion type %s", s)
	}
}
