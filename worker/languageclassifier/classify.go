package languageclassifier

import (
	"babblegraph/worker/storage"
	"babblegraph/worker/wordsmith"
	"strings"
)

func ClassifyLanguageForFile(filename storage.FileIdentifier) (*wordsmith.LanguageCode, error) {
	textBytes, err := storage.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return classify(string(textBytes))
}

const minimumPercentage int64 = 65

var classifiers []classifier = []classifier{
	spanishClassifierImpl,
}

func classify(text string) (*wordsmith.LanguageCode, error) {
	var currentHighest *int64
	var currentClassifier *classifier
	tokens := tokenize(text)
	for _, classifier := range classifiers {
		classification, err := classifier.Classify(tokens)
		if err != nil {
			return nil, err
		}
		if classification == nil || *classification < minimumPercentage {
			continue
		}
		if currentHighest == nil || *classification > *currentHighest {
			currentHighest = classification
			currentClassifier = &classifier
		}
	}
	if currentClassifier == nil {
		return nil, nil
	}
	lang := (*currentClassifier).GetLanguage()
	return &lang, nil
}

func tokenize(text string) map[string]int64 {
	tokens := make(map[string]int64)
	lines := strings.Split(text, "\n")
	for _, l := range lines {
		lineTokens := strings.Split(l, " ")
		for _, token := range lineTokens {
			count, ok := tokens[token]
			if !ok {
				count = 0
			}
			tokens[token] = count + 1
		}
	}
	return tokens
}
