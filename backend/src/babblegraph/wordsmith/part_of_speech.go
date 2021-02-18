package wordsmith

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PartOfSpeechID string

type PartOfSpeechCode string

// TODO: make this more robust
// use: https://freeling-user-manual.readthedocs.io/en/latest/tagsets/tagset-es/
func (p PartOfSpeechCode) ToDisplayCategory() string {
	category := strings.ToLower(string(p))[0]
	switch category {
	case 97:
		return "Adjective"
	case 99:
		return "Conjunction"
	case 100:
		return "Determiner"
	case 110:
		return "Noun"
	case 112:
		return "Pronoun"
	case 114:
		return "Adverb"
	case 115:
		return "Preposition"
	case 118:
		return "Verb"
	case 105:
		return "Interjection"
	default:
		log.Println(fmt.Sprintf("Unknown category: %s", string(p)))
		return ""
	}
}

type PartOfSpeech struct {
	ID       PartOfSpeechID
	Language LanguageCode
	CorpusID CorpusID
	Code     PartOfSpeechCode
}

type dbPartOfSpeech struct {
	ID       PartOfSpeechID   `db:"_id"`
	Language LanguageCode     `db:"language"`
	CorpusID CorpusID         `db:"corpus_id"`
	Code     PartOfSpeechCode `db:"code"`
}

func (d dbPartOfSpeech) ToNonDB() PartOfSpeech {
	return PartOfSpeech{
		ID:       d.ID,
		Language: d.Language,
		CorpusID: d.CorpusID,
		Code:     d.Code,
	}
}

const partsOfSpeechByIDQuery = "SELECT * FROM parts_of_speech WHERE corpus_id = '%s' AND _id IN (?)"

func GetPartOfSpeechByIDs(tx *sqlx.Tx, corpusID CorpusID, ids []PartOfSpeechID) ([]PartOfSpeech, error) {
	query, args, err := sqlx.In(fmt.Sprintf(partsOfSpeechByIDQuery, corpusID), ids)
	if err != nil {
		return nil, nil
	}
	sql := tx.Rebind(query)
	var matches []dbPartOfSpeech
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, nil
	}
	var out []PartOfSpeech
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
