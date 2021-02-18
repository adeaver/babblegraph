package wordsmith

import "github.com/jmoiron/sqlx"

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
