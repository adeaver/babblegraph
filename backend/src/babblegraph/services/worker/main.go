package main

import (
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/newsletterprocessing"
	"babblegraph/services/worker/process"
	"babblegraph/services/worker/scheduler"
	"babblegraph/util/bglog"
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
	if err := setupDatabases(); err != nil {
		bglog.Fatalf("Error initializing databases: %s", err.Error())
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.MustEnvironmentVariable("SENTRY_DSN"),
		Environment: env.MustEnvironmentName().Str(),
	}); err != nil {
		bglog.Fatalf("Error initializing sentry: %s", err.Error())
	}
	defer sentry.Flush(2 * time.Second)
	linkProcessor, err := linkprocessing.CreateLinkProcessor()
	if err != nil {
		bglog.Fatalf("Error initializing link processor: %s", err.Error())
	}
	for u, topics := range domains.GetSeedURLs() {
		if err := linkProcessor.AddURLs([]string{u}, topics); err != nil {
			bglog.Fatalf("Error initializing link processor seed domains: %s", err.Error())
		}
	}
	currentEnvironmentName := env.MustEnvironmentName()
	workerNum := 0
	ingestErrs := make(chan error, 1)
	if currentEnvironmentName != env.EnvironmentLocalTestEmail {
		for i := 0; i < numIngestWorkerThreads; i++ {
			workerThread := process.StartIngestWorkerThread(workerNum, linkProcessor, ingestErrs)
			go workerThread()
			workerNum++
		}
	}
	schedulerErrs := make(chan error, 1)
	if err := scheduler.StartScheduler(linkProcessor, schedulerErrs); err != nil {
		bglog.Fatalf("Error initializing scheduler: %s", err.Error())
	}
	newsletterProcessor, err := newsletterprocessing.CreateNewsletterProcessor()
	if err != nil {
		bglog.Fatalf("Error initializing newsletter processor: %s", err.Error())
	}
	preloadNewsletterErrs := make(chan error, 1)
	preloadWorkerNum := 0
	for i := 0; i < numNewsletterPreloadThreads; i++ {
		preloadThread := process.StartNewsletterPreloadWorkerThread(preloadWorkerNum, newsletterProcessor, preloadNewsletterErrs)
		go preloadThread()
		preloadWorkerNum++
	}
	fulfillNewsletterErrs := make(chan error, 1)
	fulfillWorkerNum := 0
	for i := 0; i < numNewsletterFulfillmentThreads; i++ {
		fulfillThread := process.StartNewsletterFulfillmentWorkerThread(fulfillWorkerNum, newsletterProcessor, fulfillNewsletterErrs)
		go fulfillThread()
		fulfillWorkerNum++
	}
	for {
		select {
		case err := <-ingestErrs:
			bglog.Infof("Saw panic: %s. Starting new worker thread.", err.Error())
			if currentEnvironmentName != env.EnvironmentLocalTestEmail {
				workerThread := process.StartIngestWorkerThread(workerNum, linkProcessor, ingestErrs)
				go workerThread()
				workerNum++
			}
		case err := <-schedulerErrs:
			bglog.Infof("Saw panic: %s in scheduler.", err.Error())
		case err := <-preloadNewsletterErrs:
			bglog.Infof("Saw panic: %s. Starting new newsletter preload thread.", err.Error())
			preloadThread := process.StartNewsletterPreloadWorkerThread(preloadWorkerNum, newsletterProcessor, preloadNewsletterErrs)
			go preloadThread()
			preloadWorkerNum++
		case err := <-fulfillNewsletterErrs:
			bglog.Infof("Saw panic: %s. Starting new newsletter fulfillment thread.", err.Error())
			fulfillThread := process.StartNewsletterFulfillmentWorkerThread(fulfillWorkerNum, newsletterProcessor, fulfillNewsletterErrs)
			go fulfillThread()
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
