package user

import (
	"babblegraph/util/email"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	getAdminUserByEmailAddressQuery = "SELECT * FROM admin_user WHERE email_address = $1"
	getAdminUserByIDQuery           = "SELECT * FROM admin_user WHERE _id = $1"
	createUserQuery                 = "INSERT INTO admin_user (email_address) VALUES ($1)"

	createAdminUserPasswordQuery = `INSERT INTO
        admin_user_password (
            admin_user_id,
            password_hash,
            salt,
        ) VALUES ($1, $2, $3)
        ON CONFLICT (admin_user_id)
        SET password_hash = $2, salt = $3
    `
	activeAdminPasswordQuery = `
        UPDATE admin_user_password
        SET is_active = TRUE
        WHERE admin_user_id = $1
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
