package sesnotifications

import (
	"babblegraph/util/ses"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

const insertSESNotificationQuery = "INSERT INTO ses_notifications (message_body) VALUES ($1)"

func InsertSESNotification(tx *sqlx.Tx, messageBody ses.Notification) error {
	bodyAsJSON, err := json.Marshal(messageBody)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(insertSESNotificationQuery, string(bodyAsJSON)); err != nil {
		return err
	}
	return nil
}
