package documents

import (
	"babblegraph/wordsmith"
)

type DocumentID string

type Document struct {
	ID               DocumentID             `json:"id"`
	URL              string                 `json:"url"`
	ReadabilityScore int64                  `json:"readability_score"`
	LanguageCode     wordsmith.LanguageCode `json:"language_code"`
	LemmatizedBody   string                 `json:"lemmatized_body"`
}
