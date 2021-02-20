package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

type MappingID string

type Mapping struct {
	ID           MappingID              `json:"id"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	UserID       users.UserID           `json:"user_id"`
	LemmaID      wordsmith.LemmaID      `json:"lemma_id"`
	IsVisible    bool                   `json:"is_visible"`
	IsActive     bool                   `json:"is_active"`
}

type dbMapping struct {
	ID           MappingID              `db:"_id"`
	LanguageCode wordsmith.LanguageCode `db:"language_code"`
	UserID       users.UserID           `db:"user_id"`
	LemmaID      wordsmith.LemmaID      `db:"lemma_id"`
	IsVisible    bool                   `db:"is_visible"`
	IsActive     bool                   `db:"is_active"`
}

func (d dbMapping) ToNonDB() Mapping {
	return Mapping{
		ID:           d.ID,
		LanguageCode: d.LanguageCode,
		UserID:       d.UserID,
		LemmaID:      d.LemmaID,
		IsVisible:    d.IsVisible,
		IsActive:     d.IsActive,
	}
}
