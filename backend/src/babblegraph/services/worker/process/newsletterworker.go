package process

import (
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/newsletter"
	"babblegraph/model/newslettersendrequests"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/services/worker/newsletterprocessing"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
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
				userNewsletterPreferences, err := usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(c, tx, sendRequest.UserID, sendRequest.LanguageCode, ptr.Time(dateOfSendUTCMidnight))
				if err != nil {
					return err
				}
				if !userNewsletterPreferences.Schedule.IsSendRequested(dateOfSendUTCMidnight.Weekday()) {
					return newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusNoSendRequested)
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
				podcastAccessor, err := newsletter.GetDefaultPodcastAccessor(c, tx, sendRequest.LanguageCode, sendRequest.UserID)
				if err != nil {
					return err
				}
				advertisementAccessor, err := newsletter.GetDefaultAdvertisementAccessor(tx, sendRequest.UserID, sendRequest.LanguageCode)
				if err != nil {
					return err
				}
				newsletter, err := newsletter.CreateNewsletterVersion2(c, dateOfSendUTCMidnight, newsletter.CreateNewsletterVersion2Input{
					WordsmithAccessor:     wordsmithAccessor,
					EmailAccessor:         emailAccessor,
					UserAccessor:          userAccessor,
					DocsAccessor:          docsAccessor,
					ContentAccessor:       contentAccessor,
					PodcastAccessor:       podcastAccessor,
					AdvertisementAccessor: advertisementAccessor,
				})
				switch {
				case err != nil:
					return err
				case newsletter == nil:
					return fmt.Errorf("No send requested, but attempted to create newsletter")
				case newsletter != nil:
					// no-op
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
				isOnBounceQuarantine, err := email.IsUserOnBounceQuarantine(tx, user.ID)
				switch {
				case err != nil:
					return err
				case isOnBounceQuarantine:
					c.Infof("User is on bounce quarantine, skipping this newsletter")
					return nil
				}
				userNewsletterPreferences, err := usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(c, tx, sendRequest.UserID, sendRequest.LanguageCode, ptr.Time(dateOfSendUTCMidnight))
				if err != nil {
					return err
				}
				data, err := s3Storage.GetData("prod-spaces-1", sendRequest.GetFileKey())
				if err != nil {
					return err
				}
				c.Debugf("Found data %s", *data)
				userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, sendRequest.UserID)
				if err != nil {
					return err
				}
				var newsletterHTML *string
				var edition newsletter.Newsletter
				var emailRecordID email.ID
				// First we try to unmarshal newsletter
				if err := json.Unmarshal([]byte(*data), &edition); err != nil {
					return err
				}
				switch {
				case len(edition.Body.ReinforcementLink) == 0:
					var newsletterVersion2 newsletter.NewsletterVersion2
					if err := json.Unmarshal([]byte(*data), &newsletterVersion2); err != nil {
						return err
					}
					c.Debugf("ID %s was version 2", sendRequest.ID)
					c.Debugf("Unmarshalled %+v", newsletterVersion2)
					emailRecordID = newsletterVersion2.EmailRecordID
					newsletterHTML, err = emailtemplates.MakeNewsletterVersion2HTML(emailtemplates.MakeNewsletterVersion2HTMLInput{
						EmailRecordID: newsletterVersion2.EmailRecordID,
						UserAccessor:  userAccessor,
						Body:          newsletterVersion2.Body,
					})
					if err != nil {
						return err
					}
				default:
					c.Debugf("Unmarshalled %+v", edition)
					emailRecordID = edition.EmailRecordID
					newsletterHTML, err = emailtemplates.MakeNewsletterHTML(emailtemplates.MakeNewsletterHTMLInput{
						EmailRecordID: edition.EmailRecordID,
						UserAccessor:  userAccessor,
						Body:          edition.Body,
					})
					if err != nil {
						return err
					}
				}
				c.Debugf("Created HTML %s", *newsletterHTML)
				today, err := userNewsletterPreferences.Schedule.ConvertUTCTimeToUserDate(c, dateOfSendUTCMidnight)
				if err != nil {
					return err
				}
				subject := fmt.Sprintf("Babblegraph Newsletter - %s %d, %d", today.Month().String(), today.Day(), today.Year())
				return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
					ID:           emailRecordID,
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
