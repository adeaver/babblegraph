package adminuser

const (
	createAdminUserPasswordQuery = "INSERT INTO admin_user_passwords (admin_user_id, password_hash, salt) VALUES ($1, $2, $3) ON CONFLICT admin_user_id SET password_hash = $2, salt = $3"
	getAdminUserPasswordQuery    = "SELECT * FROM admin_user_passwords WHERE admin_user_id = $1"

	createAdminUserQuery     = "INSERT INTO admin_users (email_address) VALUES ($1) ON CONFLICT (email_address) SET is_active = TRUE"
	deactivateAdminUserQuery = "UPDATE admin_users SET is_active = FALSE WHERE _id = $1"

	doesAdminUserHaveValidPermissionQuery = "SELECT * FROM admin_access_permission WHERE admin_user_id = $1 AND permission = $2 AND is_active = TRUE"

	addAdminUserPermissionQuery    = "INSERT INTO admin_access_permissions (admin_user_id, permission) VALUES ($1, $2) ON CONFLICT (admin_user_id, permission) SET is_active = TRUE"
	revokeAdminUserPermissionQuery = "UPDATE admin_access_permissions SET is_active = FALSE WHERE admin_user_id = $1 AND permission = $2"
)
