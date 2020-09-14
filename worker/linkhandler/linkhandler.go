package linkhandler

import (
	"babblegraph/worker/links"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func FilterLinksAndInsert(tx *sqlx.Tx, urls []string) ([]links.Link, error) {
	filteredLinks := getFilteredLinksForURLs(urls)
	if err := links.InsertLinks(tx, filteredLinks); err != nil {
		return nil, err
	}
	return filteredLinks, nil
}

var forbiddenDomains = map[string]bool{
	"google.com":    true,
	"instagram.com": true,
	"pinterest.com": true,
	"facebook.com":  true,
	"twitter.com":   true,
	"youtube.com":   true,
}

func getFilteredLinksForURLs(urls []string) []links.Link {
	// I anticipate that this function will have to go through
	// a lot of garbage. Therefore, I will absorb all the errors
	var out []links.Link
	for _, u := range urls {
		switch {
		case strings.Contains(u, "mailto"), // mailto string
			strings.Contains(u, "//") && !strings.HasPrefix(u, "http"): // wrong protocol
			continue
		}
		link, err := links.GetLinkForURL(u)
		if err != nil {
			log.Println(fmt.Sprintf("Error on making link for url %s with message: %s", u, err.Error()))
			continue
		}
		if _, isForbiddenDomain := forbiddenDomains[link.Domain.Str()]; isForbiddenDomain {
			continue
		}
		log.Println(fmt.Sprintf("Adding URL: %s", u))
		out = append(out, *link)
	}
	return out
}
