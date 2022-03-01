package scheduler

import (
	"babblegraph/model/content"
	"babblegraph/model/domains"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func refetchSeedDomainsForNewContent(c ctx.LogContext) error {
	c.Infof("Starting refetch of seed domains...")
	urlsByDomain := make(map[string][]domains.SeedURL)
	for _, domain := range domains.GetDomains() {
		urlsByDomain[domain] = []domains.SeedURL{}
	}
	for u, topics := range domains.GetSeedURLs() {
		if p := urlparser.ParseURL(u); p != nil && domains.IsURLAllowed(*p) {
			domain := p.Domain
			seedURLs, ok := urlsByDomain[domain]
			if !ok {
				c.Warnf("Error processing seed url %s: its domain of %s not found in allowable domains", u, domain)
				continue
			}
			urlsByDomain[domain] = append(seedURLs, domains.SeedURL{
				URL:    u,
				Topics: topics,
			})
		}
	}
	for len(urlsByDomain) != 0 {
		// In order to rate limit ourselves (more or less), we go through each domain one by one
		// And remove the first seed url, process it, and delete from the list.
		// Once the list for that domain is empty, we delete that list
		for domain := range urlsByDomain {
			seedURLs := urlsByDomain[domain]
			toProcess := seedURLs[0]
			if err := processSeedURL(c, toProcess); err != nil {
				c.Warnf("Error processing seed url %s: %s", toProcess.URL, err.Error())
			}
			c.Infof("Refetch finished processing seed url %s", toProcess.URL)
			if len(seedURLs) == 1 {
				delete(urlsByDomain, domain)
			} else {
				urlsByDomain[domain] = append([]domains.SeedURL{}, seedURLs[1:]...)
			}
		}
	}
	return nil
}

func processSeedURL(c ctx.LogContext, seedURL domains.SeedURL) error {
	c.Infof("Refetch is processing seed url %s", seedURL.URL)
	parsedSeedURL := urlparser.ParseURL(seedURL.URL)
	if parsedSeedURL == nil {
		return fmt.Errorf("something went wrong parsing url, got null parsed url for seed url %s", seedURL.URL)
	}
	parsedHTMLPage, err := ingesthtml.ProcessURL(seedURL.URL, parsedSeedURL.Domain)
	if err != nil {
		return err
	}
	parsedURLs := make(map[string]urlparser.ParsedURL)
	for _, l := range parsedHTMLPage.Links {
		if p := urlparser.ParseURL(l); p != nil && domains.IsURLAllowed(*p) && !domains.IsSeedURL(*p) {
			parsedURLs[p.URLIdentifier] = *p
		}
	}
	c.Infof("Inserting refetched urls for %s", seedURL.URL)
	return database.WithTx(func(tx *sqlx.Tx) error {
		var toInsert []urlparser.ParsedURL
		for _, p := range parsedURLs {
			toInsert = append(toInsert, p)
		}
		if err := links2.UpsertLinkWithEmptyFetchStatus(tx, toInsert, true); err != nil {
			return err
		}
		var mappings []urltopicmapping.TopicMappingUnion
		for _, t := range seedURL.Topics {
			topicID, err := content.GetTopicIDByContentTopic(tx, t)
			if err != nil {
				return err
			}
			topicMappingID, err := content.LookupTopicMappingIDForURL(c, tx, *parsedSeedURL, *topicID)
			switch {
			case err != nil:
				return err
			case topicMappingID != nil:
				mappings = append(mappings, urltopicmapping.TopicMappingUnion{
					Topic:          t,
					TopicMappingID: *topicMappingID,
				})
			case topicMappingID == nil:
				c.Warnf("Topic %s with ID %s did not map to anything for seed URL %s", t, *topicID, *parsedSeedURL)
			}
		}
		if len(mappings) > 0 {
			return nil
		}
		for _, u := range parsedURLs {
			if err := urltopicmapping.ApplyContentTopicsToURL(tx, u, mappings); err != nil {
				return err
			}
		}
		return nil
	})
}
