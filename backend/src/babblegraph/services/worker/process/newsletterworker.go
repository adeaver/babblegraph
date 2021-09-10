package process

import (
	"babblegraph/model/newsletter"
	"babblegraph/model/newslettersendrequests"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/services/worker/newsletterprocessing"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/storage"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

const (
	defaultPreloadWaitInterval = 1 * time.Minute
)

func getNewsletterDataBucketName() string {
	return fmt.Sprintf("worker-%s/newsletter-data", env.MustEnvironmentName().Str())
}

func StartNewsletterPreloadWorkerThread(workerNumber int, newsletterProcessor *newsletterprocessing.NewsletterProcessor, errs chan error) func() {
	return func() {
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("newsletter-preload-thread", fmt.Sprintf("init#%d", workerNumber))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Newsletter Preload Worker Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting Newsletter Preload Process")
		s3Storage := storage.NewS3StorageForEnvironment()
		for {
			sendRequest, err := newsletterProcessor.GetNextSendRequestToPreload()
			switch {
			case err != nil:
				localHub.CaptureException(err)
				continue
			case sendRequest == nil:
				log.Println("No send request available, waiting")
				time.Sleep(defaultPreloadWaitInterval)
				continue
			}
			log.Println(fmt.Sprintf("Got send request with ID %s", sendRequest.ID))
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, sendRequest.UserID)
				switch {
				case err != nil:
					return err
				case subscriptionLevel == nil:
					// no-op
				case *subscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
					*subscriptionLevel == useraccounts.SubscriptionLevelPremium:
					userScheduleForDay, err := usernewsletterschedule.LookupNewsletterDayMetadataForUserAndDay(tx, sendRequest.UserID, int(sendRequest.DateOfSend.Weekday()))
					if err != nil {
						return err
					}
					if userScheduleForDay != nil && !userScheduleForDay.IsActive {
						return newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusNoSendRequested)
					}
					if err := newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusPayloadReady); err != nil {
						return err
					}
				}
				wordsmithAccessor := newsletter.GetDefaultWordsmithAccessor()
				emailAccessor := newsletter.GetDefaultEmailAccessor(tx)
				docsAccessor := newsletter.GetDefaultDocumentsAccessor()
				userAccessor, err := newsletter.GetDefaultUserPreferencesAccessor(tx, sendRequest.UserID, sendRequest.LanguageCode, sendRequest.DateOfSend)
				if err != nil {
					return err
				}
				log.Println(fmt.Sprintf("Creating newsletter for send request with ID %s", sendRequest.ID))
				newsletter, err := newsletter.CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
				switch {
				case err != nil:
					return err
				case newsletter == nil:
					return fmt.Errorf("No send requested, but attempted to create newsletter")
				case newsletter != nil:
					// no-op
				}
				if err != nil {
					return err
				}
				newsletterBytes, err := json.Marshal(newsletter)
				if err != nil {
					return err
				}
				log.Println(fmt.Sprintf("Storing newsletter data for send request with ID %s", sendRequest.ID))
				return s3Storage.UploadData(storage.UploadDataInput{
					ContentType: storage.ContentTypeApplicationJSON,
					BucketName:  getNewsletterDataBucketName(),
					FileName:    fmt.Sprintf("%s.json", sendRequest.ID),
					Data:        string(newsletterBytes),
				})
			}); err != nil {
				log.Println(fmt.Sprintf("Got error processing send request with ID %s: %s", sendRequest.ID, err.Error()))
				localHub.CaptureException(err)
				continue
			}
			log.Println(fmt.Sprintf("finished processing send request with ID %s", sendRequest.ID))
		}
	}
}

func startNewsletterFulfillmentWorkerThread(workerNumber int, newsletterProcessor *newsletterprocessing.NewsletterProcessor, errs chan error) func() {
	return func() {
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("newsletter-fulfillment-thread", fmt.Sprintf("init#%d", workerNumber))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Newsletter Fulfillment Worker Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		for {
			sendRequest, err := newsletterProcessor.GetNextSendRequestToFulfill()
			switch {
			case err != nil:
				localHub.CaptureException(err)
				continue
			case sendRequest == nil:
				log.Println("No send request available, waiting")
				time.Sleep(defaultPreloadWaitInterval)
				continue
			}
		}
	}
}
