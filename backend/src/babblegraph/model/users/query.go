package users

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getAllUsersByStatusQuery = "SELECT * FROM users WHERE status='%s'"

	// In order of preference
	updateUserStatusByQuery        = "UPDATE users SET status = $1, last_modified_at = timezone('utc', now()) WHERE email_address = $2 and _id = $3"
	updateUserStatusByID           = "UPDATE users SET status = $1, last_modified_at = timezone('utc', now()) WHERE _id = $2"
	updateUserStatusByEmailAddress = "UPDATE users SET status = $1, last_modified_at = timezone('utc', now()) WHERE email_address = $2"

	lookupUserByEmailAddressAndID = "SELECT * FROM users WHERE _id = $1 AND email_address = $2"
	lookupUserByEmailAddressQuery = "SELECT * FROM users WHERE email_address = $1"
	lookupUserQuery               = "SELECT * FROM users WHERE _id = $1"

	insertUnverifiedUserQuery = "INSERT INTO users (email_address, status) VALUES ($1, $2) ON CONFLICT DO NOTHING"
)

func GetAllActiveUsers(tx *sqlx.Tx) ([]User, error) {
	var matches []dbUser
	if err := tx.Select(&matches, fmt.Sprintf(getAllUsersByStatusQuery, UserStatusVerified)); err != nil {
		return nil, err
	}
	var out []User
	for _, u := range matches {
		out = append(out, u.ToNonDB())
	}
	return out, nil
}

func UnsubscribeUserForIDAndEmail(tx *sqlx.Tx, userID UserID, emailAddress string) (_didUpdate bool, _err error) {
	res, err := tx.Exec(updateUserStatusByQuery, string(UserStatusUnsubscribed), emailAddress, string(userID))
	if err != nil {
		return false, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	didUpdate := numRows > 0
	return didUpdate, nil
}

func LookupUserForIDAndEmail(tx *sqlx.Tx, userID UserID, emailAddress string) (*User, error) {
	var matches []dbUser
	if err := tx.Select(&matches, lookupUserByEmailAddressAndID, userID, emailAddress); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, nil
	}
	user := matches[0].ToNonDB()
	return &user, nil
}

func AddUserToBlocklistByEmailAddress(tx *sqlx.Tx, emailAddress string, newStatus UserStatus) (_didUpdate bool, _err error) {
	switch newStatus {
	case UserStatusBlocklistBounced,
		UserStatusBlocklistComplaint:
		// no-op
	case UserStatusVerified,
		UserStatusUnverified,
		UserStatusUnsubscribed:
		return false, fmt.Errorf("User status %s not a blocklist", newStatus)
	default:
		return false, fmt.Errorf("Unrecognized status %s", newStatus)
	}
	res, err := tx.Exec(updateUserStatusByEmailAddress, string(newStatus), emailAddress)
	if err != nil {
		return false, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	didUpdate := numRows > 0
	return didUpdate, nil
}

func GetUser(tx *sqlx.Tx, id UserID) (*User, error) {
	var matches []dbUser
	if err := tx.Select(&matches, lookupUserQuery, id); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("Expecting 1 user got %d", len(matches))
	}
	user := matches[0].ToNonDB()
	return &user, nil
}

func LookupUserByEmailAddress(tx *sqlx.Tx, emailAddress string) (*User, error) {
	var matches []dbUser
	if err := tx.Select(&matches, lookupUserByEmailAddressQuery, emailAddress); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, nil
	}
	user := matches[0].ToNonDB()
	return &user, nil
}

func InsertNewUnverifiedUser(tx *sqlx.Tx, emailAddress string) error {
	_, err := tx.Exec(insertUnverifiedUserQuery, emailAddress, UserStatusUnverified)
	return err
}

func SetUserStatusToVerified(tx *sqlx.Tx, id UserID) error {
	_, err := tx.Exec(updateUserStatusByID, UserStatusVerified, id)
	return err
}
