package useraccountsnotifications

import (
	"babblegraph/model/users"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getOutstandingNotificationsToFulfillQuery = "SELECT * FROM user_account_notification_requests WHERE hold_until < CURRENT_TIMESTAMP AND fulfilled_at IS NULL"

	insertNotificationRequestForUserQuery               = "INSERT INTO user_account_notification_requests (user_id, type, hold_until) VALUES ($1, $2, $3)"
	getMostRecentNotificationRequestForUserAndTypeQuery = "SELECT * FROM user_account_notification_requests WHERE user_id = $1 AND type = $2 ORDER BY created_at DESC LIMIT 1"

	fulfillNotificationRequestQuery         = "UPDATE user_account_notification_requests SET fulfilled_at = timezone('utc', now()) WHERE _id = $1"
	insertNotificationRequestDebounceRecord = "INSERT INTO user_account_notification_request_debounce_fulfillment_records (notification_request_id) VALUES ($1)"
)

func GetNotificationsToFulfill(tx *sqlx.Tx) ([]NotificationRequest, error) {
	var matches []dbNotificationRequest
	if err := tx.Select(&matches, getOutstandingNotificationsToFulfillQuery); err != nil {
		return nil, err
	}
	var out []NotificationRequest
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func EnqueueNotificationRequest(tx *sqlx.Tx, userID users.UserID, notificationType NotificationType, holdUntil time.Time) (_didEnqueue bool, _err error) {
	var matches []dbNotificationRequest
	if err := tx.Select(&matches, getMostRecentNotificationRequestForUserAndTypeQuery, userID, notificationType); err != nil {
		return false, err
	}
	switch {
	case len(matches) == 0:
	// no-op
	case len(matches) == 1:
		minimumElapsedTimeForType, ok := minimumElapsedTimeBetweenNotificationsByType[notificationType]
		switch {
		case !ok:
			return false, fmt.Errorf("No minimum elapsed time for type %s", notificationType)
		case minimumElapsedTimeForType == nil:
			// Do not enqueue
			return false, nil
		default:
			now := time.Now()
			earliestTimeForNextMessage := matches[0].CreatedAt.Add(*minimumElapsedTimeForType)
			if now.Before(earliestTimeForNextMessage) {
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("Got %d notification requests, but expected 1.", len(matches))
	}
	if _, err := tx.Exec(insertNotificationRequestForUserQuery, userID, notificationType, holdUntil); err != nil {
		return false, err
	}
	return true, nil
}

// Call this before sending any emails
func FulfillNotificationRequest(tx *sqlx.Tx, id NotificationRequestID) error {
	if _, err := tx.Exec(insertNotificationRequestDebounceRecord, id); err != nil {
		return err
	}
	if _, err := tx.Exec(fulfillNotificationRequestQuery, id); err != nil {
		return err
	}
	return nil
}
