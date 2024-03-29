package admin

import (
	"babblegraph/util/email"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	getAllAdminUsersQuery           = "SELECT * FROM admin_user WHERE is_active = TRUE"
	getAdminUserByEmailAddressQuery = "SELECT * FROM admin_user WHERE email_address = $1 AND is_active = TRUE"
	getAdminUserByIDQuery           = "SELECT * FROM admin_user WHERE _id = $1 AND is_active = TRUE"
	createUserQuery                 = "INSERT INTO admin_user (email_address) VALUES ($1)"

	createAdminUserPasswordQuery = `INSERT INTO
        admin_user_password (
            admin_user_id,
            password_hash,
            salt
        ) VALUES ($1, $2, $3)
        ON CONFLICT (admin_user_id) DO UPDATE
        SET password_hash = $2, salt = $3
    `
	activeAdminPasswordQuery = `
        UPDATE admin_user_password
        SET is_active = TRUE
        WHERE admin_user_id = $1
    `
	validateAdminPasswordQuery = `
        SELECT * FROM admin_user_password
        WHERE admin_user_id = $1 AND is_active = TRUE
    `

	getAllActiveUserPermissionsQuery = `
        SELECT * FROM admin_access_permission
        WHERE is_active = TRUE
    `
	lookupAdminUserPermissionQuery = `
        SELECT * FROM admin_access_permission
        WHERE admin_user_id = $1 AND permission = $2 AND is_active = TRUE
    `
	updateAdminUserPermissionQuery = `
        INSERT INTO admin_access_permission (
            admin_user_id, permission, is_active
        ) VALUES (
            $1, $2, $3
        ) ON CONFLICT (admin_user_id, permission)
        DO UPDATE SET is_active=$3
    `
)

func GetAllAdminUsers(tx *sqlx.Tx) ([]Admin, error) {
	var adminUsers []dbAdmin
	if err := tx.Select(&adminUsers, getAllAdminUsersQuery); err != nil {
		return nil, err
	}
	var out []Admin
	for _, u := range adminUsers {
		out = append(out, u.ToNonDB())
	}
	return out, nil
}

func GetAdminUser(tx *sqlx.Tx, id ID) (*Admin, error) {
	var adminUsers []dbAdmin
	err := tx.Select(&adminUsers, getAdminUserByIDQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(adminUsers) == 0:
		return nil, fmt.Errorf("No admin user found")
	case len(adminUsers) > 1:
		return nil, fmt.Errorf("Expected one user, got %d", len(adminUsers))
	default:
		out := adminUsers[0].ToNonDB()
		return &out, nil
	}
}

func LookupAdminUserByEmailAddress(tx *sqlx.Tx, emailAddress string) (*Admin, error) {
	var adminUsers []dbAdmin
	err := tx.Select(&adminUsers, getAdminUserByEmailAddressQuery, emailAddress)
	switch {
	case err != nil:
		return nil, err
	case len(adminUsers) == 0:
		return nil, nil
	case len(adminUsers) > 1:
		return nil, fmt.Errorf("Expected at most one user, got %d", len(adminUsers))
	default:
		out := adminUsers[0].ToNonDB()
		return &out, nil
	}
}

func CreateAdminUser(tx *sqlx.Tx, emailAddress string) error {
	formattedEmailAddress := email.FormatEmailAddress(emailAddress)
	emailDomain := strings.Split(formattedEmailAddress, "@")[1]
	if _, ok := validAdminEmailDomains[emailDomain]; !ok {
		return fmt.Errorf("invalid email domain")
	}
	if _, err := tx.Exec(createUserQuery, formattedEmailAddress); err != nil {
		return err
	}
	return nil
}

func CreateAdminUserPassword(tx *sqlx.Tx, adminUserID ID, password string) error {
	if !validatePasswordMeetsRequirements(password) {
		return fmt.Errorf("Invalid password")
	}
	salt, err := generatePasswordSalt()
	if err != nil {
		return err
	}
	passwordHash, err := generatePasswordHash(password, *salt)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(createAdminUserPasswordQuery, adminUserID, *passwordHash, *salt); err != nil {
		return err
	}
	return nil
}

func ActivateAdminUserPassword(tx *sqlx.Tx, adminUserID ID) error {
	if _, err := tx.Exec(activeAdminPasswordQuery, adminUserID); err != nil {
		return err
	}
	return nil
}

func ValidateAdminUserPassword(tx *sqlx.Tx, adminUserID ID, password string) error {
	var passwords []dbPassword
	err := tx.Select(&passwords, validateAdminPasswordQuery, adminUserID)
	switch {
	case err != nil:
		return err
	case len(passwords) == 0:
		return fmt.Errorf("No active passwords")
	case len(passwords) > 1:
		return fmt.Errorf("Expected only one password")
	}
	return comparePasswords(passwords[0].PasswordHash, password, passwords[0].Salt)
}

func ValidateAdminUserPermission(tx *sqlx.Tx, adminUserID ID, permission Permission) error {
	var permissions []dbAccessPermission
	err := tx.Select(&permissions, lookupAdminUserPermissionQuery, adminUserID, permission)
	switch {
	case err != nil:
		return err
	case len(permissions) == 0:
		return fmt.Errorf("no permission")
	case len(permissions) > 1:
		return fmt.Errorf("expected at most one permission")
	}
	return nil
}

func GetAllActiveUserPermissions(tx *sqlx.Tx) ([]UserPermissionMapping, error) {
	var permissions []dbAccessPermission
	if err := tx.Select(&permissions, getAllActiveUserPermissionsQuery); err != nil {
		return nil, err
	}
	var out []UserPermissionMapping
	for _, p := range permissions {
		out = append(out, UserPermissionMapping{
			AdminUserID: p.AdminUserID,
			Permission:  p.Permission,
		})
	}
	return out, nil
}

func UpsertUserPermission(tx *sqlx.Tx, adminID ID, permission Permission, isActive bool) error {
	if _, err := tx.Exec(updateAdminUserPermissionQuery, adminID, permission, isActive); err != nil {
		return err
	}
	return nil
}
