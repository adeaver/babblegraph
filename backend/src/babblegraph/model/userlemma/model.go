package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

type MappingID string

type Mapping struct {
	ID           MappingID
	LanguageCode wordsmith.LanguageCode
	UserID       users.UserID
	LemmaID      wordsmith.LemmaID
	IsVisible    bool
	IsActive     bool
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
