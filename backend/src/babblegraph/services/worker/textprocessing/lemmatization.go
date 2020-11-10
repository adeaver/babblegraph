package textprocessing

func LemmatizeBody(languageCode wordsmith.LanguageCode, normalizedText string)
	tokens := collectWords(normalizedText)
	lemmaMap, err := getLemmaMapForTokens(tokens, languageCode)
	if err != nil {
		return nil, err
	}
	lemmaText := rewriteNormalizedTextWithLemmas(normalizedText, lemmaMap)
	return &lemmaText, nil
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

func getLemmaMapForTokens(tokens []string, languageCode wordsmith.LanguageCode) (map[string][]wordsmith.LemmaID, error) {
	words, err := wordsmith.GetWords(tokens, languageCode)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]wordsmith.LemmaID)
	for _, w := range words {
		lemmasForWord, ok := out[w.Word]
		if ok {
			// TODO: I need to do something more clever here
			log.Println(fmt.Sprintf("Word %s has duplicate. Using both lemmas...", w.Word))
		}
		out[w.Word] = append(lemmasForWord, w.LemmaID)
	}
	return out, nil
}

func rewriteNormalizedTextWithLemmas(normalizedText string, lemmaMap map[string][]wordsmith.LemmaID) string {
	var outTokens []string
	tokens := strings.Split(normalizedText, " ")
	for _, token := range tokens {
		if lemmaIDs, ok := lemmaMap[token]; ok {
			for _, lemmaID := range lemmaIDs {
				outTokens = append(outTokens, string(lemmaID))
			}
		}
	}
	return strings.Join(outTokens, " ")
}
