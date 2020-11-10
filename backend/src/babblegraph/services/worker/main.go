package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"strings"
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	for i := 0; i < 4; i++ {
		workerThread := startWorkerThread(errs)
		go workerThread()
	}
	for {
		select {
		case err <- errs:
			fmt.Println("Saw panic: %s. Starting new worker thread.", err.Error())
			workerThread := startWorkerThread(errs)
			go workerThread()
		}
	}
}

func startWorkerThread(errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
			}
		}()
		for {
			// Get next URL
			u := "www.google.com"
			parsedHTMLPage, err := ingesthtml.ProcessURL(u, "google.com")
			if err != nil {
				log.Println("Got error ingesting html for url %s: %s. Continuing...")
			}
			languageCode := wordsmith.LookupLanguageCodeForLanguageLabel(deref.String(parsedHTMLPage.Language, ""))
			if languageCode == nil {
				log.Println("URL %s has unsupported language code: %s", u, deref.String(parsedHTMLPage.Language, ""))
			}
			// TODO: persist links and send domain to link scheduler
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
