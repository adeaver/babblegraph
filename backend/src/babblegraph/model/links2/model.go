package links2

import "time"

type URLIdentifier string

type Link struct {
	URLIdentifier    URLIdentifier
	Domain           string
	URL              string
	LastFetchVersion *FetchVersion
	FetchedOn        *time.Time

	// This is the UNIX timestamp of the first
	// time the seed job picked up the link.
	// IMPORTANT: The crawler should not populate this field.
	FirstSeedFetchTimestamp *int64
}

type FetchVersion int64

const (
	FetchVersion1 FetchVersion = 1

	// Version 2 Updates.
	// Adds publication time, domain, and description to documents
	// Removes lemmatized body
	FetchVersion2 FetchVersion = 2

	// Version 3 Updates:
	// Fix bug with urlparser
	FetchVersion3 FetchVersion = 3

	// Version 4 Updates:
	// Fix bug with encoding of html pages
	FetchVersion4 FetchVersion = 4

	CurrentFetchVersion FetchVersion = FetchVersion4
)

type dbLink struct {
	URLIdentifier           URLIdentifier `db:"url_identifier"`
	Domain                  string        `db:"domain"`
	URL                     string        `db:"url"`
	LastFetchVersion        *FetchVersion `db:"last_fetch_version"`
	FetchedOn               *time.Time    `db:"fetched_on"`
	SeqNum                  int           `db:"seq_num"`
	FirstSeedFetchTimestamp *int64        `db:"first_seed_fetch_timestamp"`
}

func (d dbLink) ToNonDB() Link {
	return Link{
		URLIdentifier:    d.URLIdentifier,
		Domain:           d.Domain,
		URL:              d.URL,
		LastFetchVersion: d.LastFetchVersion,
		FetchedOn:        d.FetchedOn,
	}
}
