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
	StripeSubscriptionID SubscriptionID             `db:"subscription_id"`
	IsActive             bool                       `db:"is_active"`
}
