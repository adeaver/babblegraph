package wordsmith

import "github.com/jmoiron/sqlx"

type WordID string

type Word struct {
	ID             WordID
	Language       LanguageCode
	CorpusID       CorpusID
	PartOfSpeechID PartOfSpeechID
	LemmaID        LemmaID
	WordText       string
}

type dbWord struct {
	ID             WordID         `db:"_id"`
	Language       LanguageCode   `db:"language"`
	CorpusID       CorpusID       `db:"corpus_id"`
	PartOfSpeechID PartOfSpeechID `db:"part_of_speech_id"`
	LemmaID        LemmaID        `db:"lemma_id"`
	WordText       string         `db:"word_text"`
}

func (d dbWord) ToNonDB() Word {
	return Word{
		ID:             d.ID,
		Language:       d.Language,
		CorpusID:       d.CorpusID,
		PartOfSpeechID: d.PartOfSpeechID,
		LemmaID:        d.LemmaID,
		WordText:       d.WordText,
	}
}

const wordsForTextQuery = "SELECT * FROM words WHERE word_text IN (?) AND corpus = $1"

func GetWordsByText(tx *sqlx.Tx, corpus CorpusID, words []string) ([]Word, error) {
	query, args, err := sqlx.In(wordsForTextQuery, words, corpus)
	if err != nil {
		return nil, nil
	}
	sql := tx.Rebind(query)
	var matches []dbWord
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, nil
	}
	var out []Word
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
