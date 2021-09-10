package newslettersendrequests

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getSendRequestsForUsersForDayQuery = "SELECT * FROM newsletter_send_requests WHERE date_of_send = '%s' AND language_code = '%s' AND user_id IN (?)"
	insertSendRequestForUserQuery      = `INSERT INTO
        newsletter_send_requests
    (
        _id,
        user_id,
        language_code,
        date_of_send,
        payload_status
    ) VALUES (
        $1, $2, $3, $4, $5
    )`

	updateSendRequestStatusQuery = "UPDATE newsletter_send_requests SET payload_status = $1, lat_modified_at=timezone('utc', now()) WHERE _id = $2"
	insertDebounceRecordQuery    = "INSERT INTO newsletter_send_request_debounce_records (newsletter_send_request_id, to_payload_status) VALUES ($1, $2)"
)

func GetOrCreateSendRequestsForUsersForDay(tx *sqlx.Tx, userIDs []users.UserID, languageCode wordsmith.LanguageCode, day time.Time) ([]NewsletterSendRequest, error) {
	dateOfSendString := getDateOfSendForTime(day)
	query, args, err := sqlx.In(fmt.Sprintf(getSendRequestsForUsersForDayQuery, dateOfSendString, languageCode), userIDs)
	if err != nil {
		return nil, nil
	}
	sql := tx.Rebind(query)
	var matches []dbNewsletterSendRequest
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, nil
	}
	var out []NewsletterSendRequest
	usersWithSendRequests := make(map[users.UserID]bool)
	for _, m := range matches {
		usersWithSendRequests[m.UserID] = true
		out = append(out, m.ToNonDB())
	}
	for _, u := range userIDs {
		if _, ok := usersWithSendRequests[u]; !ok {
			id := makeSendRequestID(u, languageCode, dateOfSendString)
			if _, err := tx.Exec(insertSendRequestForUserQuery, u, languageCode, dateOfSendString, PayloadStatusNeedsPreload); err != nil {
				return nil, err
			}
			out = append(out, NewsletterSendRequest{
				ID:            id,
				UserID:        u,
				DateOfSend:    dateOfSendString,
				PayloadStatus: PayloadStatusNeedsPreload,
			})
		}
	}
	return out, nil
}

func UpdateSendRequestStatus(tx *sqlx.Tx, id ID, newStatus PayloadStatus) error {
	if _, err := tx.Exec(insertDebounceRecordQuery, id, newStatus); err != nil {
		return err
	}
	if _, err := tx.Exec(updateSendRequestStatusQuery, newStatus, id); err != nil {
		return err
	}
	return nil
}
