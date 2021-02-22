package wordsmith

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type LemmaID string

func (l *LemmaID) Str() string {
	return string(*l)
}

func (l LemmaID) Ptr() *LemmaID {
	return &l
}

type Lemma struct {
	ID             LemmaID
	Language       LanguageCode
	CorpusID       CorpusID
	PartOfSpeechID PartOfSpeechID
	LemmaText      string
}

type dbLemma struct {
	ID             LemmaID        `db:"_id"`
	Language       LanguageCode   `db:"language"`
	CorpusID       CorpusID       `db:"corpus_id"`
	PartOfSpeechID PartOfSpeechID `db:"part_of_speech_id"`
	LemmaText      string         `db:"lemma_text"`
}

func (d dbLemma) ToNonDB() Lemma {
	return Lemma{
		ID:             d.ID,
		Language:       d.Language,
		CorpusID:       d.CorpusID,
		PartOfSpeechID: d.PartOfSpeechID,
		LemmaText:      d.LemmaText,
	}
}

const lemmasByWordTextQuery = "SELECT * FROM lemmas WHERE _id IN (SELECT lemma_id FROM words WHERE word_text = $1 AND corpus_id = $2)"

func GetLemmasByWordText(tx *sqlx.Tx, corpus CorpusID, wordText string) ([]Lemma, error) {
	var matches []dbLemma
	if err := tx.Select(&matches, lemmasByWordTextQuery, wordText, corpus); err != nil {
		return nil, err
	}
	var out []Lemma
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

const getLemmaByIDQuery = "SELECT * FROM lemmas WHERE _id = $1"

// This function doesn't need a corpus id since lemmas are unique across corpora
func GetLemmaByID(tx *sqlx.Tx, id LemmaID) (*Lemma, error) {
	var matches []dbLemma
	if err := tx.Select(&matches, getLemmaByIDQuery, id); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("expecting exactly one match, but got %d", len(matches))
	}
	l := matches[0].ToNonDB()
	return &l, nil
}

const getLemmasByIDs = "SELECT * FROM lemmas WHERE _id IN (?)"

func GetLemmasByIDs(tx *sqlx.Tx, ids []LemmaID) ([]Lemma, error) {
	query, args, err := sqlx.In(getLemmasByIDs, ids)
	if err != nil {
		return nil, nil
	}
	sql := tx.Rebind(query)
	var matches []dbLemma
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, nil
	}
	var out []Lemma
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
