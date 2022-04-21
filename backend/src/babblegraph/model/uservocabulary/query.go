package uservocabulary

import (
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	selectVocabularyEntryForUserQuery = "SELECT * FROM user_vocabulary_entries WHERE is_visible = TRUE AND user_id = $1 AND language_code = $2"
	upsertVocabularyEntryQuery        = `INSERT INTO
        user_vocabulary_entries (
            user_id, language_code, vocabulary_id, vocabulary_type, vocabulary_display, study_note, is_active, is_visible, unique_hash
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        ) ON CONFLICT (user_id, language_code, unique_hash) DO UPDATE
        SET
            is_active = $7,
            is_visible = $8
        RETURNING _id`

	selectVocabularySpotlightRecordForUserQuery = "SELECT * FROM user_vocabulary_spotlight_records WHERE user_id = $1 AND language_code = $2 ORDER BY last_sent_on ASC"
	upsertVocabularySpotlightRecordQuery        = `INSERT INTO
        user_vocabulary_spotlight_records (
            user_id, language_code, vocabulary_entry_id, last_sent_on, number_of_times_sent
        ) VALUES (
            $1, $2, $3, timezone('utc', now()), 1
        ) ON CONFLICT (user_id, language_code, vocabulary_entry_id) DO UPDATE
        SET
        last_sent_on=timezone('utc', now()),
        number_of_times_sent=user_vocabulary_spotlight_records.number_of_times_sent+1`

	// TODO(migration): remove this query
	createVocabularySpotlightRecordQuery = `INSERT INTO
        user_vocabulary_spotlight_records (
            user_id, language_code, vocabulary_entry_id, last_sent_on, number_of_times_sent
        ) VALUES (
            $1, $2, $3, $4, $5
        ) ON CONFLICT (user_id, language_code, vocabulary_entry_id) DO UPDATE
        SET
        last_sent_on=$4,
        number_of_times_sent=$5`
)

func GetUserVocabularyEntries(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, includeDefinitions bool) ([]UserVocabularyEntry, error) {
	var matches []dbUserVocabularyEntry
	if err := tx.Select(&matches, selectVocabularyEntryForUserQuery, userID, languageCode); err != nil {
		return nil, err
	}
	var phraseDefinitionIDs []wordsmith.PhraseDefinitionID
	var lemmaIDs []wordsmith.LemmaID
	var out []UserVocabularyEntry
	for _, m := range matches {
		out = append(out, m.ToNonDB())
		switch {
		case m.VocabularyID == nil:
			// no-op
		case m.VocabularyType == VocabularyTypePhrase:
			phraseDefinitionIDs = append(phraseDefinitionIDs, wordsmith.PhraseDefinitionID(*m.VocabularyID))
		case m.VocabularyType == VocabularyTypeLemma:
			lemmaIDs = append(lemmaIDs, wordsmith.LemmaID(*m.VocabularyID))
		}
	}
	if includeDefinitions {
		phraseDefinitions := make(map[wordsmith.PhraseDefinitionID]wordsmith.PhraseDefinition)
		lemmaDefinitions := make(map[wordsmith.LemmaID][]wordsmith.DefinitionMapping)
		if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
			pDefinitions, err := wordsmith.GetPhraseDefintions(tx, phraseDefinitionIDs)
			if err != nil {
				return err
			}
			for _, d := range pDefinitions {
				phraseDefinitions[d.ID] = d
			}
			lDefinitions, err := wordsmith.GetDefinitionMappingsForLemmaIDs(tx, wordsmith.SpanishOpenDefinitions, lemmaIDs)
			if err != nil {
				return err
			}
			for _, d := range lDefinitions {
				lemmaDefinitions[d.LemmaID] = append(lemmaDefinitions[d.LemmaID], d)
			}
			return nil
		}); err != nil {
			return nil, err
		}
		for idx, e := range out {
			switch {
			case e.VocabularyID == nil:
				// no-op
			case e.VocabularyType == VocabularyTypePhrase:
				if definition, ok := phraseDefinitions[wordsmith.PhraseDefinitionID(*e.VocabularyID)]; ok {
					out[idx].Definition = ptr.String(definition.Definition)
				}
			case e.VocabularyType == VocabularyTypeLemma:
				if definitions, ok := lemmaDefinitions[wordsmith.LemmaID(*e.VocabularyID)]; ok {
					var definition []string
					for _, d := range definitions {
						definition = append(definition, strings.ToLower(d.EnglishDefinition))
					}
					out[idx].Definition = ptr.String(strings.Join(definition, "; "))
				}
			}
		}
	}
	return out, nil
}

type UpsertVocabularyEntryInput struct {
	UserID       users.UserID
	LanguageCode wordsmith.LanguageCode
	Hashable     uniqueHashable
	StudyNote    *string
	IsActive     bool
	IsVisible    bool
}

func UpsertVocabularyEntry(tx *sqlx.Tx, input UpsertVocabularyEntryInput) (*UserVocabularyEntryID, error) {
	vocabularyID := input.Hashable.getID()
	vocabularyType := input.Hashable.getType()
	vocabularyDisplay := input.Hashable.getDisplay()
	hash := GetUniqueHash(input.Hashable)
	isActive := input.IsVisible && input.IsActive
	rows, err := tx.Query(upsertVocabularyEntryQuery, input.UserID, input.LanguageCode, vocabularyID, vocabularyType, vocabularyDisplay, input.StudyNote, isActive, input.IsVisible, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userVocabularyID UserVocabularyEntryID
	for rows.Next() {
		if err := rows.Scan(&userVocabularyID); err != nil {
			return nil, err
		}
	}
	return &userVocabularyID, nil
}

func GetUserVocabularySpotlightRecords(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) ([]UserVocabularySpotlightRecord, error) {
	var matches []dbUserVocabularySpotlightRecord
	if err := tx.Select(&matches, selectVocabularySpotlightRecordForUserQuery, userID, languageCode); err != nil {
		return nil, err
	}
	var out []UserVocabularySpotlightRecord
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func UpsertUserVocabularySpotlightRecord(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, id UserVocabularyEntryID) error {
	if _, err := tx.Exec(upsertVocabularySpotlightRecordQuery, userID, languageCode, id); err != nil {
		return err
	}
	return nil
}
