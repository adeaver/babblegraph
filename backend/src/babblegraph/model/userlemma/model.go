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

type UserLemmaReinforcementSpotlightRecord struct {
	LanguageCode      wordsmith.LanguageCode
	UserID            users.UserID
	LemmaID           wordsmith.LemmaID
	LastSentOn        time.Time
	NumberOfTimesSent int64
}

type dbUserLemmaReinforcementSpotlightRecord struct {
	ID                userLemmaReinforcementReminderID `db:"_id"`
	LanguageCode      wordsmith.LanguageCode           `db:"language_code"`
	UserID            users.UserID                     `db:"user_id"`
	LemmaID           wordsmith.LemmaID                `db:"lemma_id"`
	LastSentOn        time.Time                        `db:"last_sent_on"`
	NumberOfTimesSent int64                            `db:"number_of_times_sent"`
}

func (d dbUserLemmaReinforcementSpotlightRecord) ToNonDB() UserLemmaReinforcementSpotlightRecord {
	return UserLemmaReinforcementSpotlightRecord{
		LanguageCode:      d.LanguageCode,
		UserID:            d.UserID,
		LemmaID:           d.LemmaID,
		LastSentOn:        d.LastSentOn,
		NumberOfTimesSent: d.NumberOfTimesSent,
	}
}

type PhraseMappingID string

type dbUserPhraseMapping struct {
	ID             PhraseMappingID `db:"_id"`
	CreatedAt      time.Time       `db:"created_at"`
	LastModifiedAt time.Time       `db:"last_modified_at"`
	UserID         users.UserID    `db:"user_id"`
	IsActive       bool            `db:"is_active"`
}

type dbUserPhraseOption struct {
	ID             PhraseMappingID `db:"_id"`
	CreatedAt      time.Time       `db:"created_at"`
	LastModifiedAt time.Time       `db:"last_modified_at"`
}

type UserPhraseMapping struct {
	ID         PhraseMappingID `db:"_id"`
	UserID     users.UserID    `db:"user_id"`
	IsActive   bool            `db:"is_active"`
	PhraseText string          `db:"phrase_text"`
}
