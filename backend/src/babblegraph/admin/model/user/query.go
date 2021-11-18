package user

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getAdminUserByEmailAddressQuery = "SELECT * FROM admin_user WHERE email_address = $1"
	getAdminUserByIDQuery           = "SELECT * FROM admin_user WHERE _id = $1"
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
