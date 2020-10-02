package readability

import (
	"babblegraph/util/math/decimal"
	"babblegraph/util/storage"
	"babblegraph/wordsmith"
)

type CalculateReadabilityInput struct {
	Filename     storage.FileIdentifier
	LanguageCode wordsmith.LanguageCode
}

func CalculateReadability(input CalculateReadabilityInput) (*decimal.Number, error) {
	textBytes, err := storage.ReadFile(filename)
	if err != nil {
		return err
	}
	switch input.LanguageCode {
	case wordsmith.LanguageCodeSpanish:
		return calculateReadabilityForSpanish(string(textBytes))
	default:
		panic("unsupported language")
	}
}
