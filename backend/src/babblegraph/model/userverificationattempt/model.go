package userverificationattempt

import (
	"babblegraph/model/users"
	"time"
)

type ID string

type dbUserVerificationAttempt struct {
	VerificationAttemptID ID           `db:"_id"`
	FulfilledAtTimestamp  *time.Time   `db:"fulfilled_at_timestamp"`
	UserID                users.UserID `db:"user_id"`
}
