package documents

import "babblegraph/wordsmith"

var filteredWordsForLanguageCode = map[wordsmith.LanguageCode][]string{
	wordsmith.LanguageCodeSpanish: []string{
		"arma",
		"armas",
		"homicidio",
		"homicidios",
		"asesinato",
		"asesinatos",
		"muerte",
		"muertos",
		"disparó",
		"dispararon",
		"mató",
		"mataron",
		"morió",
		"morieron",
	},
}
