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
	for u, topics := range domains.GetSeedURLs() {
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
			if err := links2.UpsertLinkWithEmptyFetchStatus(tx, parsedURLs); err != nil {
				return err
			}
			if len(topics) == 0 {
				return nil
			}
			for _, u := range parsedURLs {
				if err := contenttopics.ApplyContentTopicsToURL(tx, u.URL, topics); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			log.Println(fmt.Sprintf("Error inserting refetched urls for %s: %s", u, err.Error()))
			continue
		}
	}
	return nil
}
