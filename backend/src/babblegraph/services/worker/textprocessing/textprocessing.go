package textprocessing

import (
	"babblegraph/services/worker/textprocessing/spanishprocessing"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
)

type TextMetadata struct {
	LemmatizedText   string
	ReadabilityScore decimal.Number
}

func ProcessText(text string, language wordsmith.LanguageCode) (*TextMetadata, error) {
	normalizedText := normalizeText(text)
	switch language {
	case wordsmith.LanguageCodeSpanish:
		readabilityScore, err := spanishprocessing.CalculateReadabilityForSpanish(normalizedText)
		if err != nil {
			return nil, err
		}
		lemmatizedBody, err := lemmatizeBody(wordsmith.LanguageCodeSpanish, normalizedText)
		if err != nil {
			return nil, err
		}
		return &TextMetadata{
			LemmatizedText:   *lemmatizedBody,
			ReadabilityScore: *readabilityScore,
		}, nil
	default:
		panic("unrecognized language")
	}
}
