package uservocabulary

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
	"time"
)

type UserVocabularyEntryID string

type dbUserVocabularyEntry struct {
	CreatedAt         time.Time              `db:"created_at"`
	LastModifiedAt    time.Time              `db:"last_modified_at"`
	ID                UserVocabularyEntryID  `db:"_id"`
	UserID            users.UserID           `db:"user_id"`
	LanguageCode      wordsmith.LanguageCode `db:"language_code"`
	VocabularyID      *string                `db:"vocabulary_id"`
	VocabularyType    VocabularyType         `db:"vocabulary_type"`
	VocabularyDisplay string                 `db:"vocabulary_display"`
	StudyNote         *string                `db:"study_note"`
	IsActive          bool                   `db:"is_active"`
	IsVisible         bool                   `db:"is_visible"`

	// Phrases are unique by display text or definition ID
	// Lemmas are unique by lemma ID
	UniqueHash UniqueHash `db:"unique_hash"`
}

type UniqueHash string

func (u UniqueHash) Str() string {
	return string(u)
}

func (d dbUserVocabularyEntry) ToNonDB() UserVocabularyEntry {
	return UserVocabularyEntry{
		ID:                d.ID,
		VocabularyID:      d.VocabularyID,
		VocabularyType:    d.VocabularyType,
		VocabularyDisplay: d.VocabularyDisplay,
		StudyNote:         d.StudyNote,
		IsActive:          d.IsActive,
		IsVisible:         d.IsVisible,
		UniqueHash:        d.UniqueHash,
	}
}

type VocabularyType string

const (
	VocabularyTypeLemma  VocabularyType = "lemma"
	VocabularyTypePhrase VocabularyType = "phrase"
)

func (v VocabularyType) Ptr() *VocabularyType {
	return &v
}

func (v VocabularyType) Str() string {
	return string(v)
}

func GetVocabularyTypeFromString(s string) (*VocabularyType, error) {
	switch strings.ToLower(s) {
	case VocabularyTypeLemma.Str():
		return VocabularyTypeLemma.Ptr(), nil
	case VocabularyTypePhrase.Str():
		return VocabularyTypePhrase.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unknown vocabulary type %s", s)
	}
}

type UserVocabularyEntry struct {
	ID UserVocabularyEntryID `json:"id"`
	// Phrases that have no definition will have no ID here.
	VocabularyID      *string        `json:"vocabulary_id,omitempty"`
	VocabularyType    VocabularyType `json:"vocabulary_type"`
	VocabularyDisplay string         `json:"vocabulary_display"`
	Definition        *string        `json:"definition,omitempty"`
	StudyNote         *string        `json:"study_note,omitempty"`
	IsActive          bool           `json:"is_active"`
	IsVisible         bool           `json:"is_visible"`
	UniqueHash        UniqueHash     `json:"unique_hash"`
}

func (u UserVocabularyEntry) AsLemmaIDPhrases() ([][]wordsmith.LemmaID, error) {
	switch u.VocabularyType {
	case VocabularyTypeLemma:
		return [][]wordsmith.LemmaID{{wordsmith.LemmaID(*u.VocabularyID)}}, nil
	case VocabularyTypePhrase:
		return GetLemmaIDPhrasesForPhrase(u.VocabularyDisplay)
	default:
		return nil, fmt.Errorf("Unrecognized vocabulary type: %s", u.VocabularyType)
	}
}
