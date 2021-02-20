package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

type MappingID string

type Mapping struct {
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	LemmaID      wordsmith.LemmaID      `json:"lemma_id"`
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
		LanguageCode: d.LanguageCode,
		LemmaID:      d.LemmaID,
		IsActive:     d.IsActive,
	}
}
