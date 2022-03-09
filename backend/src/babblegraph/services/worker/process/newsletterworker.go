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
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"babblegraph/util/storage"
	"babblegraph/util/timeutils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	defaultPreloadWaitInterval = 1 * time.Minute
)

func StartNewsletterPreloadWorkerThread(newsletterProcessor *newsletterprocessing.NewsletterProcessor) func(c async.Context) {
	return func(c async.Context) {
		c.Infof("Starting Newsletter Preload Process")
		s3Storage := storage.NewS3StorageForEnvironment()
		for {
			sendRequest, err := newsletterProcessor.GetNextSendRequestToPreload(c)
			switch {
			case err != nil:
				c.Errorf("Error getting send request: %s", err.Error())
				continue
			case sendRequest == nil:
				c.Infof("No send request available, waiting")
				time.Sleep(defaultPreloadWaitInterval)
				continue
			}
			c.Infof("Got send request with ID %s", sendRequest.ID)
			dateOfSendUTCMidnight := timeutils.ConvertToMidnight(sendRequest.DateOfSend.UTC())
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, sendRequest.UserID)
				switch {
				case err != nil:
					return err
				case subscriptionLevel == nil:
					// no-op
				case *subscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
					*subscriptionLevel == useraccounts.SubscriptionLevelPremium:
					userNewsletterSchedule, err := usernewsletterschedule.GetUserNewsletterScheduleForUTCMidnight(c, tx, usernewsletterschedule.GetUserNewsletterScheduleForUTCMidnightInput{
						UserID:           sendRequest.UserID,
						LanguageCode:     sendRequest.LanguageCode,
						DayAtUTCMidnight: dateOfSendUTCMidnight,
					})
					if err != nil {
						return err
					}
					if !userNewsletterSchedule.IsSendRequested() {
						return newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusNoSendRequested)
					}
				}
				if err := newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusPayloadReady); err != nil {
					return err
				}
				wordsmithAccessor := newsletter.GetDefaultWordsmithAccessor()
				emailAccessor := newsletter.GetDefaultEmailAccessor(tx)
				docsAccessor := newsletter.GetDefaultDocumentsAccessor()
				userAccessor, err := newsletter.GetDefaultUserPreferencesAccessor(c, tx, sendRequest.UserID, sendRequest.LanguageCode, dateOfSendUTCMidnight)
				if err != nil {
					return err
				}
				contentAccessor, err := newsletter.GetDefaultContentAccessor(tx, sendRequest.LanguageCode)
				if err != nil {
					return err
				}
				c.Infof("Creating newsletter for send request with ID %s", sendRequest.ID)
				newsletter, err := newsletter.CreateNewsletter(c, newsletter.CreateNewsletterInput{
					WordsmithAccessor: wordsmithAccessor,
					EmailAccess:       emailAccessor,
					UserAccessor:      userAccessor,
					DocsAccessor:      docsAccessor,
					ContentAccessor:   contentAccessor,
				})
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
				c.Infof("Storing newsletter data for send request with ID %s", sendRequest.ID)
				return s3Storage.UploadData(storage.UploadDataInput{
					ContentType: storage.ContentTypeApplicationJSON,
					BucketName:  "prod-spaces-1",
					FileName:    sendRequest.GetFileKey(),
					Data:        string(newsletterBytes),
				})
			}); err != nil {
				c.Errorf("Got error processing send request with ID %s: %s", sendRequest.ID, err.Error())
				continue
			}
			c.Infof("finished processing send request with ID %s", sendRequest.ID)
		}
	}
}

func StartNewsletterFulfillmentWorkerThread(newsletterProcessor *newsletterprocessing.NewsletterProcessor) func(c async.Context) {
	return func(c async.Context) {
		c.Infof("Starting Newsletter Fulfillment Process")
		s3Storage := storage.NewS3StorageForEnvironment()
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		for {
			sendRequest, err := newsletterProcessor.GetNextSendRequestToFulfill(c)
			switch {
			case err != nil:
				c.Errorf("Error getting fulfillment request: %s", err.Error())
				continue
			case sendRequest == nil:
				c.Infof("No send request available, waiting")
				time.Sleep(defaultPreloadWaitInterval)
				continue
			}
			dateOfSendUTCMidnight := timeutils.ConvertToMidnight(sendRequest.DateOfSend.UTC())
			if database.WithTx(func(tx *sqlx.Tx) error {
				user, err := users.GetUser(tx, sendRequest.UserID)
				switch {
				case err != nil:
					return err
				case user.Status != users.UserStatusVerified:
					return newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusUnverifiedUser)
				}
				c.Infof("Found user %s", user.ID)
				if err := newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusSent); err != nil {
					return err
				}
				userNewsletterSchedule, err := usernewsletterschedule.GetUserNewsletterScheduleForUTCMidnight(c, tx, usernewsletterschedule.GetUserNewsletterScheduleForUTCMidnightInput{
					UserID:           sendRequest.UserID,
					LanguageCode:     sendRequest.LanguageCode,
					DayAtUTCMidnight: dateOfSendUTCMidnight,
				})
				if err != nil {
					return err
				}
				data, err := s3Storage.GetData("prod-spaces-1", sendRequest.GetFileKey())
				if err != nil {
					return err
				}
				c.Debugf("Found data %s", *data)
				var newsletter newsletter.Newsletter
				if err := json.Unmarshal([]byte(*data), &newsletter); err != nil {
					return err
				}
				c.Debugf("Unmarshalled %+v", newsletter)
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
				c.Debugf("Created HTML %s", *newsletterHTML)
				today := userNewsletterSchedule.GetSendTimeInUserTimezone()
				subject := fmt.Sprintf("Babblegraph Newsletter - %s %d, %d", today.Month().String(), today.Day(), today.Year())
				return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
					ID:           newsletter.EmailRecordID,
					EmailAddress: user.EmailAddress,
					Body:         *newsletterHTML,
					Subject:      subject,
				})
			}); err != nil {
				c.Errorf("Got error fulfilling send request with ID %s: %s", sendRequest.ID, err.Error())
				continue
			}
			c.Infof("finished fulfilling send request with ID %s", sendRequest.ID)
		}
	}
}
