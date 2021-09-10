package newslettersendrequests

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"time"
)

type ID string

type NewsletterSendRequest struct {
	ID            ID
	UserID        users.UserID
	LanguageCode  wordsmith.LanguageCode
	DateOfSend    string
	PayloadStatus PayloadStatus
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

func (d dbNewsletterSendRequest) ToNonDB() NewsletterSendRequest {
	return NewsletterSendRequest{
		ID:            d.ID,
		UserID:        d.UserID,
		LanguageCode:  d.LanguageCode,
		DateOfSend:    d.DateOfSend,
		PayloadStatus: d.PayloadStatus,
	}
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
	PayloadStatusSent            PayloadStatus = "sent"
	PayloadStatusDeleted         PayloadStatus = "deleted"
)
