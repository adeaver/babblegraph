package syllable

import (
	"babblegraph/wordsmith"
	"fmt"
)

func CountSyllablesInWord(language wordsmith.LanguageCode, word string) (*int64, error) {
	switch language {
	case wordsmith.LanguageCodeSpanish:
		return countSyllablesForSpanish(word)
	default:
		return nil, fmt.Errorf("invalid language %s", language)
	}
}
