package admin

import "time"

type ID string

type dbAdmin struct {
	CreatedAt      time.Time `db:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at"`
	ID             ID        `db:"_id"`
	EmailAddress   string    `db:"email_address"`
	IsActive       bool      `db:"is_active"`
}

func (d dbAdmin) ToNonDB() Admin {
	return Admin{
		ID:           d.ID,
		EmailAddress: d.EmailAddress,
	}
}

type Admin struct {
	ID           ID
	EmailAddress string
}

type passwordID string

type dbPassword struct {
	CreatedAt      time.Time  `db:"created_at"`
	LastModifiedAt time.Time  `db:"last_modified_at"`
	ID             passwordID `db:"_id"`
	AdminUserID    ID         `db:"admin_user_id"`
	PasswordHash   string     `db:"password_hash"`
	Salt           string     `db:"salt"`
	IsActive       bool       `db:"is_active"`
}

type permissionID string

type dbAccessPermission struct {
	CreatedAt      time.Time    `db:"created_at"`
	LastModifiedAt time.Time    `db:"last_modified_at"`
	ID             permissionID `db:"_id"`
	AdminUserID    ID           `db:"admin_user_id"`
	Permission     Permission   `db:"permission"`
	IsActive       bool         `db:"is_active"`
}

type UserPermissionMapping struct {
	AdminUserID ID         `db:"admin_user_id"`
	Permission  Permission `db:"permission"`
}

type Permission string

const (
	PermissionManagePermissions  Permission = "manage-permissions"
	PermissionViewUserMetrics    Permission = "view-user-metrics"
	PermissionWriteBlog          Permission = "write-blog"
	PermissionPublishBlog        Permission = "publish-blog"
	PermissionEditContentTopics  Permission = "edit-content-topics"
	PermissionEditContentSources Permission = "edit-content-sources"
	PermissionManageBilling      Permission = "manage-billing"
)

var validAdminEmailDomains = map[string]bool{
	"babblegraph.com": true,
}

const AccessTokenCookieName = "BG_ACCESS_TOKEN"

type dbTwoFactorAuthenticationCode struct {
	CreatedAt   time.Time  `db:"created_at"`
	Code        string     `db:"code"`
	AdminUserID ID         `db:"admin_user_id"`
	ExpiresAt   *time.Time `db:"expires_at"`
}

func (d dbTwoFactorAuthenticationCode) ToNonDB() TwoFactorAuthenticationCode {
	return TwoFactorAuthenticationCode{
		Code:        d.Code,
		AdminUserID: d.AdminUserID,
	}
}

type TwoFactorAuthenticationCode struct {
	Code        string
	AdminUserID ID
}

type dbAccessToken struct {
	CreatedAt   time.Time `db:"created_at"`
	Token       string    `db:"token"`
	ExpiresAt   time.Time `db:"expires_at"`
	AdminUserID ID        `db:"admin_user_id"`
}
