package process

import (
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/newsletter"
	"babblegraph/model/newslettersendrequests"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/services/worker/newsletterprocessing"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ses"
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
				}
				if err := newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusPayloadReady); err != nil {
					return err
				}
				// if err := newslettersendrequests.UpdateSendRequestSendAtTime(tx, sendRequest.ID,
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
					BucketName:  "prod-spaces-1",
					FileName:    sendRequest.GetFileKey(),
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

func StartNewsletterFulfillmentWorkerThread(workerNumber int, newsletterProcessor *newsletterprocessing.NewsletterProcessor, errs chan error) func() {
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
		log.Println("Starting Newsletter Fulfillment Process")
		s3Storage := storage.NewS3StorageForEnvironment()
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
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
			if database.WithTx(func(tx *sqlx.Tx) error {
				user, err := users.GetUser(tx, sendRequest.UserID)
				switch {
				case err != nil:
					return err
				case user.Status != users.UserStatusVerified:
					return newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusUnverifiedUser)
				}
				log.Println(fmt.Sprintf("Found user %s", user.ID))
				if err := newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusSent); err != nil {
					return err
				}
				data, err := s3Storage.GetData("prod-spaces-1", sendRequest.GetFileKey())
				if err != nil {
					return err
				}
				log.Println(fmt.Sprintf("Found data %s", *data))
				var newsletter newsletter.Newsletter
				if err := json.Unmarshal([]byte(*data), &newsletter); err != nil {
					return err
				}
				log.Println(fmt.Sprintf("Unmarshalled %+v", newsletter))
				userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, sendRequest.UserID)
				if err != nil {
					return err
				}
				newsletterHTML, err := emailtemplates.MakeNewsletterHTML(emailtemplates.MakeNewsletterHTMLInput{
					EmailRecordID: newsletter.EmailRecordID,
					UserAccessor:  userAccessor,
					Body:          newsletter.Body,
				})
				if err != nil {
					return err
				}
				log.Println(fmt.Sprintf("Created HTML %s", *newsletterHTML))
				today := time.Now()
				subject := fmt.Sprintf("Babblegraph Newsletter - %s %d, %d", today.Month().String(), today.Day(), today.Year())
				return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
					ID:           newsletter.EmailRecordID,
					EmailAddress: user.EmailAddress,
					Body:         *newsletterHTML,
					Subject:      subject,
				})
			}); err != nil {
				log.Println(fmt.Sprintf("Got error fulfilling send request with ID %s: %s", sendRequest.ID, err.Error()))
				localHub.CaptureException(err)
				continue
			}
			log.Println(fmt.Sprintf("finished fulfilling send request with ID %s", sendRequest.ID))
		}
	}
}
