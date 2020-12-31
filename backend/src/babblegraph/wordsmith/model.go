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
