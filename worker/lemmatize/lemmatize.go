package lemmatize

import (
	"babblegraph/worker/storage"
	"babblegraph/worker/wordsmith"
	"fmt"
	"log"
	"strings"
)

func LemmatizeWordsForFile(filename storage.FileIdentifier, languageCode wordsmith.LanguageCode) (*storage.FileIdentifier, error) {
	textBytes, err := storage.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	normalizedText := string(textBytes)
	tokens := collectWords(normalizedText)
	lemmaMap, err := getLemmaMapForTokens(tokens, languageCode)
	if err != nil {
		return nil, err
	}
	lemmaText := rewriteNormalizedTextWithLemmas(normalizedText, lemmaMap)
	log.Println(lemmaText)
	return storage.WriteFile("txt", lemmaText)
}

func collectWords(normalizedText string) []string {
	var out []string
	seenTokens := make(map[string]bool)
	tokens := strings.Split(normalizedText, " ")
	for _, token := range tokens {
		if _, ok := seenTokens[token]; !ok {
			seenTokens[token] = true
			out = append(out, token)
		}
	}
	return out
}

func getLemmaMapForTokens(tokens []string, languageCode wordsmith.LanguageCode) (map[string]wordsmith.LemmaID, error) {
	words, err := wordsmith.GetWords(tokens, languageCode)
	if err != nil {
		return nil, err
	}
	out := make(map[string]wordsmith.LemmaID)
	for _, w := range words {
		if _, ok := out[w.Word]; ok {
			// TODO: I need to do something more clever here
			log.Println(fmt.Sprintf("Word %s has duplicate. Replacing..."))
		}
		out[w.Word] = w.LemmaID
	}
	return out, nil
}

func rewriteNormalizedTextWithLemmas(normalizedText string, lemmaMap map[string]wordsmith.LemmaID) string {
	var outTokens []string
	tokens := strings.Split(normalizedText, " ")
	for _, token := range tokens {
		if lemmaID, ok := lemmaMap[token]; ok {
			outTokens = append(outTokens, string(lemmaID))
		}
	}
	return strings.Join(outTokens, " ")
}
