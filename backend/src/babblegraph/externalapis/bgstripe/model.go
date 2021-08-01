package bgstripe

import (
	"babblegraph/model/users"
	"time"
)

type CustomerID string

type StripeCustomer struct {
	BabblegraphUserID    users.UserID
	CustomerID           CustomerID
	ActiveSubscriptionID *SubscriptionID
	OtherSubscriptionID  []SubscriptionID
}

type customerRelationshipID string

type dbStripeCustomer struct {
	ID                customerRelationshipID `db:"_id"`
	CreatedAt         time.Time              `db:"created_at"`
	LastModifiedAt    time.Time              `db:"last_modified_at"`
	BabblegraphUserID users.UserID           `db:"babblegraph_user_id"`
	StripeCustomerID  CustomerID             `db:"stripe_customer_id"`
}

type subscriptionRelationshipID string

type SubscriptionID string

type dbStripeSubscription struct {
	ID                   subscriptionRelationshipID `db:"_id"`
	CreatedAt            time.Time                  `db:"created_at"`
	LastModifiedAt       time.Time                  `db:"last_modified_at"`
	BabblegraphUserID    users.UserID               `db:"babblegraph_user_id"`
	PaymentState         PaymentState               `db:"payment_state"`
	StripeSubscriptionID SubscriptionID             `db:"subscription_id"`
	StripeProductID      StripeProductID            `db:"stripe_product_id"`
	StripeClientSecret   string                     `db:"stripe_client_secret"`
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

type PaymentState int

const (
	// This happens when a user is ineligible for a free
	// trial, and has not paid their subscription
	PaymentStateCreatedUnpaid PaymentState = 0

	// This happens when a user has started a free trial but
	// has not added a payment method - this user technically has an
	// active subscription
	PaymentStateActiveNoPaymentMethod PaymentState = 1

	// This is a normal subscription or trial with an active payment method
	PaymentStateActive PaymentState = 2

	// This is a subscription that has any error with its payment
	PaymentStateErrored PaymentState = 3

	// This subscription has ended
	PaymentStateTerminated PaymentState = 4
)
