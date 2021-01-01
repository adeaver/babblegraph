package textprocessing

import (
	"babblegraph/model/documents"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"strings"
)

func getWordStatsForText(languageCode wordsmith.LanguageCode, normalizedText string) (*documents.WordStatsVersion1, error) {
	tokens := tokenizeText(normalizedText)
	uniqueWords := getUniqueWordsForText(tokens)
	tokenCounts := getTokenCounts(tokens)
	rankings, err := wordsmith.GetSortedRankingsForWords(languageCode, uniqueWords)
	if err != nil {
		return nil, err
	}
	wordExclusions := extractWordExclusionsFromRankings(rankings)
	out := &documents.WordStatsVersion1{
		AverageWordRanking:  calculateMeanWordRanking(tokenCounts, rankings),
		MedianWordRanking:   calculateMedianWordRanking(tokenCounts, rankings),
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

func getTokenCounts(tokenizedText []string) map[string]int64 {
	out := make(map[string]int64)
	for _, t := range tokenizedText {
		out[t]++
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

func calculateMedianWordRanking(tokenCounts map[string]int64, rankings []wordsmith.WordRanking) int64 {
	var sortedRankings []int64
	for _, r := range rankings {
		if count, ok := tokenCounts[r.Word]; ok {
			var i int64 = 0
			for ; i < count; i++ {
				sortedRankings = append(sortedRankings, r.CorpusRanking)
			}
		}
	}
	midIdx := len(sortedRankings) / 2
	if len(sortedRankings)%2 == 0 {
		return (sortedRankings[midIdx] + sortedRankings[midIdx+1]) / 2
	}
	return sortedRankings[midIdx]
}

func calculateMeanWordRanking(tokenCounts map[string]int64, rankings []wordsmith.WordRanking) int64 {
	var totalRanking, rankedWordsCount decimal.Number
	for _, r := range rankings {
		if count, ok := tokenCounts[r.Word]; ok {
			totalRanking.Add(decimal.FromInt64(count * r.CorpusRanking))
			rankedWordsCount.Add(decimal.FromInt64(count))
		}
	}
	return totalRanking.Divide(rankedWordsCount).ToInt64Rounded()
}
