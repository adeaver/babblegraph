package links2

import (
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type domainQuery struct {
	Domain string `db:"domain"`
}

func GetDomainsWithUnfetchedLinks(tx *sqlx.Tx) ([]string, error) {
	var domainsWithUnfetchedLinks []domainQuery
	if err := tx.Select(&domainsWithUnfetchedLinks, "SELECT DISTINCT(domain) FROM links2 WHERE last_fetch_version IS DISTINCT FROM $1", CurrentFetchVersion); err != nil {
		return nil, err
	}
	var out []string
	for _, domain := range domainsWithUnfetchedLinks {
		out = append(out, domain.Domain)
	}
	return out, nil
}

const insertLinkQuery = "INSERT INTO links2 (url_identifier, domain, url) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"

func InsertLinks(tx *sqlx.Tx, urls []urlparser.ParsedURL) error {
	queryBuilder, err := database.NewBulkInsertQueryBuilder("links2", "url_identifier", "domain", "url")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("DO NOTHING")
	for _, u := range urls {
		if err := queryBuilder.AddValues(u.URLIdentifier, u.Domain, u.URL); err != nil {
			log.Println(fmt.Sprintf("Error inserting url with identifier %s: %s", u.URLIdentifier, err.Error()))
		}
	}
	return queryBuilder.Execute(tx)
}

func SetURLAsFetched(tx *sqlx.Tx, urlIdentifier URLIdentifier) error {
	if _, err := tx.Exec("UPDATE links2 SET last_fetch_version=$1 WHERE url_identifier=$2", CurrentFetchVersion, string(urlIdentifier)); err != nil {
		return err
	}
	return nil
}

func LookupUnfetchedLinkForDomain(tx *sqlx.Tx, domain string) (*Link, error) {
	out, err := LookupBulkUnfetchedLinksForDomain(tx, domain, 1)
	if err != nil {
		return nil, err
	}
	if len(out) != 1 {
		return nil, nil
	}
	return &out[0], nil
}

func LookupBulkUnfetchedLinksForDomain(tx *sqlx.Tx, domain string, chunkSize int) ([]Link, error) {
	var matches []dbLink
	if err := tx.Select(&matches, "SELECT * FROM links2 WHERE last_fetch_version IS DISTINCT FROM $1 AND domain=$2 ORDER BY seq_num ASC LIMIT $3", CurrentFetchVersion, domain, chunkSize); err != nil {
		return nil, err
	}
	var out []Link
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}

func UpsertLinkWithEmptyFetchStatus(tx *sqlx.Tx, urls []urlparser.ParsedURL, includeTimestamp bool) error {
	queryBuilder, err := database.NewBulkInsertQueryBuilder("links2", "url_identifier", "domain", "url", "first_seed_fetch_timestamp")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("(url_identifier) DO UPDATE SET last_fetch_version = NULL")
	var firstSeedFetchTimestamp *int64
	if includeTimestamp {
		firstSeedFetchTimestamp = ptr.Int64(time.Now().Unix())
	}
	for _, u := range urls {
		if err := queryBuilder.AddValues(u.URLIdentifier, u.Domain, u.URL, firstSeedFetchTimestamp); err != nil {
			log.Println(fmt.Sprintf("Error inserting url with identifier %s: %s", u.URLIdentifier, err.Error()))
		}
	}
	return queryBuilder.Execute(tx)
}
