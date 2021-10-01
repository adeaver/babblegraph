package process

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/opengraph"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func StartIngestWorkerThread(workerNumber int, linkProcessor *linkprocessing.LinkProcessor, errs chan error) func() {
	return func() {
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("worker-thread", fmt.Sprintf("init#%d", workerNumber))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Worker Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		fmt.Println("Starting Ingest Worker Process")
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
				ParsedHTMLPage:         *parsedHTMLPage,
				TextMetadata:           *textMetadata,
				LanguageCode:           languageCode,
				DocumentVersion:        documents.CurrentDocumentVersion,
				URL:                    urlparser.MustParseURL(u),
				TopicsForURL:           topicsForURL,
				SeedJobIngestTimestamp: link.SeedJobIngestTimestamp,
			})
			if err != nil {
				log.Println(fmt.Sprintf("Got error indexing document for url %s: %s. Continuing...", u, err.Error()))
				localHub.CaptureException(err)
				continue
			}
		}
	}
}
