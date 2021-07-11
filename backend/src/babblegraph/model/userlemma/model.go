package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"time"
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

type userLemmaReinforcementReminderID string

type UserLemmaReinforcementRecord struct {
	LanguageCode      wordsmith.LanguageCode
	UserID            users.UserID
	LemmaID           wordsmith.LemmaID
	LastSentOn        time.Time
	NumberOfTimesSent int64
}

type dbUserLemmaReinforcementRecord struct {
	ID                userLemmaReinforcementReminderID `db:"_id"`
	LanguageCode      wordsmith.LanguageCode           `db:"language_code"`
	UserID            users.UserID                     `db:"user_id"`
	LemmaID           wordsmith.LemmaID                `db:"lemma_id"`
	LastSentOn        time.Time                        `db:"last_sent_on"`
	NumberOfTimesSent int64                            `db:"number_of_times_sent"`
}

func (d dbUserLemmaReinforcementRecord) ToNonDB() UserLemmaReinforcementRecord {
	return UserLemmaReinforcementRecord{
		LanguageCode:      d.LanguageCode,
		UserID:            d.UserID,
		LemmaID:           d.LemmaID,
		LastSentOn:        d.LastSentOn,
		NumberOfTimesSent: d.NumberOfTimesSent,
	}
}
