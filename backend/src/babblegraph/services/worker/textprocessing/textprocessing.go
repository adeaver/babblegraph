package textprocessing

import (
	"babblegraph/services/worker/textprocessing/spanishprocessing"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
)

type TextMetadata struct {
	ReadabilityScore      decimal.Number
	LemmatizedDescription *string
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
		var lemmatizedDescription *string
		if normalizedDescription != nil {
			lemmatizedDescription = ptr.String(spanishprocessing.LemmatizeText(*normalizedDescription))
		}
		return &TextMetadata{
			ReadabilityScore:      *readabilityScore,
			LemmatizedDescription: lemmatizedDescription,
		}, nil
	default:
		panic("unrecognized language")
	}
}
