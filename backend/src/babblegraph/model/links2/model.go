package links2

import (
	"babblegraph/model/content"
	"babblegraph/util/urlparser"
	"time"

	"github.com/jmoiron/sqlx"
)

type URLIdentifier string

type Link struct {
	URLIdentifier    URLIdentifier
	Domain           string
	SourceID         *content.SourceID
	URL              string
	LastFetchVersion *FetchVersion
	FetchedOn        *time.Time

	// This is the UNIX timestamp of the first
	// time the seed job picked up the link.
	// IMPORTANT: The crawler should not populate this field.
	// This is used as an approximation for publication date time.
	SeedJobIngestTimestamp *int64
}

func (l *Link) GetSourceID(tx *sqlx.Tx) (*content.SourceID, error) {
	if l.SourceID != nil {
		return l.SourceID, nil
	}
	return content.GetSourceIDForParsedURL(tx, urlparser.MustParseURL(l.URL))
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

	// Version 5 Updates:
	// Add Paywall Detection
	FetchVersion5 FetchVersion = 5

	// Version 6 Updates:
	// Fixes some bugs with Paywall Detection
	FetchVersion6 FetchVersion = 6

	// Version 7 Updates:
	// Added new elasticsearch instance, slowly reindexing everything
	FetchVersion7 FetchVersion = 7

	// Version 8 Updates:
	// Added new fields onto model/documents
	FetchVersion8 FetchVersion = 8

	CurrentFetchVersion FetchVersion = FetchVersion8
)

type dbLink struct {
	URLIdentifier          URLIdentifier     `db:"url_identifier"`
	Domain                 string            `db:"domain"`
	SourceID               *content.SourceID `db:"source_id"`
	URL                    string            `db:"url"`
	LastFetchVersion       *FetchVersion     `db:"last_fetch_version"`
	FetchedOn              *time.Time        `db:"fetched_on"`
	SeqNum                 int64             `db:"seq_num"`
	SeedJobIngestTimestamp *int64            `db:"seed_job_ingest_timestamp"`
}

func (d dbLink) ToNonDB() Link {
	return Link{
		URLIdentifier:          d.URLIdentifier,
		Domain:                 d.Domain,
		SourceID:               d.SourceID,
		URL:                    d.URL,
		LastFetchVersion:       d.LastFetchVersion,
		FetchedOn:              d.FetchedOn,
		SeedJobIngestTimestamp: d.SeedJobIngestTimestamp,
	}
}
