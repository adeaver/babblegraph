package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
	"fmt"
	"log"
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
	if err := linkProcessor.AddURLs([]string{
		"elmundo.es",
		"https://cnnespanol.cnn.com/",
		"https://elpais.com/",
	}); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	for i := 0; i < 4; i++ {
		workerThread := startWorkerThread(linkProcessor, errs)
		go workerThread()
	}
	for {
		select {
		case err <- errs:
			fmt.Println("Saw panic: %s. Starting new worker thread.", err.Error())
			workerThread := startWorkerThread(linkProcessor, errs)
			go workerThread()
		}
	}
}

func startWorkerThread(linkProcessor *linkprocessing.LinkProcessor, errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
			}
		}()
		for {
			var u string
			link, waitTime, err := linkProcessor.GetLink()
			switch {
			case err != nil:
				log.Println("Error getting link...")
				continue
			case waitTime != nil:
				log.Println("No link available. Sleeping...")
				time.Sleep(waitTime)
				continue
			case link != nil:
				u = link.URL
			default:
				log.Println("No error, but no wait time. Continuing...")
				continue
			}
			parsedHTMLPage, err := ingesthtml.ProcessURL(u, "google.com")
			if err != nil {
				log.Println("Got error ingesting html for url %s: %s. Continuing...")
			}
			languageCode := wordsmith.LookupLanguageCodeForLanguageLabel(deref.String(parsedHTMLPage.Language, ""))
			if languageCode == nil {
				log.Println("URL %s has unsupported language code: %s", u, deref.String(parsedHTMLPage.Language, ""))
			}
			if err := linkProcessor.AddURLs(parsedHTMLPage.Links); err != nil {
				log.Println("Error saving urls %+v for url %s: %s", parsedHTMLPage.Links, u, err.Error())
			}
			if strings.ToLower(deref.String(parsedHTMLPage.PageType, "")) != "article" {
				log.Println("URL %s is not an article. Continuing...", u)
				continue
			}
			textMetadata, err := textprocessing.ProcessText(parsedHtmlPage.BodyText, *languageCode)
			if err != nil {
				log.Println("Got error processing text for url %s: %s. Continuing...")
				continue
			}
			err = indexing.IndexDocument(indexing.IndexDocumentInput{
				ParsedHTMLPage:  *parsedHTMLPage,
				TextMetadata:    *textMetadata,
				LanguageCode:    *languageCode,
				DocumentVersion: documents.Version2,
				URL:             u,
			})
			if err != nil {
				log.Println("Got error indexing document for url %s: %s. Continuing...")
				continue
			}
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
