package text

import "strings"

func Tokenize(text string) []string {
	return strings.Split(text, " ")
}

func TokenizeUnique(text string) []string {
	tokenHash := make(map[string]bool)
	var tokenSet []string
	for _, token := range Tokenize(text) {
		if _, ok := tokenHash[token]; !ok {
			tokenHash[token] = true
			tokenSet = append(tokenSet, token)
		}
	}
	return tokenSet
}
