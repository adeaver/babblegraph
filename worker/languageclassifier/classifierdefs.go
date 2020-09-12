package languageclassifier

import (
	"babblegraph/worker/wordsmith"
)

type classifier interface {
	Classify(tokens map[string]int64) (*int, error)
	GetLanguage() wordsmith.LanguageCode
}

type spanishClassifier struct{}

func (s spanishClassifier) Classify(tokens map[string]int64) (*int, error) {
	var totalTokens int
	var keys []string
	for token, count := range tokens {
		keys = append(keys, token)
		totalTokens += count
	}
	words, err := wordsmith.GetWords(keys, wordsmith.LanguageCodeSpanish)
	if err != nil {
		return err
	}
	var count int
	for _, w := range words {
		if c, ok := tokens[w.Word]; ok {
			count += c
		}
	}
	percent := float64(count) / float64(totalTokens)
	percentAsInt := int(percent * 100)
	return &percentAsInt, nil
}

func (s spanishClassifier) GetLanguage() wordsmith.LanguageCode {
	return wordsmith.LanguageCodeSpanish
}
