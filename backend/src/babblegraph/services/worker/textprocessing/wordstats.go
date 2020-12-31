package textprocessing

import (
	"babblegraph/model/documents"
	"babblegraph/wordsmith"
	"strings"
)

func GetWordStatsForText(languageCode wordsmith.LanguageCode, normalizedText string) (*documents.WordStatsVersion1, error) {
	tokens := tokenizeText(normalizedText)
	uniqueWords := getUniqueWordsForText(tokens)
	rankings, err := wordsmith.GetSortedRankingsForWords(languageCode, uniqueWords)
	if err != nil {
		return nil, err
	}
	wordExclusions := extractWordExclusionsFromRankings(rankings)
	out := &documents.WordStatsVersion1{
		TotalNumberOfWords:  int64(len(tokens)),
		NumberOfUniqueWords: int64(len(uniqueWords)),
	}
	if len(rankings) > 0 {
		out.LeastFrequentWordRanking = rankings[len(rankings)-1].CorpusRanking
	}
	for i := len(wordExclusions) - 1; i >= 0; i-- {
		switch i {
		case len(wordExclusions) - 1:
			out.LeastFrequentWordExclusion = &wordExclusions[i]
		case len(wordExclusions) - 2:
			out.SecondLeastFrequentWordExclusion = &wordExclusions[i]
		case len(wordExclusions) - 3:
			out.ThirdLeastFrequentWordExclusion = &wordExclusions[i]
		}
	}
	return out, err
}

func tokenizeText(normalizedText string) []string {
	lines := strings.Split(normalizedText, "\n")
	var out []string
	for _, line := range lines {
		out = append(out, strings.Split(line, " ")...)
	}
	return out
}

func getUniqueWordsForText(tokenizedText []string) []string {
	var tokenSet map[string]bool
	var out []string
	for _, token := range tokenizedText {
		if _, ok := tokenSet[token]; !ok {
			tokenSet[token] = true
			out = append(out, token)
		}
	}
	return out
}

func extractWordExclusionsFromRankings(rankings []wordsmith.WordRanking) []documents.WordExclusion {
	var out []documents.WordExclusion
	for i := len(rankings) - 1; i >= 0; i-- {
		var newRanking int64 = 0
		if i > 0 {
			newRanking = rankings[i-1].CorpusRanking
		}
		out = append(out, documents.WordExclusion{
			WordText:                        rankings[i].Word,
			LeastFrequentRankingWithoutWord: newRanking,
		})
	}
	return out
}
