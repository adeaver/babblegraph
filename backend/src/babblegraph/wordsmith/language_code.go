package wordsmith

import "fmt"

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

func GetLanguageCodeFromString(s string) (*LanguageCode, error) {
	switch s {
	case LanguageCodeSpanish.Str():
		return LanguageCodeSpanish.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unrecognized language code: %s", s)
	}
}
