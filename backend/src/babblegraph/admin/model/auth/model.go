package auth

import (
	"babblegraph/admin/model/user"
	"time"
)

const AccessTokenCookieName = "BG_ACCESS_TOKEN"

type dbAdmin2FACode struct {
	CreatedAt   time.Time    `db:"created_at"`
	Code        string       `db:"code"`
	AdminUserID user.AdminID `db:"admin_user_id"`
	ExpiresAt   *time.Time   `db:"expires_at"`
}

func (d dbAdmin2FACode) ToNonDB() Admin2FACode {
	return Admin2FACode{
		Code:        d.Code,
		AdminUserID: d.AdminUserID,
	}
}

type Admin2FACode struct {
	Code        string
	AdminUserID user.AdminID
}

type dbAdminAccessToken struct {
	CreatedAt   time.Time    `db:"created_at"`
	Token       string       `db:"token"`
	ExpiresAt   time.Time    `db:"expires_at"`
	AdminUserID user.AdminID `db:"admin_user_id"`
}
