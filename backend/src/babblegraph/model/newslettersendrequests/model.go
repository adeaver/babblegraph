package newslettersendrequests

import (
	"babblegraph/model/users"
	"babblegraph/util/deref"
	"babblegraph/util/env"
	"babblegraph/wordsmith"
	"fmt"
	"time"
)

type ID string

type NewsletterSendRequest struct {
	ID            ID
	UserID        users.UserID
	LanguageCode  wordsmith.LanguageCode
	DateOfSend    time.Time
	PayloadStatus PayloadStatus
}

func (req *NewsletterSendRequest) GetFileKey() string {
	return fmt.Sprintf("worker-%s/newsletter-data/%s.json", env.MustEnvironmentName().Str(), req.ID)
}

type dbNewsletterSendRequest struct {
	ID                        ID                     `db:"_id"`
	CreatedAt                 time.Time              `db:"created_at"`
	LastModifiedAt            time.Time              `db:"last_modified_at"`
	UserID                    users.UserID           `db:"user_id"`
	LanguageCode              wordsmith.LanguageCode `db:"language_code"`
	DateOfSend                string                 `db:"date_of_send"`
	HourToSendIndexUTC        *int64                 `db:"hour_to_send_index_utc"`
	QuarterHourToSendIndexUTC *int64                 `db:"quarter_hour_to_send_index_utc"`
	PayloadStatus             PayloadStatus          `db:"payload_status"`
}

func (d dbNewsletterSendRequest) ToNonDB() (*NewsletterSendRequest, error) {
	utcDate, err := getUTCMidnightDateOfSend(d.DateOfSend)
	if err != nil {
		return nil, err
	}
	hourToSend := int(deref.Int64(d.HourToSendIndexUTC, 0))
	minuteToSend := int(deref.Int64(d.QuarterHourToSendIndexUTC, 0) * 15)
	sendAtTime := time.Date(utcDate.Year(), utcDate.Month(), utcDate.Day(), hourToSend, minuteToSend, 0, 0, time.UTC)
	return &NewsletterSendRequest{
		ID:            d.ID,
		UserID:        d.UserID,
		LanguageCode:  d.LanguageCode,
		DateOfSend:    sendAtTime,
		PayloadStatus: d.PayloadStatus,
	}, nil
}

type debounceRecordID string

type dbNewsletterSendRequestProcessingDebounceRecord struct {
	ID                      debounceRecordID `db:"_id"`
	CreatedAt               time.Time        `db:"created_at"`
	LastModifiedAt          time.Time        `db:"last_modified_at"`
	NewsletterSendRequestID ID               `db:"newsletter_send_request_id"`
	ToPayloadStatus         PayloadStatus    `db:"to_payload_status"`
}

type PayloadStatus string

const (
	PayloadStatusNeedsPreload    PayloadStatus = "needs-preload"
	PayloadStatusNoSendRequested PayloadStatus = "no-send-requested"
	PayloadStatusPayloadReady    PayloadStatus = "payload-ready"
	PayloadStatusUnverifiedUser  PayloadStatus = "user-not-verified"
	PayloadStatusSent            PayloadStatus = "sent"
	PayloadStatusDeleted         PayloadStatus = "deleted"
)
