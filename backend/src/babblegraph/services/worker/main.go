package main

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/scheduler"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/elastic"
	"babblegraph/util/env"
	"babblegraph/util/opengraph"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
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
	errs := make(chan error, 1)
	for i := 0; i < 3; i++ {
		workerThread := startWorkerThread(linkProcessor, errs)
		go workerThread()
	}
	schedulerErrs := make(chan error, 1)
	if err := scheduler.StartScheduler(linkProcessor, schedulerErrs); err != nil {
		log.Fatal(err.Error())
	}
	for {
		select {
		case err := <-errs:
			log.Println(fmt.Sprintf("Saw panic: %s. Starting new worker thread.", err.Error()))
			workerThread := startWorkerThread(linkProcessor, errs)
			go workerThread()
		case err := <-schedulerErrs:
			log.Println(fmt.Sprintf("Saw panic: %s in scheduler.", err.Error()))
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
	return nil
}

func startWorkerThread(linkProcessor *linkprocessing.LinkProcessor, errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			debug.PrintStack()
			err := fmt.Errorf("Encountered worker panic: %v\n", x)
			errs <- err
		}()
		for {
			var u, domain string
			link, waitTime, err := linkProcessor.GetLink()
			switch {
			case err != nil:
				log.Println(fmt.Sprintf("Error getting link... %s", err.Error()))
				continue
			case waitTime != nil:
				log.Println("No link available. Sleeping...")
				time.Sleep(*waitTime)
				continue
			case link != nil:
				u = link.URL
				domain = link.Domain
			default:
				log.Println("No error, but no wait time. Continuing...")
				continue
			}
			if p := urlparser.ParseURL(u); p != nil && domains.IsSeedURL(*p) {
				log.Println(fmt.Sprintf("Received url %s, which is a seed url. Skipping...", u))
				continue
			}
			log.Println(fmt.Sprintf("Processing URL %s with identifier %s", u, link.URLIdentifier))
			parsedHTMLPage, err := ingesthtml.ProcessURL(u, domain)
			if err != nil {
				log.Println(fmt.Sprintf("Got error ingesting html for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
			domainMetadata, err := domains.GetDomainMetadata(domain)
			if err != nil {
				log.Println(fmt.Sprintf("Got error getting metadata for domain %s on url %s: %s. Continuing...", domain, u, err.Error()))
				continue
			}
			languageCode := domainMetadata.LanguageCode
			if err := linkProcessor.AddURLs(parsedHTMLPage.Links, domainMetadata.Topics); err != nil {
				log.Println(fmt.Sprintf("Error saving urls %+v for url %s: %s", parsedHTMLPage.Links, u, err.Error()))
				continue
			}
			if strings.ToLower(deref.String(parsedHTMLPage.PageType, "")) != "article" {
				log.Println(fmt.Sprintf("URL %s is not an article. Continuing...", u))
				continue
			}
			log.Println(fmt.Sprintf("Processing text for url %s", u))
			var description *string
			if d, ok := parsedHTMLPage.Metadata[opengraph.DescriptionTag.Str()]; ok {
				description = ptr.String(d)
			}
			textMetadata, err := textprocessing.ProcessText(textprocessing.ProcessTextInput{
				BodyText:     parsedHTMLPage.BodyText,
				Description:  description,
				LanguageCode: languageCode,
			})
			if err != nil {
				log.Println(fmt.Sprintf("Got error processing text for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
			var topicsForURL []contenttopics.ContentTopic
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				topicsForURL, err = contenttopics.GetTopicsForURL(tx, u)
				return err
			}); err != nil {
				log.Println(fmt.Sprintf("Error getting topics for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
			log.Println(fmt.Sprintf("Indexing text for URL %s", u))
			err = indexing.IndexDocument(indexing.IndexDocumentInput{
				ParsedHTMLPage:  *parsedHTMLPage,
				TextMetadata:    *textMetadata,
				LanguageCode:    languageCode,
				DocumentVersion: documents.CurrentDocumentVersion,
				URL:             urlparser.MustParseURL(u),
				TopicsForURL:    topicsForURL,
			})
			if err != nil {
				log.Println(fmt.Sprintf("Got error indexing document for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
		}
	}
}
