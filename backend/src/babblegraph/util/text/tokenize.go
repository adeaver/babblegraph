package text

import "strings"

func Tokenize(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		out = append(out, strings.Split(line, " ")...)
	}
	return out
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
