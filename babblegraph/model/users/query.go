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
