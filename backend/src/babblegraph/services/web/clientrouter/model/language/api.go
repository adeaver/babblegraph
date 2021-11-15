package language

import "babblegraph/wordsmith"

type Lemma struct {
	Text         string            `json:"text"`
	ID           wordsmith.LemmaID `json:"id"`
	PartOfSpeech PartOfSpeech      `json:"part_of_speech"`
	Definitions  []Definition      `json:"definitions"`
}

type PartOfSpeech struct {
	ID   wordsmith.PartOfSpeechID `json:"id"`
	Name string                   `json:"name"`
}

type Definition struct {
	Text      string  `json:"text"`
	ExtraInfo *string `json:"extra_info,omitempty"`
}

func GetLemmasForWordText(wordText string) ([]Lemma, error) {
	wrappedLemmas, err := getWrappedLemmasForWordText(wordText)
	if err != nil {
		return nil, err
	}
	var out []Lemma
	for _, wrappedLemma := range wrappedLemmas {
		out = append(out, wrappedLemma.ToAPI())
	}
	return out, nil
}

func GetLemmasForLemmaIDs(lemmaIDs []wordsmith.LemmaID) ([]Lemma, error) {
	wrappedLemmas, err := getWrappedLemmasForLemmaIDs(lemmaIDs)
	if err != nil {
		return nil, err
	}
	var out []Lemma
	for _, wrappedLemma := range wrappedLemmas {
		out = append(out, wrappedLemma.ToAPI())
	}
	return out, nil
}

func (w wrappedLemma) ToAPI() Lemma {
	var definitions []Definition
	for _, d := range w.DefinitionMappings {
		definitions = append(definitions, convertDefinitionMappingToAPIFormat(d))
	}
	return Lemma{
		Text: w.Lemma.LemmaText,
		ID:   w.Lemma.ID,
		PartOfSpeech: PartOfSpeech{
			ID:   w.PartOfSpeech.ID,
			Name: w.PartOfSpeech.Code.ToDisplayCategory(),
		},
		Definitions: definitions,
	}
}

func convertDefinitionMappingToAPIFormat(d wordsmith.DefinitionMapping) Definition {
	return Definition{
		Text:      d.EnglishDefinition,
		ExtraInfo: d.ExtraInfo,
	}
}
