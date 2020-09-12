package wordsmith

import "fmt"

type LemmaID string

type Lemma struct {
	ID           LemmaID
	Lemma        string
	PartOfSpeech string
	Language     LanguageCode
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
	ID           WordID  `db:"_id"`
	Word         string  `db:"word"`
	LemmaID      LemmaID `db:"lemma_id"`
	PartOfSpeech string  `db:"part_of_speech"`
	Language     string  `db:"language"`
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

func MustLanguageCodeForString(code string) LanguageCode {
	switch code {
	case LanguageCodeSpanish.Str():
		return LanguageCodeSpanish
	default:
		panic(fmt.Sprintf("unrecognized language code: %s", code))
	}
}
