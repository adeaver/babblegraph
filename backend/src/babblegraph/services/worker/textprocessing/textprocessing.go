package textprocessing

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/textprocessing/spanishprocessing"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
)

type TextMetadata struct {
	ReadabilityScore decimal.Number
	WordStats        documents.WordStatsVersion1
}

func ProcessText(text string, language wordsmith.LanguageCode) (*TextMetadata, error) {
	normalizedText := normalizeText(text)
	wordStats, err := getWordStatsForText(language, normalizedText)
	if err != nil {
		return nil, err
	}
	switch language {
	case wordsmith.LanguageCodeSpanish:
		readabilityScore, err := spanishprocessing.CalculateReadabilityForSpanish(normalizedText)
		if err != nil {
			return nil, err
		}
		return &TextMetadata{
			ReadabilityScore: *readabilityScore,
			WordStats:        *wordStats,
		}, nil
	default:
		panic("unrecognized language")
	}
}
