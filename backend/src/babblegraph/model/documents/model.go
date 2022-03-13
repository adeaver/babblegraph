package documents

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/wordsmith"
	"time"
)

type Version int64

const (
	Version1 Version = 1
	Version2 Version = 2

	// Version 3 adds content topics
	Version3 Version = 3

	// Version 4 adds lemmatized description
	Version4 Version = 4

	// Version 5 fixes an issue with URL parser
	Version5 Version = 5

	// Version 6 adds paywall tracking
	Version6 Version = 6

	// Version 7 adds content topics length filter
	Version7 Version = 7

	// Version 8 adds source ID and topic mapping
	Version8 Version = 8

	CurrentDocumentVersion Version = Version8
)

func (v Version) Ptr() *Version {
	return &v
}

type Type string

const (
	TypeArticle Type = "article"
)

func (t Type) Ptr() *Type {
	return &t
}

func (t Type) Str() string {
	return string(t)
}

type Metadata struct {
	Title              *string    `json:"title,omitempty"`
	Image              *string    `json:"image,omitempty"`
	URL                *string    `json:"url,omitempty"`
	Description        *string    `json:"description,omitempty"`
	PublicationTimeUTC *time.Time `json:"publication_time_utc,omitempty"`
}

type DocumentID string

func (d DocumentID) Str() string {
	return string(d)
}

type Document struct {
	ID               DocumentID             `json:"id"`
	Version          Version                `json:"version"`
	URL              string                 `json:"url"`
	SourceID         *content.SourceID      `json:"source_id,omitempty"`
	ReadabilityScore int64                  `json:"readability_score"`
	LanguageCode     wordsmith.LanguageCode `json:"language_code"`
	DocumentType     Type                   `json:"document_type"`
	Metadata         Metadata               `json:"metadata"`
	// We need topic IDs for relevance queries, but topic mapping IDs to filter
	TopicIDs                           []content.TopicID        `json:"topic_ids"`
	TopicMappingIDs                    []content.TopicMappingID `json:"topic_mapping_ids"`
	TopicsLength                       *int64                   `json:"topics_length,omitempty"`
	SeedJobIngestTimestamp             *int64                   `json:"seed_job_ingest_timestamp,omitempty"`
	LemmatizedDescription              *string                  `json:"lemmatized_description,omitempty"`
	HasPaywall                         *bool                    `json:"has_paywall"`
	LemmatizedDescriptionIndexMappings []int                    `json:"lemmatized_description_index_mappings,omitempty"`

	// These will all be deprecated
	Domain                   string                       `json:"domain"`
	Topics                   []contenttopics.ContentTopic `json:"content_topics"`
	LemmatizedBodyDEPRECATED *string                      `json:"lemmatized_body,omitempty"`
}
