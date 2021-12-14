package scheduler

import (
	"babblegraph/model/newslettersendrequests"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/storage"
	"time"

	"github.com/jmoiron/sqlx"
)

const cleanUpPeriod = 30 * 24 * time.Hour // 30 Days

func handleCleanupOldNewsletter(c async.Context) {
	s3Storage := storage.NewS3StorageForEnvironment()
	cleanUpDate := time.Now().Add(-1 * cleanUpPeriod)
	var sendRequests []newslettersendrequests.NewsletterSendRequest
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sendRequests, err = newslettersendrequests.GetNonDeletedSendRequestsOlderThan(tx, cleanUpDate)
		return err
	}); err != nil {
		c.Errorf("Error getting newsletter send requests: %s", err.Error())
		return
	}
	for _, req := range sendRequests {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			switch req.PayloadStatus {
			case newslettersendrequests.PayloadStatusNoSendRequested,
				newslettersendrequests.PayloadStatusUnverifiedUser:
				return newslettersendrequests.UpdateSendRequestStatus(tx, req.ID, newslettersendrequests.PayloadStatusDeleted)
			case newslettersendrequests.PayloadStatusNeedsPreload:
				c.Warnf("Got preload with ID %s that was never sent.", req.ID)
				return newslettersendrequests.UpdateSendRequestStatus(tx, req.ID, newslettersendrequests.PayloadStatusDeleted)
			case newslettersendrequests.PayloadStatusDeleted:
				return nil
			case newslettersendrequests.PayloadStatusSent:
				// no-op
			case newslettersendrequests.PayloadStatusPayloadReady:
				c.Warnf("Got ready send request with ID %s that was never sent.", req.ID)
				// no-op
			}
			if err := newslettersendrequests.UpdateSendRequestStatus(tx, req.ID, newslettersendrequests.PayloadStatusDeleted); err != nil {
				return err
			}
			return s3Storage.DeleteData("prod-spaces-1", req.GetFileKey())
		}); err != nil {
			c.Errorf("Error deleting newsletter send request with ID %s: %s", req.ID, err.Error())
			continue
		}
		c.Infof("Successfully deleted send request with ID %s", req.ID)
	}
}
