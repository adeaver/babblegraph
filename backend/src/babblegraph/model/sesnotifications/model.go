package sesnotifications

import (
	"time"
)

type NotificationID string

type dbNotification struct {
	ID          NotificationID `db:"_id"`
	CreatedAt   time.Time      `db:"created_at"`
	MessageBody string         `db:"message_body"`
}
