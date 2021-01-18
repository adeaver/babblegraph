package userverificationattempt

import (
	"babblegraph/model/users"
	"babblegraph/util/ptr"

	"github.com/jmoiron/sqlx"
)

const (
	getAllFulfilledAttemptsForUserQuery   = "SELECT * FROM user_verification_attempts WHERE user_id = $1 AND fulfilled_at_timestamp IS NOT NULL"
	insertVerificationAttemptForUserQuery = "INSERT INTO user_verification_attempts (user_id) VALUES ($1) ON CONFLICT DO NOTHING"
	fulfillVerificationAttemptByIDQuery   = "UPDATE user_verification_attempts SET fulfilled_at_timestamp = (now() at time zone 'utc') WHERE _id = $1"
)

func GetNumberOfFulfilledVerificationAttemptsForUser(tx *sqlx.Tx, userID users.UserID) (*int, error) {
	var matches []dbUserVerificationAttempt
	if err := tx.Select(&matches, getAllFulfilledAttemptsForUserQuery, userID); err != nil {
		return nil, err
	}
	return ptr.Int(len(matches)), nil
}

func InsertVerificationAttemptForUser(tx *sqlx.Tx, userID users.UserID) error {
	_, err := tx.Exec(insertVerificationAttemptForUserQuery, userID)
	return err
}

func MarkVerificationAttemptAsFulfilled(tx *sqlx.Tx, id ID) error {
	_, err := tx.Exec(fulfillVerificationAttemptByIDQuery, id)
	return err
}
