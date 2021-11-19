package user

import (
	"babblegraph/util/email"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
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

	lookupAdminUserPermissionQuery = `
        SELECT * FROM admin_user_permission
        WHERE admin_user_id = $1 AND is_active = TRUE
    `
	createAdminUserPermissionQuery = `
        INSERT INTO admin_user_permission (
            admin_user_id, permission, is_active
        ) VALUES ($1, $2, TRUE)
    `
	deactivateAdminUserPermissionQuery = `
        UPDATE admin_user_permission SET
            is_active = FALSE
        WHERE
            admin_user_id = $1 AND permission = TRUE
    `
)

func GetAdminUser(tx *sqlx.Tx, id AdminID) (*AdminUser, error) {
	var adminUsers []dbAdminUser
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

func LookupAdminUserByEmailAddress(tx *sqlx.Tx, emailAddress string) (*AdminUser, error) {
	var adminUsers []dbAdminUser
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

func CreateAdminUserPassword(tx *sqlx.Tx, adminUserID AdminID, password string) error {
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

func ActivateAdminUserPassword(tx *sqlx.Tx, adminUserID AdminID) error {
	if _, err := tx.Exec(activeAdminPasswordQuery, adminUserID); err != nil {
		return err
	}
	return nil
}

func ValidateAdminUserPassword(tx *sqlx.Tx, adminUserID AdminID, password string) error {
	var passwords []dbAdminUserPassword
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

func ValidateAdminUserPermission(tx *sqlx.Tx, adminUserID AdminID, permission Permission) error {
	var permissions []dbAdminAccessPermission
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
