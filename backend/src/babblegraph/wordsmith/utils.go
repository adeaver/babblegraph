package wordsmith

import "strings"

// TODO: get rid of this with new website model
// A language is considered supported if there is a code available for it in wordsmith
var htmlLanguageValueToLanguageCode = map[string]LanguageCode{
	"es":    LanguageCodeSpanish,
	"es-es": LanguageCodeSpanish,
	"es-mx": LanguageCodeSpanish,
	"es-co": LanguageCodeSpanish,
	"es-gt": LanguageCodeSpanish,
}

func LookupLanguageCodeForLanguageLabel(languageLabel string) *LanguageCode {
	if code, ok := htmlLanguageValueToLanguageCode[strings.ToLower(languageLabel)]; ok {
		return &code
	}
	return nil
}

func LookupLanguageLabelsForLanguageCode(languageCode LanguageCode) []string {
	var out []string
	for label, code := range htmlLanguageValueToLanguageCode {
		if code == languageCode {
			out = append(out, label)
		}
	}
	return out
}
