package links2

import (
	"babblegraph/model/content"
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

// TODO: this is deprecated
func InsertLinks(tx *sqlx.Tx, urls []urlparser.ParsedURL) error {
	queryBuilder, err := database.NewBulkInsertQueryBuilder("links2", "url_identifier", "domain", "url", "source_id")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("DO NOTHING")
	for _, u := range urls {
		sourceID, err := content.LookupSourceIDForParsedURL(tx, u)
		switch {
		case err != nil:
			return err
		case sourceID == nil:
			return fmt.Errorf("No source ID found for url %s", u.URL)
		}
		if err := queryBuilder.AddValues(u.URLIdentifier, u.Domain, u.URL, *sourceID); err != nil {
			log.Println(fmt.Sprintf("Error inserting url with identifier %s: %s", u.URLIdentifier, err.Error()))
		}
	}
	return queryBuilder.Execute(tx)
}

type URLWithSourceMapping struct {
	URL      urlparser.ParsedURL
	SourceID content.SourceID
}

func InsertLinksWithSourceID(tx *sqlx.Tx, urls []URLWithSourceMapping) error {
	queryBuilder, err := database.NewBulkInsertQueryBuilder("links2", "url_identifier", "domain", "url", "source_id")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("DO NOTHING")
	for _, u := range urls {
		if err := queryBuilder.AddValues(input.URL.URLIdentifier, input.URL.Domain, input.URL.URL, input.SourceID); err != nil {
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
	if err := tx.Select(&matches, "SELECT * FROM links2 WHERE last_fetch_version IS DISTINCT FROM $1 AND domain=$2 ORDER BY seed_job_ingest_timestamp DESC NULLS LAST, seq_num ASC LIMIT $3", CurrentFetchVersion, domain, chunkSize); err != nil {
		return nil, err
	}
	var out []Link
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}

func LookupBulkUnfetchedLinksForSourceID(tx *sqlx.Tx, sourceID content.SourceID, chunkSize int) ([]Link, error) {
	var matches []dbLink
	if err := tx.Select(&matches, "SELECT * FROM links2 WHERE last_fetch_version IS DISTINCT FROM $1 AND source_id=$2 ORDER BY seed_job_ingest_timestamp DESC NULLS LAST, seq_num ASC LIMIT $3", CurrentFetchVersion, sourceID, chunkSize); err != nil {
		return nil, err
	}
	var out []Link
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}

func UpsertLinkWithEmptyFetchStatus(tx *sqlx.Tx, urls []urlparser.ParsedURL, includeTimestamp bool) error {
	queryBuilder, err := database.NewBulkInsertQueryBuilder("links2", "url_identifier", "domain", "url", "source_id", "seed_job_ingest_timestamp")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("(url_identifier) DO UPDATE SET last_fetch_version = NULL")
	var firstSeedFetchTimestamp *int64
	if includeTimestamp {
		firstSeedFetchTimestamp = ptr.Int64(time.Now().Unix())
	}
	for _, u := range urls {
		sourceID, err := content.LookupSourceIDForParsedURL(tx, u)
		switch {
		case err != nil:
			return err
		case sourceID == nil:
			return fmt.Errorf("No source available for URL %s", u.URL)
		}
		if err := queryBuilder.AddValues(u.URLIdentifier, u.Domain, u.URL, *sourceID, firstSeedFetchTimestamp); err != nil {
			log.Println(fmt.Sprintf("Error inserting url with identifier %s: %s", u.URLIdentifier, err.Error()))
		}
	}
	return queryBuilder.Execute(tx)
}

// This is useful for reindexing
func GetLinksCursor(tx *sqlx.Tx, fn func(link Link) (bool, error)) error {
	rows, err := tx.Queryx("SELECT * FROM links2 WHERE source_id IS NOT NULL ORDER BY seq_num ASC")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var l dbLink
		if err := rows.StructScan(&l); err != nil {
			return err
		}
		shouldEndIteration, err := fn(l.ToNonDB())
		switch {
		case err != nil && shouldEndIteration:
			return err
		case err != nil && !shouldEndIteration:
			log.Println(fmt.Sprintf("Got error: %s, continuing...", err.Error()))
		}
	}
	return nil
}

// TODO(content-migration): get rid of this
func UpdateLinkSource(tx *sqlx.Tx, u urlparser.ParsedURL, sourceID content.SourceID) error {
	if _, err := tx.Exec("UPDATE links2 SET source_id = $1 WHERE url_identifier = $2", sourceID, u.URLIdentifier); err != nil {
		return err
	}
	return nil
}
