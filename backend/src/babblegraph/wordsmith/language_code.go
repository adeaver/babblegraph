package wordsmith

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
