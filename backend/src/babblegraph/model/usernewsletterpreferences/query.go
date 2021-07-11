package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getLemmaReinforcementSpotlightPreferencesQuery    = "SELECT * FROM user_lemma_reinforcement_spotlight_preferences WHERE user_id = $1 AND language_code = $2"
	updateLemmaReinforcementSpotlightPreferencesQuery = "INSERT INTO user_lemma_reinforcement_spotlight_preferences (user_id, language_code, should_include_lemma_reinforcement_spotlight) VALUES ($1, $2, $3) ON CONFLICT (user_id, language_code) SET should_include_lemma_reinforcement_spotlight = $3"
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
}

func UpdateUserNewsletterPreferences(tx *sqlx.Tx, input UpdateUserNewsletterPreferencesInput) error {
	return updateLemmaReinforcementSpotlightPreferences(tx, input.UserID, input.LanguageCode, input.IsLemmaReinforcementSpotlightActive)
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
