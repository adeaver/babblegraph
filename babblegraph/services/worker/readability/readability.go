package readability

import (
	"babblegraph/model/documents"
	"babblegraph/util/storage"
	"babblegraph/wordsmith"
)

type CalculateReadabilityInput struct {
	DocumentID   documents.DocumentID
	Filename     storage.FileIdentifier
	LanguageCode wordsmith.LanguageCode
}

func CalculateReadability(input CalculateReadabilityInput) error {
	textBytes, err := storage.ReadFile(filename)
	if err != nil {
		return err
	}
	switch input.LanguageCode {
	case wordsmith.LanguageCodeSpanish:
		readabilityScore, err := calculateReadabilityForSpanish(string(textBytes))
		if err != nil {
			return err
		}
		// insert
	default:
		panic("unsupported language")
	}
}
