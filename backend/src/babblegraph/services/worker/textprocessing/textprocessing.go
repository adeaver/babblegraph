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
		return &TextMetadata{
			ReadabilityScore:      *readabilityScore,
			LemmatizedDescription: lemmatizedDescription,
		}, nil
	default:
		panic("unrecognized language")
	}
}
