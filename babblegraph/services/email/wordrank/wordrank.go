package wordrank

import (
	"babblegraph/model/index"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"sort"

	"github.com/jmoiron/sqlx"
)

type RankedWord struct {
	LemmaID           wordsmith.LemmaID `json:"lemma_id"`
	Rank              int               `json:"rank"`
	TermFrequency     int64             `json:"term_frequency"`
	DocumentFrequency int64             `json:"document_frequency"`
}

func GetRankedWords(language wordsmith.LanguageCode) (map[wordsmith.LemmaID]RankedWord, error) {
	var orderedTerms []index.TermWithStats
	lemmaIDs, err := getAllLemmaIDsForLanguage(language)
	if err != nil {
		return nil, err
	}
	chunks := len(lemmaIDs) / 25000
	for i := 0; i < chunks; i++ {
		var orderedTermsChunk []index.TermWithStats
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			start := i * 25000
			end := (i + 1) * 25000
			if end > len(lemmaIDs) {
				end = len(lemmaIDs)
			}
			var err error
			orderedTermsChunk, err = index.GetOrderedTermsForLemmaIDs(tx, lemmaIDs[start:end])
			return err
		}); err != nil {
			return nil, err
		}
		orderedTerms = append(orderedTerms, orderedTermsChunk...)
	}
	sort.Sort(byCount(orderedTerms))
	out := make(map[wordsmith.LemmaID]RankedWord)
	for rank, termWithCount := range orderedTerms {
		out[termWithCount.TermID] = RankedWord{
			LemmaID:           termWithCount.TermID,
			Rank:              rank + 1,
			TermFrequency:     termWithCount.TotalCount,
			DocumentFrequency: termWithCount.DocumentCount,
		}
	}
	return out, nil
}

func getAllLemmaIDsForLanguage(language wordsmith.LanguageCode) ([]wordsmith.LemmaID, error) {
	lemmasForLanguage, err := wordsmith.GetAllLemmasForLanguage(language)
	if err != nil {
		return nil, err
	}
	var out []wordsmith.LemmaID
	for _, lemma := range lemmasForLanguage {
		out = append(out, lemma.ID)
	}
	return out, nil
}

type byCount []index.TermWithStats

func (c byCount) Len() int           { return len(c) }
func (c byCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byCount) Less(i, j int) bool { return c[i].TotalCount > c[j].TotalCount }
