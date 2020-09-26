package wordrank

import "github.com/adeaver/babblegraph/lib/wordsmith"

type RankedWord struct {
	LemmaID wordsmith.LemmaID `json:"lemma_id"`
	Rank    int64             `json:"rank"`
}

func RankWords(language wordsmith.LanguageCode) []RankedWord {

}
