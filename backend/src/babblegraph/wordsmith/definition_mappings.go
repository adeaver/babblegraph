package wordsmith

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DefinitionMappingID string

type DefinitionMapping struct {
	ID                DefinitionMappingID
	Language          LanguageCode
	CorpusID          CorpusID
	LemmaID           LemmaID
	PartOfSpeechID    PartOfSpeechID
	EnglishDefinition string
	ExtraInfo         *string
}

type dbDefinitionMapping struct {
	ID                DefinitionMappingID `db:"_id"`
	Language          LanguageCode        `db:"language"`
	CorpusID          CorpusID            `db:"corpus_id"`
	LemmaID           LemmaID             `db:"lemma_id"`
	PartOfSpeechID    PartOfSpeechID      `db:"part_of_speech_id"`
	EnglishDefinition string              `db:"english_definition"`
	ExtraInfo         *string             `db:"extra_info"`
}

func (d dbDefinitionMapping) ToNonDB() DefinitionMapping {
	return DefinitionMapping{
		ID:                d.ID,
		Language:          d.Language,
		CorpusID:          d.CorpusID,
		LemmaID:           d.LemmaID,
		PartOfSpeechID:    d.PartOfSpeechID,
		EnglishDefinition: d.EnglishDefinition,
		ExtraInfo:         d.ExtraInfo,
	}
}

const definitionMappingsForLemmaIDsQuery = "SELECT * FROM definition_mappings WHERE corpus_id = '%s' AND lemma_id IN (?)"

func GetDefinitionMappingsForLemmaIDs(tx *sqlx.Tx, corpusID CorpusID, lemmaIDs []LemmaID) ([]DefinitionMapping, error) {
	query, args, err := sqlx.In(fmt.Sprintf(definitionMappingsForLemmaIDsQuery, corpusID), lemmaIDs)
	if err != nil {
		return nil, nil
	}
	sql := tx.Rebind(query)
	var matches []dbDefinitionMapping
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, nil
	}
	var out []DefinitionMapping
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
