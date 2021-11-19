package user

import "time"

type AdminID string

type dbAdminUser struct {
	CreatedAt      time.Time `db:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at"`
	AdminID        AdminID   `db:"_id"`
	EmailAddress   string    `db:"email_address"`
	IsActive       bool      `db:"is_active"`
}

func (d dbAdminUser) ToNonDB() AdminUser {
	return AdminUser{
		AdminID:      d.AdminID,
		EmailAddress: d.EmailAddress,
	}
}

type AdminUser struct {
	AdminID      AdminID
	EmailAddress string
}

type adminPasswordID string

type dbAdminUserPassword struct {
	CreatedAt      time.Time       `db:"created_at"`
	LastModifiedAt time.Time       `db:"last_modified_at"`
	ID             adminPasswordID `db:"_id"`
	AdminUserID    AdminID         `db:"admin_user_id"`
	PasswordHash   string          `db:"password_hash"`
	Salt           string          `db:"salt"`
	IsActive       bool            `db:"is_active"`
}

type adminPermissionID string

type dbAdminAccessPermission struct {
	CreatedAt      time.Time         `db:"created_at"`
	LastModifiedAt time.Time         `db:"last_modified_at"`
	ID             adminPermissionID `db:"_id"`
	AdminUserID    AdminID           `db:"admin_user_id"`
	Permission     Permission        `db:"permission"`
	IsActive       bool              `db:"is_active"`
}

type Permission string

var validAdminEmailDomains = map[string]bool{
	"babblegraph.com": true,
}
