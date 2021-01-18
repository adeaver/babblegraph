package userverificationattempt

import (
	"babblegraph/model/users"
	"time"
)

type id string

type dbUserVerificationAttempt struct {
	VerificationAttemptID id           `db:"_id"`
	FulfilledAtTimestamp  *time.Time   `db:"fulfilled_at_timestamp"`
	UserID                users.UserID `db:"user_id"`
}
