package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

// TODO: move schedule into here maybe?

type UserNewsletterPreferences struct {
	UserID                                   users.UserID
	LanguageCode                             wordsmith.LanguageCode
	ShouldIncludeLemmaReinforcementSpotlight bool
}

type userLemmaReinforcementSpotlightPreferencesID string

type dbUserLemmaReinforcementSpotlightPreferences struct {
	ID                                       userLemmaReinforcementSpotlightPreferencesID `db:"_id"`
	LanguageCode                             wordsmith.LanguageCode                       `db:"language_code"`
	UserID                                   users.UserID                                 `db:"user_id"`
	ShouldIncludeLemmaReinforcementSpotlight bool                                         `db:"should_include_lemma_reinforcement_spotlight"`
}
