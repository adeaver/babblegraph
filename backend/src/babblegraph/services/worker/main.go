package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/domains"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/scheduler"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/elastic"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	linkProcessor, err := linkprocessing.CreateLinkProcessor()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := linkProcessor.AddURLs(domains.GetSeedURLs()); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	for i := 0; i < 3; i++ {
		workerThread := startWorkerThread(linkProcessor, errs)
		go workerThread()
	}
	schedulerErrs := make(chan error, 1)
	if err := scheduler.StartScheduler(schedulerErrs); err != nil {
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
			if err, ok := x.(error); ok {
				errs <- err
				debug.PrintStack()
			}
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
			log.Println(fmt.Sprintf("Processing URL %s with identifier %s", u, link.URLIdentifier))
			parsedHTMLPage, err := ingesthtml.ProcessURL(u, domain)
			if err != nil {
				log.Println(fmt.Sprintf("Got error ingesting html for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
			languageLabel := deref.String(parsedHTMLPage.Language, "")
			languageCode := wordsmith.LookupLanguageCodeForLanguageLabel(languageLabel)
			if languageCode == nil {
				log.Println(fmt.Sprintf("URL %s has unsupported language code: %s", u, languageLabel))
				// This is only allowed because we're restricting domains
				// If domains are ever non-restricted, this needs to be removed
				// and made more robust
				languageCode = wordsmith.LanguageCodeSpanish.Ptr()
			}
			log.Println(fmt.Sprintf("Got language code %s for label %s on URL %s. Processing...", languageCode.Str(), languageLabel, u))
			if err := linkProcessor.AddURLs(parsedHTMLPage.Links); err != nil {
				log.Println(fmt.Sprintf("Error saving urls %+v for url %s: %s", parsedHTMLPage.Links, u, err.Error()))
				continue
			}
			if strings.ToLower(deref.String(parsedHTMLPage.PageType, "")) != "article" {
				log.Println(fmt.Sprintf("URL %s is not an article. Continuing...", u))
				continue
			}
			log.Println(fmt.Sprintf("Processing text for url %s", u))
			textMetadata, err := textprocessing.ProcessText(parsedHTMLPage.BodyText, *languageCode)
			if err != nil {
				log.Println(fmt.Sprintf("Got error processing text for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
			log.Println(fmt.Sprintf("Indexing text for URL %s", u))
			err = indexing.IndexDocument(indexing.IndexDocumentInput{
				ParsedHTMLPage:  *parsedHTMLPage,
				TextMetadata:    *textMetadata,
				LanguageCode:    *languageCode,
				DocumentVersion: documents.CurrentDocumentVersion,
				URL:             urlparser.MustParseURL(u),
			})
			if err != nil {
				log.Println(fmt.Sprintf("Got error indexing document for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
		}
	}
}
