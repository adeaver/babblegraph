package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

// TODO: move schedule into here maybe?

type UserNewsletterPreferences struct {
	UserID                          users.UserID
	LanguageCode                    wordsmith.LanguageCode
	ShouldIncludeLemmaReinforcement bool
}

type userLemmaReinforcementPreferencesID string

type dbUserLemmaReinforcementPreferences struct {
	ID                              userLemmaReinforcementPreferencesID `db:"_id"`
	LanguageCode                    wordsmith.LanguageCode              `db:"language_code"`
	UserID                          users.UserID                        `db:"user_id"`
	ShouldIncludeLemmaReinforcement bool                                `db:"should_include_lemma_reinforcement"`
}
