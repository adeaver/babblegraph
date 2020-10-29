package users

type UserID string

type UserStatus string

const (
	UserStatusVerified     UserStatus = "verified"
	UserStatusUnverified   UserStatus = "unverified"
	UserStatusUnsubscribed UserStatus = "unsubscribed"
)

type dbUser struct {
	ID           UserID     `db:"_id"`
	EmailAddress string     `db:"email_address"`
	Status       UserStatus `db:"status"`
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
