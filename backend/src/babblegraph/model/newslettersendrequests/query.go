package newslettersendrequests

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getSendRequestsForUsersForDayQuery      = "SELECT * FROM newsletter_send_requests WHERE date_of_send = '%s' AND language_code = '%s' AND user_id IN (?)"
	getSendRequestsOlderThanWithStatusQuery = "SELECT * FROM newsletter_send_requests WHERE created_at <= $1 AND payload_status != $2"
	insertSendRequestForUserQuery           = `INSERT INTO
        newsletter_send_requests
    (
        _id,
        user_id,
        language_code,
        date_of_send,
        payload_status,
        hour_to_send_index_utc,
        quarter_hour_to_send_index_utc
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )`

	updateSendRequestSendAtTimeQuery = "UPDATE newsletter_send_requests SET hour_to_send_index_utc=$1, quarter_hour_to_send_index_utc=$2 WHERE _id = $3"

	updateSendRequestStatusQuery = "UPDATE newsletter_send_requests SET payload_status = $1, last_modified_at=timezone('utc', now()) WHERE _id = $2"
	insertDebounceRecordQuery    = "INSERT INTO newsletter_send_request_debounce_records (newsletter_send_request_id, to_payload_status) VALUES ($1, $2)"

	defaultUSEasternHourToSend        = 7
	defaultUSEasternQuarterHourToSend = 0
)

func GetOrCreateSendRequestsForUsersForDay(tx *sqlx.Tx, userIDs []users.UserID, languageCode wordsmith.LanguageCode, day time.Time) ([]NewsletterSendRequest, error) {
	usEastern, err := time.LoadLocation("US/Eastern")
	if err != nil {
		return nil, err
	}
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
		req, err := m.ToNonDB()
		if err != nil {
			return nil, err
		}
		out = append(out, *req)
	}
	for _, u := range userIDs {
		if _, ok := usersWithSendRequests[u]; !ok {
			dateOfSend := time.Date(day.Year(), day.Month(), day.Day(), defaultUSEasternHourToSend, defaultUSEasternQuarterHourToSend, 0, 0, usEastern).UTC()
			id := makeSendRequestID(u, languageCode, dateOfSendString)
			if _, err := tx.Exec(insertSendRequestForUserQuery, id, u, languageCode, dateOfSendString, PayloadStatusNeedsPreload, dateOfSend.Hour(), dateOfSend.Minute()/15); err != nil {
				return nil, err
			}
			out = append(out, NewsletterSendRequest{
				ID:            id,
				UserID:        u,
				LanguageCode:  languageCode,
				DateOfSend:    dateOfSend,
				PayloadStatus: PayloadStatusNeedsPreload,
			})
		}
	}
	return out, nil
}

func GetNonDeletedSendRequestsOlderThan(tx *sqlx.Tx, t time.Time) ([]NewsletterSendRequest, error) {
	var matches []dbNewsletterSendRequest
	if err := tx.Select(&matches, getSendRequestsOlderThanWithStatusQuery, t, PayloadStatusDeleted); err != nil {
		return nil, err
	}
	var out []NewsletterSendRequest
	for _, m := range matches {
		req, err := m.ToNonDB()
		if err != nil {
			return nil, err
		}
		out = append(out, *req)
	}
	return out, nil
}

func UpdateSendRequestSendAtTime(tx *sqlx.Tx, id ID, sendAtTime time.Time) error {
	utcSendTime := sendAtTime.UTC()
	sendAtHourIndexUTC := utcSendTime.Hour()
	sendAtQuarterHourIndexUTC := utcSendTime.Minute() / 15
	if _, err := tx.Exec(updateSendRequestSendAtTimeQuery, sendAtHourIndexUTC, sendAtQuarterHourIndexUTC, id); err != nil {
		return err
	}
	return nil
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
