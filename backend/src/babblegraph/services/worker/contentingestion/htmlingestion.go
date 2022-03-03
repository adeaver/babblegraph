package contentingestion

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/services/worker/contentingestion/ingesthtml"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/opengraph"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func processWebsiteHTML1Link(c ctx.LogContext, link links2.Link) error {
	shouldMarkAsComplete := true
	defer func() {
		if shouldMarkAsComplete {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return links2.SetURLAsFetched(tx, link.URLIdentifier)
			}); err != nil {
				c.Errorf("Error marking link %s complete: %s", link.URLIdentifier, err.Error())
			}
		}
	}()
	c.Infof("Processing new link with URL %s", link.URL)
	u := link.URL
	p := urlparser.ParseURL(u)
	if p == nil {
		c.Infof("Received link that did not parse")
		return nil
	}
	var shouldExit bool
	var source *content.Source
	var sourceFilter *content.SourceFilter
	err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		shouldExit, err = content.IsParsedURLASeedURL(tx, *p)
		if err != nil {
			return err
		}
		source, err = content.LookupActiveSourceForSourceID(c, tx, *link.SourceID)
		if err != nil {
			return err
		}
		sourceFilter, err = content.LookupSourceFilterForSource(tx, *link.SourceID)
		return err
	})
	switch {
	case err != nil:
		c.Warnf("Error verifying if url %s is a seed url: %s", u, err.Error())
		return err
	case shouldExit == true:
		c.Infof("Link is a seed url")
		return nil
	case source == nil:
		c.Infof("No source found")
		return nil
	default:
		// no-op
	}
	parsedHTMLPage, err := ingesthtml.ProcessURL(ingesthtml.ProcessURLInput{
		URL:          link.URL,
		Source:       *source,
		SourceFilter: sourceFilter,
	})
	if err != nil {
		c.Errorf("Error parsing html for link %s: %s", link.URL, err.Error())
		return nil
	}
	if err := insertLinks(parsedHTMLPage.Links); err != nil {
		return err
	}
	var description *string
	if d, ok := parsedHTMLPage.Metadata[opengraph.DescriptionTag.Str()]; ok {
		description = ptr.String(d)
	}
	textMetadata, err := textprocessing.ProcessText(textprocessing.ProcessTextInput{
		BodyText:     parsedHTMLPage.BodyText,
		Description:  description,
		LanguageCode: source.LanguageCode,
	})
	if err != nil {
		c.Warnf("Got error processing text for url %s: %s. Continuing...", u, err.Error())
		return nil
	}
	var topicsForURL []contenttopics.ContentTopic
	var topicMappingIDs []content.TopicMappingID
	var topicIDs []content.TopicID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topicsForURL, topicMappingIDs, err = urltopicmapping.GetTopicsAndMappingIDsForURL(tx, u)
		if err != nil {
			return err
		}
		var sourceSeedTopicMappings []content.SourceSeedTopicMappingID
		for _, topicMappingID := range topicMappingIDs {
			sourceSeedTopicMapping, sourceTopicMapping, err := topicMappingID.GetOriginID()
			switch {
			case err != nil:
				return err
			case sourceTopicMapping != nil:
				c.Warnf("Found source topic mapping from ID %s, which is unsupported", topicMappingID)
			case sourceSeedTopicMapping != nil:
				sourceSeedTopicMappings = append(sourceSeedTopicMappings, *sourceSeedTopicMapping)
			default:
				return fmt.Errorf("unreachable")
			}
		}
		if len(sourceSeedTopicMappings) > 0 {
			topicIDs, err = content.LookupTopicsForSourceSeedMappingIDs(tx, sourceSeedTopicMappings)
			if err != nil {
				return err
			}
		}
		return err
	}); err != nil {
		c.Warnf("Error getting topics for url %s: %s. Continuing...", u, err.Error())
		return nil
	}
	c.Debugf("Indexing text for URL %s", u)
	err = indexing.IndexDocument(c, indexing.IndexDocumentInput{
		ParsedHTMLPage:         *parsedHTMLPage,
		TextMetadata:           *textMetadata,
		LanguageCode:           source.LanguageCode,
		DocumentVersion:        documents.CurrentDocumentVersion,
		URL:                    *p,
		SourceID:               source.ID.Ptr(),
		TopicsForURL:           topicsForURL,
		TopicIDs:               topicIDs,
		TopicMappingIDs:        topicMappingIDs,
		SeedJobIngestTimestamp: link.SeedJobIngestTimestamp,
	})
	if err != nil {
		c.Warnf("Got error indexing document for url %s: %s. Continuing...", u, err.Error())
	}
	return nil
}

func insertLinks(urls []string) error {
	var parsedURLs []urlparser.ParsedURL
	for _, u := range urls {
		parsedURL := urlparser.ParseURL(u)
		if parsedURL == nil {
			continue
		}
		parsedURLs = append(parsedURLs, *parsedURL)
	}
	return database.WithTx(func(tx *sqlx.Tx) error {
		var filteredURLs []links2.URLWithSourceMapping
		for _, u := range parsedURLs {
			sourceID, _, err := content.LookupSourceIDForDomain(tx, u.Domain)
			switch {
			case err != nil:
				return err
			case sourceID == nil:
				continue
			default:
				filteredURLs = append(filteredURLs, links2.URLWithSourceMapping{
					URL:      u,
					SourceID: *sourceID,
				})
			}
		}
		if len(filteredURLs) == 0 {
			return nil
		}
		return links2.InsertLinksWithSourceID(tx, filteredURLs)
	})
}
