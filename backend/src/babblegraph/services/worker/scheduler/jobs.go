package scheduler

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/domains"
	"babblegraph/model/links2"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func RefetchSeedDomainsForNewContent() error {
	log.Println(fmt.Sprintf("Starting refetch of seed domains..."))
	urlsByDomain := make(map[string][]domains.SeedURL)
	for _, domain := range domains.GetDomains() {
		urlsByDomain[domain] = []domains.SeedURL{}
	}
	for u, topics := range domains.GetSeedURLs() {
		if p := urlparser.ParseURL(u); p != nil && domains.IsURLAllowed(*p) {
			domain := p.Domain
			seedURLs, ok := urlsByDomain[domain]
			if !ok {
				log.Println(fmt.Sprintf("Error processing seed url %s: its domain of %s not found in allowable domains", u, domain))
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
			if err := processSeedURL(toProcess); err != nil {
				log.Println(fmt.Sprintf("Error processing seed url %s: %s", toProcess.URL, err.Error()))
			}
			log.Println(fmt.Sprintf("Refetch finished processing seed url %s", toProcess.URL))
			if len(seedURLs) == 1 {
				delete(urlsByDomain, domain)
			} else {
				urlsByDomain[domain] = append([]domains.SeedURL{}, seedURLs[1:]...)
			}
		}
	}
	return nil
}

func processSeedURL(seedURL domains.SeedURL) error {
	log.Println(fmt.Sprintf("Refetch is processing seed url %s", seedURL.URL))
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
	log.Println(fmt.Sprintf("Inserting refetched urls for %s", seedURL.URL))
	return database.WithTx(func(tx *sqlx.Tx) error {
		var toInsert []urlparser.ParsedURL
		for _, p := range parsedURLs {
			toInsert = append(toInsert, p)
		}
		if err := links2.UpsertLinkWithEmptyFetchStatus(tx, toInsert); err != nil {
			return err
		}
		if len(seedURL.Topics) == 0 {
			return nil
		}
		for _, u := range parsedURLs {
			if err := contenttopics.ApplyContentTopicsToURL(tx, u.URL, seedURL.Topics); err != nil {
				return err
			}
		}
		return nil
	})
}
