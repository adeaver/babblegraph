package textprocessing

import (
	"babblegraph/services/worker/textprocessing/spanishprocessing"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
	"strings"
)

type TextMetadata struct {
	ReadabilityScore      decimal.Number
	LemmatizedDescription *LemmatizedDescription
}

type LemmatizedDescription struct {
	LemmatizedText string
	IndexMappings  []int
}

type ProcessTextInput struct {
	BodyText     string
	Description  *string
	LanguageCode wordsmith.LanguageCode
}

func ProcessText(input ProcessTextInput) (*TextMetadata, error) {
	var normalizedDescription *string
	if input.Description != nil {
		normalizedDescription = ptr.String(text.Normalize(*input.Description))
	}
	normalizedBodyText := text.Normalize(input.BodyText)
	switch input.LanguageCode {
	case wordsmith.LanguageCodeSpanish:
		readabilityScore, err := spanishprocessing.CalculateReadabilityForSpanish(normalizedBodyText)
		if err != nil {
			return nil, err
		}
		var lemmatizedDescription *LemmatizedDescription
		if normalizedDescription != nil {
			lemmatizedTokens, err := spanishprocessing.LemmatizeText(*normalizedDescription)
			if err != nil {
				return nil, err
			}
			var indexMappings []int
			var lemmatizedTextTokens []string
			for idx, lemmaToken := range lemmatizedTokens {
				if lemmaToken != nil {
					indexMappings = append(indexMappings, idx)
					lemmatizedTextTokens = append(lemmatizedTextTokens, lemmaToken.Str())
				}
			}
			lemmatizedDescription = &LemmatizedDescription{
				LemmatizedText: strings.Join(lemmatizedTextTokens, " "),
				IndexMappings:  indexMappings,
			}
		}
		return &TextMetadata{
			ReadabilityScore:      *readabilityScore,
			LemmatizedDescription: lemmatizedDescription,
		}, nil
	default:
		panic("unrecognized language")
	}
}
