package uservocabulary

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	upsertVocabularyEntryQuery = `INSERT INTO
        user_vocabulary_entries (
            user_id, language_code, vocabulary_id, vocabulary_type, vocabulary_display, study_note, is_active, is_visible, unique_hash
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        ) ON CONFLICT (user_id, language_code, unique_hash) DO UPDATE
        SET
            is_active = $7,
            is_visible = $8
        RETURNING _id`
)

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
