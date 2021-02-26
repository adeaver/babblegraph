package wordsmith

import "github.com/jmoiron/sqlx"

type WordBigramCountID string

type WordBigramCount struct {
	ID         WordBigramCountID
	Language   LanguageCode
	CorpusID   CorpusID
	FirstWord  BigramWord
	SecondWord BigramWord
	Count      int64
}

type BigramWord struct {
	Text    string
	LemmaID LemmaID
}

type dbWordBigramCount struct {
	ID                WordBigramCountID `db:"_id"`
	Language          LanguageCode      `db:"language"`
	CorpusID          CorpusID          `db:"corpus_id"`
	FirstWordText     string            `db:"first_word_text"`
	FirstWordLemmaID  LemmaID           `db:"first_word_lemma_id"`
	SecondWordText    string            `db:"second_word_text"`
	SecondWordLemmaID LemmaID           `db:"second_word_lemma_id"`
	Count             int64             `db:"count"`
}

func (d dbWordBigramCount) ToNonDB() WordBigramCount {
	return WordBigramCount{
		ID:       d.ID,
		Language: d.Language,
		CorpusID: d.CorpusID,
		FirstWord: BigramWord{
			Text:    d.FirstWordText,
			LemmaID: d.FirstWordLemmaID,
		},
		SecondWord: BigramWord{
			Text:    d.SecondWordText,
			LemmaID: d.SecondWordLemmaID,
		},
		Count: d.Count,
	}
}

const wordBigramCountForWordQuery = "SELECT * FROM word_bigram_counts WHERE corpus_id = $1 AND (first_word_text = $2 AND second_word_text = $3)"

func GetWordBigramCountsByWordText(tx *sqlx.Tx, corpusID CorpusID, firstWord string, secondWord string) ([]WordBigramCount, error) {
	var matches []dbWordBigramCount
	if err := tx.Select(&matches, wordBigramCountForWordQuery, corpusID, firstWord, secondWord); err != nil {
		return nil, err
	}
	var out []WordBigramCount
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
