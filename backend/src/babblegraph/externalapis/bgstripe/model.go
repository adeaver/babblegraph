package bgstripe

import (
	"babblegraph/model/users"
	"time"
)

type CustomerID string

type StripeCustomer struct {
	BabblegraphUserID    users.UserID
	CustomerID           CustomerID
	DefaultPaymentMethod *PaymentMethodID
}

type customerRelationshipID string

type dbStripeCustomer struct {
	ID                     customerRelationshipID `db:"_id"`
	CreatedAt              time.Time              `db:"created_at"`
	LastModifiedAt         time.Time              `db:"last_modified_at"`
	BabblegraphUserID      users.UserID           `db:"babblegraph_user_id"`
	StripeCustomerID       CustomerID             `db:"stripe_customer_id"`
	DefaultPaymentMethodID *PaymentMethodID       `db:"default_payment_method"`
}

type subscriptionRelationshipID string

type SubscriptionID string

type dbStripeSubscription struct {
	ID                   subscriptionRelationshipID `db:"_id"`
	CreatedAt            time.Time                  `db:"created_at"`
	LastModifiedAt       time.Time                  `db:"last_modified_at"`
	BabblegraphUserID    users.UserID               `db:"babblegraph_user_id"`
	PaymentState         PaymentState               `db:"payment_state"`
	StripeSubscriptionID SubscriptionID             `db:"stripe_subscription_id"`
	StripeProductID      StripeProductID            `db:"stripe_product_id"`
}

type StripeProductID string

const (
	StripeProductIDYearlySubscriptionProd  StripeProductID = "price_1JIMqNJscBSiX47SxOGRUX1p"
	StripeProductIDMonthlySubscriptionProd StripeProductID = "price_1JIMqNJscBSiX47SnYtkOVv6"

	StripeProductIDYearlySubscriptionTest  StripeProductID = "price_1JIMr1JscBSiX47SEEUzRf0e"
	StripeProductIDMonthlySubscriptionTest StripeProductID = "price_1JIMr1JscBSiX47SReF6SdJj"
)

func (s StripeProductID) Str() string {
	return string(s)
}

func (s StripeProductID) Ptr() *StripeProductID {
	return &s
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

type paymentMethodRelationshipID string

type PaymentMethodID string

type dbStripePaymentMethod struct {
	ID                    paymentMethodRelationshipID `db:"_id"`
	CreatedAt             time.Time                   `db:"created_at"`
	LastModifiedAt        time.Time                   `db:"last_modified_at"`
	BabblegraphUserID     users.UserID                `db:"babblegraph_user_id"`
	StripePaymentMethodID PaymentMethodID             `db:"stripe_payment_method_id"`
	CardType              string                      `db:"card_type"`
	LastFourDigits        string                      `db:"last_four_digits"`
	ExpirationMonth       string                      `db:"expiration_month"`
	ExpirationYear        string                      `db:"expiration_year"`
}

type PaymentMethod struct {
	StripePaymentMethodID PaymentMethodID `json:"stripe_payment_method_id"`
	CardType              string          `json:"card_type"`
	LastFourDigits        string          `json:"last_four_digits"`
	ExpirationMonth       string          `json:"expiration_month"`
	ExpirationYear        string          `json:"expiration_year"`
	IsDefault             bool            `json:"is_default"`
}
