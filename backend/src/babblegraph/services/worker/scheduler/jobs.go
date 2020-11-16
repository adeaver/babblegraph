package scheduler

import (
	"babblegraph/model/links2"
	"babblegraph/services/worker/domains"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func RefetchSeedDomainsForNewContent() error {
	log.Println(fmt.Sprintf("Starting refetch of seed domains..."))
	for _, u := range domains.GetSeedURLs() {
		log.Println(fmt.Sprintf("Refetch processing seed domain %s", u))
		parsedHTMLPage, err := ingesthtml.ProcessURL(u, u)
		if err != nil {
			log.Println(fmt.Sprintf("Got error ingesting html for url %s: %s. Continuing...", u, err.Error()))
			continue
		}
		var parsedURLs []urlparser.ParsedURL
		for _, l := range parsedHTMLPage.Links {
			if p := urlparser.ParseURL(l); p != nil && domains.IsURLAllowed(*p) {
				parsedURLs = append(parsedURLs, *p)
			}
		}
		log.Println(fmt.Sprintf("Inserting refetched urls for %s", u))
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			return links2.UpsertLinkWithEmptyFetchStatus(tx, parsedURLs)
		}); err != nil {
			log.Println(fmt.Sprintf("Error inserting refetched urls for %s: %s", u, err.Error()))
			continue
		}
	}
	return nil
}
