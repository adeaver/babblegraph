package email

import (
	"babblegraph/model/users"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getBounceRecordForUserQuery    = "SELECT * FROM email_bounce_records WHERE user_id = $1"
	upsertBounceRecordForUserQuery = `INSERT INTO
        email_bounce_records (
            user_id, attempt_number, last_bounce_at
        ) VALUES (
            $1, $2, timezone('utc', now())
        ) ON CONFLICT (user_id) DO UPDATE SET
            attempt_number = $2,
            last_bounce_at = timezone('utc', now()),
            last_modified_at = timezone('utc', now())
        `

	maximumNumberOfBounces = 4
)

var daysSinceLastBouncePerAttempt = []time.Duration{
	24 * time.Hour,      // First bounce, wait a full day before sending another email
	3 * 24 * time.Hour,  // Second bounce, wait three days before sending another email
	7 * 24 * time.Hour,  // Third bounce, wait a week before sending another email
	14 * 24 * time.Hour, // Fourth bounce, wait two weeks before sending another email
}

func lookupBounceRecordForUser(tx *sqlx.Tx, userID users.UserID) (*dbBounceRecord, error) {
	var matches []dbBounceRecord
	err := tx.Select(&matches, getBounceRecordForUserQuery, userID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one bounce record for user %s, but got %d", userID, len(matches))
	default:
		m := matches[0]
		return &m, nil
	}
}

func IsUserOnBounceQuarantine(tx *sqlx.Tx, userID users.UserID) (bool, error) {
	bounceRecord, err := lookupBounceRecordForUser(tx, userID)
	switch {
	case err != nil:
		return false, err
	case bounceRecord == nil:
		return false, nil
	default:
		quarantinePeriod := daysSinceLastBouncePerAttempt[bounceRecord.AttemptNumber]
		return time.Now().Before(bounceRecord.LastBounceAt.Add(quarantinePeriod)), nil
	}
}

func HandleBouncedEmail(tx *sqlx.Tx, userID users.UserID) (_shouldBlocklist bool, _err error) {
	var shouldBlocklist bool
	var nextAttemptNumber int
	bounceRecord, err := lookupBounceRecordForUser(tx, userID)
	switch {
	case err != nil:
		return false, err
	case bounceRecord == nil:
		nextAttemptNumber = 1
	default:
		quarantinePeriod := daysSinceLastBouncePerAttempt[bounceRecord.AttemptNumber]
		if time.Now().After(bounceRecord.LastBounceAt.Add(quarantinePeriod).Add(2 * 24 * time.Hour)) {
			// It has been two days since the last quarantine period ended, we want to reset the counter
			nextAttemptNumber = 1
		} else {
			nextAttemptNumber = bounceRecord.AttemptNumber + 1
		}
		shouldBlocklist = nextAttemptNumber > maximumNumberOfBounces
	}
	if _, err := tx.Exec(upsertBounceRecordForUserQuery, userID, nextAttemptNumber); err != nil {
		return false, err
	}
	return shouldBlocklist, nil
}
