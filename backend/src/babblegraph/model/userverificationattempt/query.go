package userverificationattempt

import (
	"babblegraph/model/users"
	"babblegraph/util/ptr"

	"github.com/jmoiron/sqlx"
)

const (
	getAllFulfilledAttemptsForUserQuery    = "SELECT * FROM user_verification_attempts WHERE user_id = $1 AND fulfilled_at_timestamp IS NOT NULL"
	insertVerificationAttemptForUserQuery  = "INSERT INTO user_verification_attempts (user_id) VALUES ($1) ON CONFLICT DO NOTHING"
	fulfillVerificationAttemptByIDQuery    = "UPDATE user_verification_attempts SET fulfilled_at_timestamp = (now() at time zone 'utc') WHERE user_id = $1 AND fulfilled_at_timestamp IS NULL"
	getAllPendingVerificationAttemptsQuery = "SELECT * FROM user_verification_attempts WHERE fulfilled_at_timestamp IS NULL"
)

func GetNumberOfFulfilledVerificationAttemptsForUser(tx *sqlx.Tx, userID users.UserID) (*int, error) {
	var matches []dbUserVerificationAttempt
	if err := tx.Select(&matches, getAllFulfilledAttemptsForUserQuery, userID); err != nil {
		return nil, err
	}
	return ptr.Int(len(matches)), nil
}

func GetUserIDsWithPendingVerificationAttempts(tx *sqlx.Tx) ([]users.UserID, error) {
	var matches []dbUserVerificationAttempt
	if err := tx.Select(&matches, getAllPendingVerificationAttemptsQuery); err != nil {
		return nil, err
	}
	var out []users.UserID
	for _, m := range matches {
		out = append(out, m.UserID)
	}
	return out, nil
}

func InsertVerificationAttemptForUser(tx *sqlx.Tx, userID users.UserID) error {
	_, err := tx.Exec(insertVerificationAttemptForUserQuery, userID)
	return err
}

func MarkVerificationAttemptAsFulfilledByUserID(tx *sqlx.Tx, userID users.UserID) error {
	_, err := tx.Exec(fulfillVerificationAttemptByIDQuery, userID)
	return err
}
