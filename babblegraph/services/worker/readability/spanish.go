package readability

import (
	"babblegraph/util/language/syllable"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"strings"
)

func calculateReadabilityForSpanish(text string) (*decimal.Number, error) {
	sentences := strings.Split(text, "\n")
	var wordCount, syllableCount, sentenceCount decimal.Number
	for _, sentence := range sentences {
		sentenceCount = sentenceCount.Add(1)
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
	syllableTerm := decimal.FromFloat64(60.0).Multiply(syllableCount.Divide(wordCount))
	wordTerm := decimal.FromFloat64(1.02).Multiply(sentenceCount.Divide(wordCount))
	return decimal.FromFloat64(206.84).Subtract(syllableTerm).Subtract(wordTerm)
}
