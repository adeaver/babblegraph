package wordsmith

import "fmt"

type CorpusID string

type LemmaID string

type Lemma struct {
	ID           LemmaID
	Lemma        string
	PartOfSpeech string
	Language     LanguageCode
}

type dbLemma struct {
	ID             LemmaID  `db:"_id"`
	CorpusID       CorpusID `db:"corpus_id"`
	LemmaText      string   `db:"lemma_text"`
	PartOfSpeechID string   `db:"part_of_speech_id"`
	Language       string   `db:"language"`
}

func (d dbLemma) ToNonDB() Lemma {
	return Lemma{
		ID:           d.ID,
		Lemma:        d.LemmaText,
		PartOfSpeech: d.PartOfSpeechID,
		Language:     MustLanguageCodeForString(d.Language),
	}
}

type WordID string

type Word struct {
	ID           WordID
	Word         string
	LemmaID      LemmaID
	PartOfSpeech string
	Language     LanguageCode
}

type dbWord struct {
	ID           WordID   `db:"_id"`
	CorpusID     CorpusID `db:"corpus_id"`
	Word         string   `db:"word_text"`
	LemmaID      LemmaID  `db:"lemma_id"`
	PartOfSpeech string   `db:"part_of_speech_id"`
	Language     string   `db:"language"`
}

func (d *dbWord) ToWord() Word {
	return Word{
		ID:           d.ID,
		Word:         d.Word,
		LemmaID:      d.LemmaID,
		PartOfSpeech: d.PartOfSpeech,
		Language:     MustLanguageCodeForString(d.Language),
	}
}

type LanguageCode string

const (
	LanguageCodeSpanish LanguageCode = "es"
)

func (c LanguageCode) Str() string {
	return string(c)
}

func (c LanguageCode) Ptr() *LanguageCode {
	return &c
}

func MustLanguageCodeForString(code string) LanguageCode {
	switch code {
	case LanguageCodeSpanish.Str():
		return LanguageCodeSpanish
	default:
		panic(fmt.Sprintf("unrecognized language code: %s", code))
	}
}
