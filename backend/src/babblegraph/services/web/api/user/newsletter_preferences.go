package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/services/web/util/routetoken"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/wordsmith"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type getUserNewsletterPreferencesRequest struct {
	EmailAddress                string `json:"email_address"`
	LanguageCode                string `json:"language_code"`
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type getUserNewsletterPreferencesResponse struct {
	LanguageCode wordsmith.LanguageCode    `json:"language_code"`
	Preferences  userNewsletterPreferences `json:"preferences"`
}

type userNewsletterPreferences struct {
	IsLemmaReinforcementSpotlightActive bool `json:"is_lemma_reinforcement_spotlight_active"`
}

func getUserNewsletterPreferences(body []byte) (interface{}, error) {
	var req getUserNewsletterPreferencesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	userID, err := routetoken.ValidateTokenAndEmailAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey, formattedEmailAddress)
	if err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	var prefs *usernewsletterpreferences.UserNewsletterPreferences
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		prefs, err = usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(tx, *userID, *languageCode)
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

type updateUserNewsletterPreferencesRequest struct {
	EmailAddress                string                    `json:"email_address"`
	LanguageCode                string                    `json:"language_code"`
	SubscriptionManagementToken string                    `json:"subscription_management_token"`
	Preferences                 userNewsletterPreferences `json:"preferences"`
}

type updateUserNewsletterPreferencesResponse struct {
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	Success      bool                   `json:"success"`
}

func updateUserNewsletterPreferences(body []byte) (interface{}, error) {
	var req updateUserNewsletterPreferencesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	userID, err := routetoken.ValidateTokenAndEmailAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey, formattedEmailAddress)
	if err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return usernewsletterpreferences.UpdateUserNewsletterPreferences(tx, usernewsletterpreferences.UpdateUserNewsletterPreferencesInput{
			UserID:                              *userID,
			LanguageCode:                        *languageCode,
			IsLemmaReinforcementSpotlightActive: req.Preferences.IsLemmaReinforcementSpotlightActive,
		})
	}); err != nil {
		return nil, err
	}
	return updateUserNewsletterPreferencesResponse{
		LanguageCode: *languageCode,
		Success:      true,
	}, nil
}
