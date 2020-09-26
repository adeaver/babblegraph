package wordrank

import (
	"babblegraph/model/index"
	"babblegraph/util/database"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type RankedWord struct {
	LemmaID           wordsmith.LemmaID `json:"lemma_id"`
	Rank              int64             `json:"rank"`
	TermFrequency     int64             `json:"term_frequency"`
	DocumentFrequency int64             `json:"document_frequency"`
}

func GetRankedWords(language wordsmith.LanguageCode) ([]RankedWord, error) {
	var orderedTerms []index.TermWithStats
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		orderedTerms, err = index.GetOrderedTermsForLanguage(language)
		return err
	}); err != nil {
		return nil, err
	}
	var out []RankedWord
	for rank, termWithCount := range orderedTerms {
		out = append(out, RankedWord{
			LemmaID:           termWithCount.TermID,
			Rank:              rank + 1,
			TermFrequency:     termWithCount.TotalCount,
			DocumentFrequency: termWithCount.DocumentFrequency,
		})
	}
	return out
}
