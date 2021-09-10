package main

import (
	"babblegraph/model/newsletter"
	"babblegraph/model/newslettersendrequests"
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

func startNewsletterPreloadWorkerThread(workerNumber int, newsletterProcessor *newsletterprocessing.NewsletterProcessor, errs chan error) func() {
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
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				wordsmithAccessor := newsletter.GetDefaultWordsmithAccessor()
				emailAccessor := newsletter.GetDefaultEmailAccessor(tx)
				docsAccessor := newsletter.GetDefaultDocumentsAccessor()
				userAccessor, err := newsletter.GetDefaultUserPreferencesAccessor(tx, sendRequest.UserID, sendRequest.LanguageCode)
				if err != nil {
					return err
				}
				newsletter, err := newsletter.CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
				switch {
				case err != nil:
					return err
				case newsletter == nil:
					return newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusNoSendRequested)
				case newsletter != nil:
					err := newslettersendrequests.UpdateSendRequestStatus(tx, sendRequest.ID, newslettersendrequests.PayloadStatusPayloadReady)
					if err != nil {
						return err
					}
				}
				if err != nil {
					return err
				}
				newsletterBytes, err := json.Marshal(newsletter)
				if err != nil {
					return err
				}
				return s3Storage.UploadData(getNewsletterDataBucketName(), fmt.Sprintf("%s.json", sendRequest.ID), string(newsletterBytes))
			}); err != nil {

			}
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
