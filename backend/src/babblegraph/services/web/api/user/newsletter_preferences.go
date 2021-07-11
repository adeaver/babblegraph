package user

import (
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type getUserNewsletterPreferencesRequest struct {
	LanguageCode                string `json:"language_code"`
	EmailAddress                string `json:"email_address"`
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type getUserNewsletterPreferencesResponse struct {
	LanguageCode wordsmith.LanguageCode    `json:"language_code"`
	Preferences  userNewsletterPreferences `json:"preferences"`
}

type userNewsletterPreferences struct {
	IsLemmaReinforcementSpotlightActive bool `json:"is_lemma_reinforcement_spotlight_active"`
}

func getUserNewsletterPreferences(userID users.UserID, body []byte) (interface{}, error) {
	var req getUserNewsletterPreferencesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	tokenUserID, err := parseSubscriptionManagementToken(req.SubscriptionManagementToken, ptr.String(req.EmailAddress))
	switch {
	case err != nil:
		return nil, err
	case *tokenUserID != userID:
		return nil, fmt.Errorf("invalid token")
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	var prefs *usernewsletterpreferences.UserNewsletterPreferences
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		prefs, err = usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(tx, userID, *languageCode)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserNewsletterPreferencesResponse{
		LanguageCode: *languageCode,
		Preferences: userNewsletterPreferences{
			IsLemmaReinforcementSpotlightActive: prefs.ShouldIncludeLemmaReinforcementSpotlight,
		},
	}, nil
}
