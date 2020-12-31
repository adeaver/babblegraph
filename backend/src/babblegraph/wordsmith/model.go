package wordsmith

import "fmt"

type CorpusID string

type LemmaID string

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

type WordRankingID string

type WordRanking struct {
	ID            WordRankingID
	LanguageCode  LanguageCode
	CorpusID      CorpusID
	Word          string
	CorpusRanking int64
	CorpusCount   int64
}

type dbWordRanking struct {
	ID            WordRankingID `db:"_id"`
	LanguageCode  LanguageCode  `db:"language"`
	CorpusID      CorpusID      `db:"corpus_id"`
	Word          string        `db:"word"`
	CorpusRanking int64         `db:"ranking"`
	CorpusCount   int64         `db:"count"`
}

func (d dbWordRanking) ToNonDB() WordRanking {
	return WordRanking{
		ID:            d.ID,
		LanguageCode:  d.LanguageCode,
		CorpusID:      d.CorpusID,
		Word:          d.Word,
		CorpusRanking: d.CorpusRanking,
		CorpusCount:   d.CorpusCount,
	}
}
