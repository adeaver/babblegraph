package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getLemmaReinforcementPreferencesQuery    = "SELECT * FROM user_lemma_reinforcement_preferences WHERE user_id = $1 AND language_code = $2"
	updateLemmaReinforcementPreferencesQuery = "INSERT INTO user_lemma_reinforcement_preferences (user_id, language_code, should_include_lemma_reinforcement) VALUES ($1, $2, $3) ON CONFLICT (user_id, language_code) SET should_include_lemma_reinforcement= $3"
)

func GetUserNewsletterPrefrencesForLanguage(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*UserNewsletterPreferences, error) {
	shouldIncludeLemmaReinforcement := true
	lemmaReinforcementPreferences, err := lookupLemmaReinforcementPreferences(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	if lemmaReinforcementPreferences != nil {
		shouldIncludeLemmaReinforcement = lemmaReinforcementPreferences.ShouldIncludeLemmaReinforcement
	}
	return &UserNewsletterPreferences{
		UserID:                          userID,
		LanguageCode:                    languageCode,
		ShouldIncludeLemmaReinforcement: shouldIncludeLemmaReinforcement,
	}, nil
}

func lookupLemmaReinforcementPreferences(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*dbUserLemmaReinforcementPreferences, error) {
	var matches []dbUserLemmaReinforcementPreferences
	if err := tx.Select(&matches, getLemmaReinforcementPreferencesQuery, userID, languageCode); err != nil {
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

func UpdateIsUserLemmaReinforcementActiveForLanguage(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, isActive bool) error {
	if _, err := tx.Exec(updateLemmaReinforcementPreferencesQuery, userID, languageCode, isActive); err != nil {
		return err
	}
	return nil
}
