package process

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/opengraph"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func StartIngestWorkerThread(linkProcessor *linkprocessing.LinkProcessor) func(c async.Context) {
	return func(c async.Context) {
		c.Infof("Starting Ingest Worker Process")
		for {
			var u, domain string
			link, waitTime, err := linkProcessor.GetLink()
			switch {
			case err != nil:
				c.Infof("Error getting link... %s", err.Error())
				continue
			case waitTime != nil:
				c.Infof("No link available. Sleeping...")
				time.Sleep(*waitTime)
				continue
			case link != nil:
				u = link.URL
				domain = link.Domain
			default:
				c.Infof("No error, but no wait time. Continuing...")
				continue
			}
			if p := urlparser.ParseURL(u); p != nil && domains.IsSeedURL(*p) {
				c.Infof("Received url %s, which is a seed url. Skipping...", u)
				continue
			}
			c.Infof("Processing URL %s with identifier %s", u, link.URLIdentifier)
			parsedHTMLPage, err := ingesthtml.ProcessURL(u, domain)
			if err != nil {
				c.Infof("Got error ingesting html for url %s: %s. Continuing...", u, err.Error())
				continue
			}
			domainMetadata, err := domains.GetDomainMetadata(domain)
			if err != nil {
				c.Infof("Got error getting metadata for domain %s on url %s: %s. Continuing...", domain, u, err.Error())
				continue
			}
			languageCode := domainMetadata.LanguageCode
			if err := linkProcessor.AddURLs(parsedHTMLPage.Links, domainMetadata.Topics); err != nil {
				c.Infof("Error saving urls %+v for url %s: %s", parsedHTMLPage.Links, u, err.Error())
				continue
			}
			if strings.ToLower(deref.String(parsedHTMLPage.PageType, "")) != "article" {
				c.Infof("URL %s is not an article. Continuing...", u)
				continue
			}
			c.Infof("Processing text for url %s", u)
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
				c.Infof("Got error processing text for url %s: %s. Continuing...", u, err.Error())
				continue
			}
			var topicsForURL []contenttopics.ContentTopic
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				topicsForURL, err = contenttopics.GetTopicsForURL(tx, u)
				return err
			}); err != nil {
				c.Infof("Error getting topics for url %s: %s. Continuing...", u, err.Error())
				continue
			}
			c.Infof("Indexing text for URL %s", u)
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
				c.Errorf("Got error indexing document for url %s: %s. Continuing...", u, err.Error())
				continue
			}
		}
	}
}
