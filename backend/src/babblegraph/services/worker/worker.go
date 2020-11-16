package main

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/deref"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

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
			log.Println(fmt.Sprintf("Processing URL %s", u))
			parsedHTMLPage, err := ingesthtml.ProcessURL(u, domain)
			if err != nil {
				log.Println(fmt.Sprintf("Got error ingesting html for url %s: %s. Continuing...", u, err.Error()))
				continue
			}
			languageCode := wordsmith.LookupLanguageCodeForLanguageLabel(deref.String(parsedHTMLPage.Language, ""))
			if languageCode == nil {
				log.Println(fmt.Sprintf("URL %s has unsupported language code: %s", u, deref.String(parsedHTMLPage.Language, "")))
				// This is only allowed because we're restricting domains
				// If domains are ever non-restricted, this needs to be removed
				// and made more robust
				languageCode = wordsmith.LanguageCodeSpanish.Ptr()
			}
			log.Println(fmt.Sprintf("Got language code %s for label %s on URL %s. Processing...", languageCode.Str(), *parsedHTMLPage.Language, u))
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
