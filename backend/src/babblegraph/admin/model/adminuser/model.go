package adminuser

import "time"

type ID string

type dbAdminUser struct {
	CreatedAt      time.Time `db:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at"`
	ID             ID        `db:"_id"`
	EmailAddress   string    `db:"email_address"`
	IsActive       bool      `db:"is_active"`
}

type adminPasswordID string

type dbAdminUserPassword struct {
	CreatedAt    time.Time       `db:"created_at"`
	ID           adminPasswordID `db:"_id"`
	AdminUserID  ID              `db:"admin_user_id"`
	PasswordHash string          `db:"password_hash"`
	Salt         string          `db:"salt"`
}

type adminPermissionID string

type dbAdminAccessPermission struct {
	CreatedAt      time.Time         `db:"created_at"`
	LastModifiedAt time.Time         `db:"last_modified_at"`
	ID             adminPermissionID `db:"_id"`
	AdminUserID    ID                `db:"admin_user_id"`
	Permission     Permission        `db:"permission"`
	IsActive       bool              `db:"is_active"`
}

type Permission string
