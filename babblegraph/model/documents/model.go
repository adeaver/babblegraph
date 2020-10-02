package documents

import (
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
)

type DocumentID string

type Document struct {
	ID               DocumentID             `json:"id"`
	URL              string                 `json:"url"`
	ReadabilityScore decimal.Number         `json:"readability_score"`
	LanguageCode     wordsmith.LanguageCode `json:"language_code"`
	LemmatizedBody   string                 `json:"lemmatized_body"`
}
