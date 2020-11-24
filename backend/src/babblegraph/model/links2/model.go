package links2

import "time"

type URLIdentifier string

type Link struct {
	URLIdentifier    URLIdentifier
	Domain           string
	URL              string
	LastFetchVersion *FetchVersion
	FetchedOn        *time.Time
}

type FetchVersion int64

const (
	FetchVersion1 FetchVersion = 1

	// Version 2 Updates.
	// Adds publication time, domain, and description to documents
	// Removes lemmatized body
	FetchVersion2 FetchVersion = 2

	CurrentFetchVersion FetchVersion = FetchVersion2
)

type dbLink struct {
	URLIdentifier    URLIdentifier `db:"url_identifier"`
	Domain           string        `db:"domain"`
	URL              string        `db:"url"`
	LastFetchVersion *FetchVersion `db:"last_fetch_version"`
	FetchedOn        *time.Time    `db:"fetched_on"`
	SeqNum           int           `db:"seq_num"`
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
