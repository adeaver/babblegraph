package users

import "time"

type UserID string

type UserStatus string

const (
	UserStatusVerified           UserStatus = "verified"
	UserStatusUnverified         UserStatus = "unverified"
	UserStatusUnsubscribed       UserStatus = "unsubscribed"
	UserStatusBlocklistBounced   UserStatus = "blocklist-bounced"
	UserStatusBlocklistComplaint UserStatus = "blocklist-complaint"
)

type dbUser struct {
	CreatedAt      time.Time  `db:"created_at"`
	LastModifiedAt time.Time  `db:"last_modified_at"`
	ID             UserID     `db:"_id"`
	EmailAddress   string     `db:"email_address"`
	Status         UserStatus `db:"status"`
}

func (d dbUser) ToNonDB() User {
	return User{
		ID:           d.ID,
		EmailAddress: d.EmailAddress,
		Status:       d.Status,
	}
}

type User struct {
	ID           UserID     `json:"id"`
	EmailAddress string     `json:"email_address"`
	Status       UserStatus `json:"status"`
}
