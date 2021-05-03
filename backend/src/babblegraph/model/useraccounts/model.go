package useraccounts

import (
	"babblegraph/model/users"
	"time"
)

type userPasswordID string

type dbUserPassword struct {
	ID           userPasswordID `db:"_id"`
	PasswordHash string         `db:"password_hash"`
	CreatedAt    time.Time      `db:"created_at"`
	UserID       users.UserID   `db:"user_id"`
	Salt         string         `db:"salt"`
}

type userSubscriptionID string

type dbUserSubscription struct {
	ID                userSubscriptionID `db:"_id"`
	CreatedAt         time.Time          `db:"created_at"`
	SubscriptionLevel SubscriptionLevel  `db:"subscription_level"`
	ExpiresAt         time.Time          `db:"expires_at"`
	IsActive          bool               `db:"is_active"`
}

type SubscriptionLevel string

const (
	SubscriptionLevelBetaPremium SubscriptionLevel = "Beta-Premium"
)

func (s SubscriptionLevel) Ptr() *SubscriptionLevel {
	return &s
}
