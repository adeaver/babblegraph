package links2

import (
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

type domainQuery struct {
	Domain string `db:"domain"`
}

func GetDomainsWithUnfetchedLinks(tx *sqlx.Tx) ([]string, error) {
	var domainsWithUnfetchedLinks []domainQuery
	if err := tx.Select(&domainsWithUnfetchedLinks, "SELECT DISTINCT(domain) FROM links2 WHERE last_fetch_version IS DISTINCT FROM $1", FetchVersion1); err != nil {
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
	for _, u := range urls {
		if _, err := tx.Exec(insertLinkQuery, u.URLIdentifier, u.Domain, u.URL); err != nil {
			return err
		}
	}
	return nil
}

func SetURLAsFetched(tx *sqlx.Tx, urlIdentifier URLIdentifier) error {
	if _, err := tx.Exec("UPDATE links2 SET last_fetch_version=$1 WHERE url_identifier=$2", FetchVersion1, string(urlIdentifier)); err != nil {
		return err
	}
	return nil
}

func LookupUnfetchedLinkForDomain(tx *sqlx.Tx, domain string) (*Link, error) {
	var matches []dbLink
	if err := tx.Select(&matches, "SELECT * FROM links2 WHERE last_fetch_version IS DISTINCT FROM $1 AND domain=$2 ORDER BY seq_num ASC LIMIT 1", FetchVersion1, domain); err != nil {
		return nil, err
	}
	if len(matches) < 1 {
		return nil, nil
	}
	l := matches[0].ToNonDB()
	return &l, nil
}

const upsertLinkQuery = "INSERT INTO links2 (url_identifier, domain, url, last_fetch_version) VALUES ($1, $2, $3, NULL) ON CONFLICT (url_identifier) DO UPDATE SET last_fetch_version = NULL"

func UpsertLinkWithEmptyFetchStatus(tx *sqlx.Tx, urls []urlparser.ParsedURL) error {
	for _, u := range urls {
		if _, err := tx.Exec(upsertLinkQuery, u.URLIdentifier, u.Domain, u.URL); err != nil {
			return err
		}
	}
	return nil
}
