package readability

import (
	"babblegraph/util/language/syllable"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"strings"
)

func calculateReadabilityForSpanish(text string) (*decimal.Number, error) {
	sentences := strings.Split(text, "\n")
	var wordCount, syllableCount decimal.Number
	for _, sentence := range sentences {
		words := strings.Split(sentence, "\n")
		wordCount = wordCount.Add(decimal.FromInt64(len(words)))
		for _, word := range words {
			count, err := syllable.CountSyllablesInWord(wordsmith.LanguageCode, word)
			if err != nil {
				return nil, err
			}
			syllableCount = syllableCount.Add(decimal.FromInt64(*count))
		}
	}
	// score := decimal.FromFloat64(
	// return
}
