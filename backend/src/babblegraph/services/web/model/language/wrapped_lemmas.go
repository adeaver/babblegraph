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

func getWrappedLemmasForWordText(wordText string) ([]wrappedLemma, error) {
	var lemmas []wordsmith.Lemma
	var definitionMappings []wordsmith.DefinitionMapping
	var partsOfSpeech []wordsmith.PartOfSpeech
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		var err error
		lemmas, err = wordsmith.GetLemmasByWordText(tx, wordsmith.SpanishUPCWikiCorpus, wordText)
		if err != nil {
			return err
		}
		definitionMappings, partsOfSpeech, err = getDefinitionsAndPartsOfSpeechForLemmas(tx, lemmas)
		return err
	}); err != nil {
		return nil, err
	}
	return getWrappedLemmas(getWrappedLemmasInput{
		lemmas:             lemmas,
		partsOfSpeech:      partsOfSpeech,
		definitionMappings: definitionMappings,
	}), nil
}

func getWrappedLemmasForLemmaIDs(lemmaIDs []wordsmith.LemmaID) ([]wrappedLemma, error) {
	var lemmas []wordsmith.Lemma
	var definitionMappings []wordsmith.DefinitionMapping
	var partsOfSpeech []wordsmith.PartOfSpeech
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		var err error
		lemmas, err = wordsmith.GetLemmasByIDs(tx, lemmaIDs)
		if err != nil {
			return err
		}
		definitionMappings, partsOfSpeech, err = getDefinitionsAndPartsOfSpeechForLemmas(tx, lemmas)
		return err
	}); err != nil {
		return nil, err
	}
	return getWrappedLemmas(getWrappedLemmasInput{
		lemmas:             lemmas,
		partsOfSpeech:      partsOfSpeech,
		definitionMappings: definitionMappings,
	}), nil
}

func getDefinitionsAndPartsOfSpeechForLemmas(tx *sqlx.Tx, lemmas []wordsmith.Lemma) ([]wordsmith.DefinitionMapping, []wordsmith.PartOfSpeech, error) {
	var lemmaIDs []wordsmith.LemmaID
	var partOfSpeechIDs []wordsmith.PartOfSpeechID
	for _, l := range lemmas {
		lemmaIDs = append(lemmaIDs, l.ID)
		partOfSpeechIDs = append(partOfSpeechIDs, l.PartOfSpeechID)
	}
	definitionMappings, err := wordsmith.GetDefinitionMappingsForLemmaIDs(tx, wordsmith.SpanishOpenDefinitions, lemmaIDs)
	if err != nil {
		return nil, nil, err
	}
	partsOfSpeech, err := wordsmith.GetPartOfSpeechByIDs(tx, wordsmith.SpanishUPCWikiCorpus, partOfSpeechIDs)
	if err != nil {
		return nil, nil, err
	}
	return definitionMappings, partsOfSpeech, nil
}

type getWrappedLemmasInput struct {
	partsOfSpeech      []wordsmith.PartOfSpeech
	definitionMappings []wordsmith.DefinitionMapping
	lemmas             []wordsmith.Lemma
}

func getWrappedLemmas(input getWrappedLemmasInput) []wrappedLemma {
	partsOfSpeechByID := make(map[wordsmith.PartOfSpeechID]wordsmith.PartOfSpeech)
	for _, p := range input.partsOfSpeech {
		partsOfSpeechByID[p.ID] = p
	}
	wrappedLemmas := make(map[wordsmith.LemmaID]wrappedLemma)
	for _, lemma := range input.lemmas {
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
	for _, definitionMapping := range input.definitionMappings {
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
		if len(wrappedLemma.DefinitionMappings) == 0 {
			continue
		}
		out = append(out, wrappedLemma)
	}
	return out
}
