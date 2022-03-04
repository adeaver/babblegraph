package usernewsletterpreferences

import (
	"babblegraph/model/content"
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"time"
)

// TODO: move schedule into here maybe?

type UserNewsletterPreferences struct {
	UserID                                   users.UserID
	LanguageCode                             wordsmith.LanguageCode
	ShouldIncludeLemmaReinforcementSpotlight bool
	PodcastPreferences                       PodcastPreferences
}

type PodcastPreferences struct {
	ArePodcastsEnabled         bool
	IncludeExplicitPodcasts    bool
	MinimumDurationNanoseconds *time.Duration
	MaximumDurationNanoseconds *time.Duration
	ExcludedSourceIDs          []content.SourceID
}

type userLemmaReinforcementSpotlightPreferencesID string

type dbUserLemmaReinforcementSpotlightPreferences struct {
	ID                                       userLemmaReinforcementSpotlightPreferencesID `db:"_id"`
	LanguageCode                             wordsmith.LanguageCode                       `db:"language_code"`
	UserID                                   users.UserID                                 `db:"user_id"`
	ShouldIncludeLemmaReinforcementSpotlight bool                                         `db:"should_include_lemma_reinforcement_spotlight"`
}

type userPodcastPreferecesID string

type dbUserPodcastPreferences struct {
	CreatedAt                  time.Time               `db:"created_at"`
	LastModifiedAt             time.Time               `db:"last_modified_at"`
	ID                         userPodcastPreferecesID `db:"_id"`
	LanguageCode               wordsmith.LanguageCode  `db:"language_code"`
	UserID                     users.UserID            `db:"user_id"`
	ArePodcastsEnabled         bool                    `db:"podcasts_enabled"`
	IncludeExplicitPodcasts    bool                    `db:"include_explicit_podcasts"`
	MinimumDurationNanoseconds *time.Duration          `db:"minimum_duration_nanoseconds"`
	MaximumDurationNanoseconds *time.Duration          `db:"maximum_duration_nanoseconds"`
}

type dbUserPodcastSourceMappingID string

type dbUserPodcastSourcePreferences struct {
	CreatedAt      time.Time                    `db:"created_at"`
	LastModifiedAt time.Time                    `db:"last_modified_at"`
	ID             dbUserPodcastSourceMappingID `db:"_id"`
	LanguageCode   wordsmith.LanguageCode       `db:"language_code"`
	UserID         users.UserID                 `db:"user_id"`
	SourceID       content.SourceID             `db:"source_id"`
	IsActive       bool                         `db:"is_active"`
}
