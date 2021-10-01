package scheduler

import (
	"babblegraph/model/newslettersendrequests"
	"babblegraph/util/database"
	"babblegraph/util/storage"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

const cleanUpPeriod = 30 * 24 * time.Hour // 30 Days

func handleCleanupOldNewsletter(localSentryHub *sentry.Hub, s3Storage *storage.S3Storage) error {
	cleanUpDate := time.Now().Add(-1 * cleanUpPeriod)
	var sendRequests []newslettersendrequests.NewsletterSendRequest
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sendRequests, err = newslettersendrequests.GetNonDeletedSendRequestsOlderThan(tx, cleanUpDate)
		return err
	}); err != nil {
		localSentryHub.CaptureException(err)
		return err
	}
	for _, req := range sendRequests {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			switch req.PayloadStatus {
			case newslettersendrequests.PayloadStatusNoSendRequested,
				newslettersendrequests.PayloadStatusUnverifiedUser:
				return newslettersendrequests.UpdateSendRequestStatus(tx, req.ID, newslettersendrequests.PayloadStatusDeleted)
			case newslettersendrequests.PayloadStatusNeedsPreload:
				localSentryHub.CaptureException(fmt.Errorf("Got preload with ID %s that was never sent.", req.ID))
				return newslettersendrequests.UpdateSendRequestStatus(tx, req.ID, newslettersendrequests.PayloadStatusDeleted)
			case newslettersendrequests.PayloadStatusDeleted:
				return nil
			case newslettersendrequests.PayloadStatusSent:
				// no-op
			case newslettersendrequests.PayloadStatusPayloadReady:
				localSentryHub.CaptureException(fmt.Errorf("Got ready send request with ID %s that was never sent.", req.ID))
				// no-op
			}
			if err := newslettersendrequests.UpdateSendRequestStatus(tx, req.ID, newslettersendrequests.PayloadStatusDeleted); err != nil {
				return err
			}
			return s3Storage.DeleteData("prod-spaces-1", req.GetFileKey())
		}); err != nil {
			localSentryHub.CaptureException(err)
			continue
		}
		log.Println(fmt.Sprintf("Successfully deleted send request"))
	}
	return nil
}
