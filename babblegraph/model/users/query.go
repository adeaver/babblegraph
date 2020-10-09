package users

import "github.com/jmoiron/sqlx"

const getAllUsersByStatusQuery = "SELECT * FROM users WHERE status=?"

func GetAllActiveUsers(tx *sqlx.Tx) ([]User, error) {
	var matches []dbUser
	if err := tx.Select(&matches, getAllUsersByStatusQuery, UserStatusVerified); err != nil {
		return nil, err
	}
	var out []User
	for _, u := range matches {
		out = append(out, u.ToNonDB())
	}
	return out, nil
}
