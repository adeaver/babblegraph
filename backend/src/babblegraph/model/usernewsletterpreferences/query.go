package usernewsletterpreferences

import (
	"babblegraph/model/content"
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

	getUserPodcastPreferencesQuery    = "SELECT * FROM user_podcast_preferences WHERE user_id = $1 AND language_code = $2"
	upsertUserPodcastPreferencesQuery = `INSERT INTO user_podcast_preferences
        (user_id, language_code, podcasts_enabled, include_explicit_podcasts, minimum_duration_nanoseconds, maximum_duration_nanoseconds)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (user_id, language_code)
    DO UPDATE SET
        podcasts_enabled=$3,
        include_explicit_podcasts=$4,
        minimum_duration_nanoseconds=$5,
        maximum_duration_nanoseconds=$6,
        last_modified_at=timezone('utc', now())`

	getPodcastSourcePreferencesQuery    = "SELECT * FROM user_podcast_source_preferences WHERE user_id = $1 AND language_code = $2 AND is_active = FALSE"
	upsertPodcastSourcePreferencesQuery = `INSERT INTO user_podcast_source_preferences
        (user_id, language_code, source_id, is_active)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, language_code, source_id)
    DO UPDATE SET
        is_active = $4,
        last_modified_at=timezone('utc', now())`
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
	dbPodcastPreferences, err := lookupPodcastPreferences(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	podcastPreferences := PodcastPreferences{
		ArePodcastsEnabled:      true,
		IncludeExplicitPodcasts: true,
	}
	if dbPodcastPreferences != nil {
		podcastPreferences = PodcastPreferences{
			ArePodcastsEnabled:         dbPodcastPreferences.ArePodcastsEnabled,
			IncludeExplicitPodcasts:    dbPodcastPreferences.IncludeExplicitPodcasts,
			MinimumDurationNanoseconds: dbPodcastPreferences.MinimumDurationNanoseconds,
			MaximumDurationNanoseconds: dbPodcastPreferences.MaximumDurationNanoseconds,
		}
	}
	podcastPreferences.ExcludedSourceIDs, err = lookupInactiveSourceIDs(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	return &UserNewsletterPreferences{
		UserID:                                   userID,
		LanguageCode:                             languageCode,
		ShouldIncludeLemmaReinforcementSpotlight: shouldIncludeLemmaReinforcementSpotlight,
		PodcastPreferences:                       podcastPreferences,
	}, nil
}

type UpdateUserNewsletterPreferencesInput struct {
	UserID                              users.UserID
	LanguageCode                        wordsmith.LanguageCode
	IsLemmaReinforcementSpotlightActive bool
	PodcastPreferences                  *PodcastPreferencesInput
}

type PodcastPreferencesInput struct {
	ArePodcastsEnabled         bool
	IncludeExplicitPodcasts    bool
	MinimumDurationNanoseconds *time.Duration
	MaximumDurationNanoseconds *time.Duration
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
		return updatePodcastPreferences(tx, input.UserID, input.LanguageCode, *input.PodcastPreferences)
	}
	return nil
}

func updateLemmaReinforcementSpotlightPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, isActive bool) error {
	if _, err := tx.Exec(updateLemmaReinforcementSpotlightPreferencesQuery, userID, languageCode, isActive); err != nil {
		return err
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

func updatePodcastPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, input PodcastPreferencesInput) error {
	if _, err := tx.Exec(upsertUserPodcastPreferencesQuery, userID, languageCode, input.ArePodcastsEnabled, input.IncludeExplicitPodcasts, input.MinimumDurationNanoseconds, input.MaximumDurationNanoseconds); err != nil {
		return err
	}
	return nil
}

func lookupPodcastPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*dbUserPodcastPreferences, error) {
	var matches []dbUserPodcastPreferences
	err := tx.Select(&matches, getUserPodcastPreferencesQuery, userID, languageCode)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		return &matches[0], nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one result for user podcast preferences (user id %s, language code %s) but got %d", userID, languageCode, len(matches))
	default:
		panic("unreachable")
	}
}

func updatePodcastSourcePreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, updates []SourceIDUpdate) error {
	for _, u := range updates {
		if _, err := tx.Exec(upsertPodcastSourcePreferencesQuery, userID, languageCode, u.SourceID, u.IsActive); err != nil {
			return err
		}
	}
	return nil
}

func lookupInactiveSourceIDs(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) ([]content.SourceID, error) {
	var matches []dbUserPodcastSourcePreferences
	if err := tx.Select(&matches, getPodcastSourcePreferencesQuery, userID, languageCode); err != nil {
		return nil, err
	}
	var out []content.SourceID
	for _, m := range matches {
		out = append(out, m.SourceID)
	}
	return out, nil
}
