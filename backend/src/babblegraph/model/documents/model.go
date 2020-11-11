package documents

import (
	"babblegraph/wordsmith"
)

type Version int64

const (
	Version1 Version = 1
	Version2 Version = 2

	CurrentDocumentVersion Version = Version2
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
	Title string `json:"title"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

type DocumentID string

type Document struct {
	ID               DocumentID             `json:"id"`
	Version          Version                `json:"version"`
	URL              string                 `json:"url"`
	ReadabilityScore int64                  `json:"readability_score"`
	LanguageCode     wordsmith.LanguageCode `json:"language_code"`
	LemmatizedBody   string                 `json:"lemmatized_body"`

	// Added in version 2
	DocumentType *Type     `json:"document_type"`
	Metadata     *Metadata `json:"metadata"`
}
