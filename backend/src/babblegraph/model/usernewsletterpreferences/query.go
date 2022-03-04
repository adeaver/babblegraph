package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getLemmaReinforcementSpotlightPreferencesQuery    = "SELECT * FROM user_lemma_reinforcement_spotlight_preferences WHERE user_id = $1 AND language_code = $2"
	updateLemmaReinforcementSpotlightPreferencesQuery = `INSERT INTO
        user_lemma_reinforcement_spotlight_preferences (user_id, language_code, should_include_lemma_reinforcement_spotlight)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id, language_code)
    DO UPDATE SET
        should_include_lemma_reinforcement_spotlight = $3`

	upsertPodcastSourcePreferencesQuery = `INSERT INTO user_podcast_source_preferences
        (user_id, language_code, source_id, is_active)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, language_code, source_id)
    DO UPDATE SET
        is_active = $4`

	upsertUserPodcastPreferencesQuery = `INSERT INTO user_podcast_preferences
        (user_id, language_code, include_explicit_podcasts, minimum_duration_nanoseconds, maximum_duration_nanoseconds)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (user_id, language_code)
    DO UPDATE SET
        include_explicit_podcasts=$3,
        minimum_duration_nanoseconds=$4,
        maximum_duration_nanoseconds=$5`
)

func GetUserNewsletterPrefrencesForLanguage(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*UserNewsletterPreferences, error) {
	shouldIncludeLemmaReinforcementSpotlight := true
	lemmaReinforcementSpotlightPreferences, err := lookupLemmaReinforcementSpotlightPreferences(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	if lemmaReinforcementSpotlightPreferences != nil {
		shouldIncludeLemmaReinforcementSpotlight = lemmaReinforcementSpotlightPreferences.ShouldIncludeLemmaReinforcementSpotlight
	}
	return &UserNewsletterPreferences{
		UserID:                                   userID,
		LanguageCode:                             languageCode,
		ShouldIncludeLemmaReinforcementSpotlight: shouldIncludeLemmaReinforcementSpotlight,
	}, nil
}

type UpdateUserNewsletterPreferencesInput struct {
	UserID                              users.UserID
	LanguageCode                        wordsmith.LanguageCode
	IsLemmaReinforcementSpotlightActive bool
	PodcastPreferences                  *PodcastPreferencesInput
}

// TODO: get PodcastPreferencesInput
type PodcastPreferencesInput struct {
	IncludeExplicitPodcasts    bool
	MinimumDurationNanoseconds *time.Duration
	MaximumDurationNanoseconds *time.Duration
	SourceIDUpdates            []SourceIDUpdate
}

type SourceIDUpdate struct {
	SourceID bool
	IsActive bool
}

func UpdateUserNewsletterPreferences(tx *sqlx.Tx, input UpdateUserNewsletterPreferencesInput) error {
	err := updateLemmaReinforcementSpotlightPreferences(tx, input.UserID, input.LanguageCode, input.IsLemmaReinforcementSpotlightActive)
	if err != nil {
		return err
	}
	if input.PodcastPreferences != nil {
		err = updatePodcastPreferences(tx, input.UserID, input.LanguageCode, *input.PodcastPreferences)
		if err != nil {
			return err
		}
		return updatePodcastSourcePreferences(tx, input.UserID, input.LanguageCode, input.PodcastPreferences.SourceIDUpdates)
	}
	return nil
}

func updateLemmaReinforcementSpotlightPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, isActive bool) error {
	if _, err := tx.Exec(updateLemmaReinforcementSpotlightPreferencesQuery, userID, languageCode, isActive); err != nil {
		return err
	}
	return nil
}

func updatePodcastPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, input PodcastPreferencesInput) error {
	if _, err := tx.Exec(upsertUserPodcastPreferencesQuery, userID, languageCode, input.IncludeExplicitPodcasts, input.MinimumDurationNanoseconds, input.MaximumDurationNanoseconds); err != nil {
		return err
	}
	return nil
}

func updatePodcastSourcePreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, updates []SourceIDUpdate) error {
	for _, u := range updates {
		if _, err := tx.Exec(upsertPodcastSourcePreferencesQuery, userID, languageCode, u.SourceID, u.IsActive); err != nil {
			return err
		}
	}
	return nil
}

func lookupLemmaReinforcementSpotlightPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*dbUserLemmaReinforcementSpotlightPreferences, error) {
	var matches []dbUserLemmaReinforcementSpotlightPreferences
	if err := tx.Select(&matches, getLemmaReinforcementSpotlightPreferencesQuery, userID, languageCode); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		m := matches[0]
		return &m, nil
	default:
		return nil, fmt.Errorf("expected at most one record, but got %d", len(matches))
	}
}
