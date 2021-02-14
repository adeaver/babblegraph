package documents

import (
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

	CurrentDocumentVersion Version = Version4
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

type Document struct {
	ID                                 DocumentID                   `json:"id"`
	Version                            Version                      `json:"version"`
	URL                                string                       `json:"url"`
	ReadabilityScore                   int64                        `json:"readability_score"`
	LanguageCode                       wordsmith.LanguageCode       `json:"language_code"`
	DocumentType                       Type                         `json:"document_type"`
	Metadata                           Metadata                     `json:"metadata"`
	Domain                             string                       `json:"domain"`
	Topics                             []contenttopics.ContentTopic `json:"content_topics"`
	LemmatizedDescription              *string                      `json:"lemmatized_description,omitempty"`
	LemmatizedDescriptionIndexMappings []int                        `json:"lemmatized_description_index_mappings,omitempty"`

	LemmatizedBodyDEPRECATED *string `json:"lemmatized_body,omitempty"`
}
