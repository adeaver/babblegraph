package users

import "github.com/jmoiron/sqlx"

const (
	getUserStatusAggregationQuery = "SELECT status, COUNT(*) count FROM users GROUP BY status"
)

type UserStatusAggregation struct {
	Status UserStatus `db:"status"`
	Count  int64      `db:"count"`
}

func GetUserStatusAggregation(tx *sqlx.Tx) ([]UserStatusAggregation, error) {
	var out []UserStatusAggregation
	if err := tx.Select(&out, getUserStatusAggregationQuery); err != nil {
		return nil, err
	}
	return out, nil
}
