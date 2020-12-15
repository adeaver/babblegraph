package users

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getAllUsersByStatusQuery       = "SELECT * FROM users WHERE status='%s'"
	updateUserStatusByQuery        = "UPDATE users SET status = $1 WHERE email_address = $2 and _id = $3"
	updateUserStatusByEmailAddress = "UPDATE users SET status = $1 WHERE email_address = $2" // prefer update by query
	lookupUserByEmailAddressAndID  = "SELECT * FROM users WHERE _id = $1 AND email_address = $2"
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
