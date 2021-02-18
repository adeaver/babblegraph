package language

import (
	"babblegraph/wordsmith"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type wrappedLemma struct {
	Lemma              wordsmith.Lemma
	DefinitionMappings []wordsmith.DefinitionMapping
	PartOfSpeech       wordsmith.PartOfSpeech
}

func getWrappedLemmas(wordText string) ([]wrappedLemma, error) {
	var lemmas []wordsmith.Lemma
	var definitionMappings []wordsmith.DefinitionMapping
	var partsOfSpeech []wordsmith.PartOfSpeech
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		var err error
		lemmas, err = wordsmith.GetLemmasByWordText(tx, wordsmith.SpanishUPCWikiCorpus, wordText)
		if err != nil {
			return err
		}
		var lemmaIDs []wordsmith.LemmaID
		var partOfSpeechIDs []wordsmith.PartOfSpeechID
		for _, l := range lemmas {
			lemmaIDs = append(lemmaIDs, l.ID)
			partOfSpeechIDs = append(partOfSpeechIDs, l.PartOfSpeechID)
		}
		definitionMappings, err = wordsmith.GetDefinitionMappingsForLemmaIDs(tx, wordsmith.SpanishOpenDefinitions, lemmaIDs)
		if err != nil {
			return err
		}
		partsOfSpeech, err = wordsmith.GetPartOfSpeechByIDs(tx, wordsmith.SpanishUPCWikiCorpus, partOfSpeechIDs)
		return err
	}); err != nil {
		return nil, err
	}
	partsOfSpeechByID := make(map[wordsmith.PartOfSpeechID]wordsmith.PartOfSpeech)
	for _, p := range partsOfSpeech {
		partsOfSpeechByID[p.ID] = p
	}
	wrappedLemmas := make(map[wordsmith.LemmaID]wrappedLemma)
	for _, lemma := range lemmas {
		partOfSpeech, ok := partsOfSpeechByID[lemma.PartOfSpeechID]
		if !ok {
			log.Println(fmt.Sprintf("no part of speech for id %s. continuing...", lemma.PartOfSpeechID))
			continue
		}
		wrappedLemmas[lemma.ID] = wrappedLemma{
			Lemma:        lemma,
			PartOfSpeech: partOfSpeech,
		}
	}
	for _, definitionMapping := range definitionMappings {
		wrappedLemma, ok := wrappedLemmas[definitionMapping.LemmaID]
		if !ok {
			log.Println(fmt.Sprintf("no lemma for id %s. continuing...", definitionMapping.LemmaID))
			continue
		}
		wrappedLemma.DefinitionMappings = append(wrappedLemma.DefinitionMappings, definitionMapping)
		wrappedLemmas[definitionMapping.LemmaID] = wrappedLemma
	}
	var out []wrappedLemma
	for _, wrappedLemma := range wrappedLemmas {
		out = append(out, wrappedLemma)
	}
	return out, nil
}

func (w wrappedLemma) ToAPI() lemma {
	var definitions []definition
	for _, d := range w.DefinitionMappings {
		definitions = append(definitions, convertDefinitionMappingToAPIFormat(d))
	}
	return lemma{
		Text: w.Lemma.LemmaText,
		ID:   w.Lemma.ID,
		PartOfSpeech: partOfSpeech{
			ID:   w.PartOfSpeech.ID,
			Name: w.PartOfSpeech.Code.ToDisplayCategory(),
		},
		Definitions: definitions,
	}
}

func convertDefinitionMappingToAPIFormat(d wordsmith.DefinitionMapping) definition {
	return definition{
		Text:      d.EnglishDefinition,
		ExtraInfo: d.ExtraInfo,
	}
}
