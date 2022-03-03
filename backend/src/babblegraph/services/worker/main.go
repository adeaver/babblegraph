package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/contentingestion"
	"babblegraph/services/worker/newsletterprocessing"
	"babblegraph/services/worker/process"
	"babblegraph/services/worker/scheduler"
	"babblegraph/util/async"
	"babblegraph/util/bglog"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/util/env"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	numIngestWorkerThreads          = 5
	numNewsletterPreloadThreads     = 2
	numNewsletterFulfillmentThreads = 1
)

func main() {
	bglog.InitLogger()
	currentEnvironmentName := env.MustEnvironmentName()
	if err := setupDatabases(); err != nil {
		bglog.Fatalf("Error initializing databases: %s", err.Error())
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.MustEnvironmentVariable("SENTRY_DSN"),
		Environment: currentEnvironmentName.Str(),
	}); err != nil {
		bglog.Fatalf("Error initializing sentry: %s", err.Error())
	}
	defer sentry.Flush(2 * time.Second)
	ingestErrs := make(chan error, 1)
	if currentEnvironmentName != env.EnvironmentLocalTestEmail {
		async.WithContext(ingestErrs, "link-processor", contentingestion.StartIngestion()).Start()
	}
	schedulerErrs := make(chan error, 1)
	if err := scheduler.StartScheduler(schedulerErrs); err != nil {
		bglog.Fatalf("Error initializing scheduler: %s", err.Error())
	}
	newsletterProcessor, err := newsletterprocessing.CreateNewsletterProcessor(ctx.GetDefaultLogContext())
	if err != nil {
		bglog.Fatalf("Error initializing newsletter processor: %s", err.Error())
	}
	preloadNewsletterErrs := make(chan error, 1)
	preloadWorkerNum := 0
	for i := 0; i < numNewsletterPreloadThreads; i++ {
		async.WithContext(preloadNewsletterErrs, fmt.Sprintf("preload-worker-%d", preloadWorkerNum), process.StartNewsletterPreloadWorkerThread(newsletterProcessor)).Start()
		preloadWorkerNum++
	}
	fulfillNewsletterErrs := make(chan error, 1)
	fulfillWorkerNum := 0
	for i := 0; i < numNewsletterFulfillmentThreads; i++ {
		async.WithContext(fulfillNewsletterErrs, fmt.Sprintf("fulfillment-worker-%d", fulfillWorkerNum), process.StartNewsletterFulfillmentWorkerThread(newsletterProcessor)).Start()
		fulfillWorkerNum++
	}
	for {
		select {
		case _ = <-ingestErrs:
			panic("Error on ingest, forcing restart")
		case err = <-schedulerErrs:
			bglog.Infof("Saw panic: %s in scheduler.", err.Error())
		case _ = <-preloadNewsletterErrs:
			async.WithContext(preloadNewsletterErrs, fmt.Sprintf("preload-worker-%d", preloadWorkerNum), process.StartNewsletterPreloadWorkerThread(newsletterProcessor)).Start()
			preloadWorkerNum++
		case _ = <-fulfillNewsletterErrs:
			async.WithContext(fulfillNewsletterErrs, fmt.Sprintf("fulfillment-worker-%d", fulfillWorkerNum), process.StartNewsletterFulfillmentWorkerThread(newsletterProcessor)).Start()
			fulfillWorkerNum++
		}
	}
}

func setupDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("Error setting up main-db: %s", err.Error())
	}
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up wordsmith: %s", err.Error())
	}
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up elasticsearch: %s", err.Error())
	}
	if err := documents.CreateDocumentMappings(); err != nil {
		return fmt.Errorf("Error setting up documents: %s", err.Error())
	}
	return nil
}
