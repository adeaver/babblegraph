package wordsmith

import (
	"github.com/jmoiron/sqlx"
)

type PhraseDefinitionID string

type dbPhraseDefinition struct {
	ID         PhraseDefinitionID `db:"_id"`
	Language   LanguageCode       `db:"language"`
	CorpusID   CorpusID           `db:"corpus_id"`
	Phrase     string             `db:"phrase"`
	Definition string             `db:"definition"`
}

func (d dbPhraseDefinition) ToNonDB() PhraseDefinition {
	return PhraseDefinition{
		ID:           d.ID,
		LanguageCode: d.Language,
		Phrase:       d.Phrase,
		Definition:   d.Definition,
	}
}

type PhraseDefinition struct {
	ID           PhraseDefinitionID `json:"id"`
	LanguageCode LanguageCode       `json:"language"`
	Phrase       string             `json:"phrase"`
	Definition   string             `json:"defintion"`
}

type dbLemmaPhraseDefinitionMapping struct {
	ID                 string             `db:"_id"`
	Language           LanguageCode       `db:"language"`
	CorpusID           CorpusID           `db:"corpus_id"`
	LemmaPhrase        string             `db:"lemma_phrase"`
	PhraseDefinitionID PhraseDefinitionID `db:"phrase_definition_id"`
}

const getPhraseDefinitionsForLemmaPhrasesQuery = "SELECT * FROM phrase_definitions WHERE _id IN (SELECT phrase_definition_id FROM lemma_phrase_definition_mappings WHERE corpus_id = $1 AND lemma_phrase = $2)"

func GetPhraseDefinitionsForLemmaPhrases(tx *sqlx.Tx, corpusID CorpusID, lemmaPhrases []string) ([]PhraseDefinition, error) {
	var out []PhraseDefinition
	phraseDefinitionIDs := make(map[PhraseDefinitionID]bool)
	for _, p := range lemmaPhrases {
		var matches []dbPhraseDefinition
		if err := tx.Select(&matches, getPhraseDefinitionsForLemmaPhrasesQuery, corpusID, p); err != nil {
			return nil, err
		}
		for _, match := range matches {
			if _, ok := phraseDefinitionIDs[match.ID]; ok {
				continue
			}
			phraseDefinitionIDs[match.ID] = true
			out = append(out, match.ToNonDB())
		}
	}
	return out, nil
}

const getPhraseDefintionsQuery = "SELECT * FROM phrase_definitions WHERE  _id IN (?)"

func GetPhraseDefintions(tx *sqlx.Tx, definitionIDs []PhraseDefinitionID) ([]PhraseDefinition, error) {
	query, args, err := sqlx.In(getPhraseDefintionsQuery, definitionIDs)
	if err != nil {
		return nil, nil
	}
	sql := tx.Rebind(query)
	var matches []dbPhraseDefinition
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, nil
	}
	var out []PhraseDefinition
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}
