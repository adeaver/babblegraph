package wordsmith

// A language is considered supported if there is a code available for it in wordsmith
var htmlLanguageValueToLanguageCode = map[string]LanguageCode{
	"es": LanguageCodeSpanish,
}

func LookupLanguageCodeForLanguageLabel(languageLabel string) *LanguageCode {
	if code, ok := htmlLanguageValueToLanguageCode[languageLabel]; ok {
		return &code
	}
	return nil
}
