package text

import (
	"babblegraph/wordsmith"
	"fmt"
	"strings"
)

var nonTitleCaseWords = map[wordsmith.LanguageCode][]string{
	wordsmith.LanguageCodeSpanish: []string{
		"de",
		"del",
		"a",
		"al",
		"el",
		"los",
		"las",
		"la",
	},
}

func ToTitleCaseForLanguage(text string, languageCode wordsmith.LanguageCode) string {
	tokens := Tokenize(text)
	var out []string
	for idx, token := range tokens {
		lowerCaseToken := strings.ToLower(token)
		if idx == 0 || !isNonTitleCaseWord(lowerCaseToken, languageCode) {
			out = append(out, titleCaseToken(lowerCaseToken))
			continue
		}
		out = append(out, lowerCaseToken)
	}
	return strings.Join(out, " ")
}

func titleCaseToken(text string) string {
	uppercaseFirstLetter := strings.ToUpper(string(text[0]))
	return fmt.Sprintf("%s%s", uppercaseFirstLetter, text[1:])
}

func isNonTitleCaseWord(text string, languageCode wordsmith.LanguageCode) bool {
	nonTitleCaseWordsForLanguage, ok := nonTitleCaseWords[languageCode]
	if !ok {
		return false
	}
	for _, nonTitleCaseWord := range nonTitleCaseWordsForLanguage {
		if nonTitleCaseWord == text {
			return true
		}
	}
	return false
}
