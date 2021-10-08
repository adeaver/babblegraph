package main

import (
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/newsletterprocessing"
	"babblegraph/services/worker/process"
	"babblegraph/services/worker/scheduler"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/util/env"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	numIngestWorkerThreads          = 5
	numNewsletterPreloadThreads     = 2
	numNewsletterFulfillmentThreads = 1
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.MustEnvironmentVariable("SENTRY_DSN"),
		Environment: env.MustEnvironmentName().Str(),
	}); err != nil {
		log.Fatal(err.Error())
	}
	defer sentry.Flush(2 * time.Second)

	linkProcessor, err := linkprocessing.CreateLinkProcessor()
	if err != nil {
		log.Fatal(err.Error())
	}
	for u, topics := range domains.GetSeedURLs() {
		if err := linkProcessor.AddURLs([]string{u}, topics); err != nil {
			log.Fatal(err.Error())
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
		log.Fatal(err.Error())
	}
	newsletterProcessor, err := newsletterprocessing.CreateNewsletterProcessor()
	if err != nil {
		log.Fatal(err.Error())
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
			log.Println(fmt.Sprintf("Saw panic: %s. Starting new worker thread.", err.Error()))
			if currentEnvironmentName != env.EnvironmentLocalTestEmail {
				workerThread := process.StartIngestWorkerThread(workerNum, linkProcessor, ingestErrs)
				go workerThread()
				workerNum++
			}
		case err := <-schedulerErrs:
			log.Println(fmt.Sprintf("Saw panic: %s in scheduler.", err.Error()))
		case err := <-preloadNewsletterErrs:
			log.Println(fmt.Sprintf("Saw panic: %s. Starting new newsletter preload thread.", err.Error()))
			preloadThread := process.StartNewsletterPreloadWorkerThread(preloadWorkerNum, newsletterProcessor, preloadNewsletterErrs)
			go preloadThread()
			preloadWorkerNum++
		case err := <-fulfillNewsletterErrs:
			log.Println(fmt.Sprintf("Saw panic: %s. Starting new newsletter fulfillment thread.", err.Error()))
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
