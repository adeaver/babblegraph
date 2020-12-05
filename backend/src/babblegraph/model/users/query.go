package users

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const getAllUsersByStatusQuery = "SELECT * FROM users WHERE status='%s'"

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

const updateUserStatusByQuery = "UPDATE users SET status = $1 WHERE email_address = $2 and _id = $3"

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
