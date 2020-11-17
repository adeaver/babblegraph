package documents

import (
	"babblegraph/wordsmith"
)

type Version int64

const (
	Version1 Version = 1

	CurrentDocumentVersion Version = Version1
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
	DocumentType     *Type                  `json:"document_type"`
	Metadata         *Metadata              `json:"metadata"`
}
