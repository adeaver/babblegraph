package labels

import (
	"babblegraph/wordsmith"
)

// This package defines a set of hardcoded labels
// in order to reduce the complexity of finding relevant
// documents. Everything in this package currently assumes
// that the language is Spanish.

type LabelName string

const (
	LabelNameScience    LabelName = "Science"
	LabelNameTechnology LabelName = "Technology"
)

var keywordSearchTermsForLabelNames = map[LabelName][]string{
	LabelNameScience: []string{
		"cienca",
		"biología",
		"química",
		"física",
		"matemáticas",
	},
	LabelNameTechnology: []string{
		"tecnología",
		"computadora",
		"máquina",
	},
}

func GetLemmaIDsForLabelNames() (map[LabelName][]wordsmith.LemmaID, error) {
	out := make(map[LabelName][]wordsmith.LemmaID)
	for labelName, searchTerm := range keywordSearchTermsForLabelNames {
		lemmas, err := wordsmith.LookupLemmas(searchTerm, wordsmith.LanguageCodeSpanish)
		if err != nil {
			return nil, err
		}
		var lemmaIDs []wordsmith.LemmaID
		for _, lemma := range lemmas {
			lemmaIDs = append(lemmaIDs, lemma.ID)
		}
		out[labelName] = lemmaIDs
	}
	return out, nil
}
