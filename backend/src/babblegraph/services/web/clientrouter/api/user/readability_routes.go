package user

import (
	"babblegraph/model/userreadability"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type getUserPreferencesForTokenRequest struct {
	Token string `json:"token"`
}

type getUserPreferencesForTokenResponse struct {
	ClassificationsByLanguage []readingLevelClassificationForLanguage `json:"classifications_by_language"`
}

type readingLevelClassificationForLanguage struct {
	LanguageCode               wordsmith.LanguageCode                     `json:"language_code"`
	ReadingLevelClassification userreadability.ReadingLevelClassification `json:"reading_level_classification"`
}

func handleGetUserPreferencesForToken(body []byte) (interface{}, error) {
	var req getUserPreferencesForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := parseSubscriptionManagementToken(req.Token, nil)
	if err != nil {
		return nil, err
	}
	var readingLevelsByLanguageCode map[wordsmith.LanguageCode]userreadability.ReadingLevelClassification
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		readingLevelsByLanguageCode, err = userreadability.GetReadingLevelClassificationsForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	var classifications []readingLevelClassificationForLanguage
	for languageCode, classification := range readingLevelsByLanguageCode {
		classifications = append(classifications, readingLevelClassificationForLanguage{
			LanguageCode:               languageCode,
			ReadingLevelClassification: classification,
		})
	}
	return getUserPreferencesForTokenResponse{
		ClassificationsByLanguage: classifications,
	}, nil
}

type updateUserPreferencesForTokenRequest struct {
	Token                     string                                  `json:"token"`
	EmailAddress              string                                  `json:"email_address"`
	ClassificationsByLanguage []readingLevelClassificationForLanguage `json:"classifications_by_language"`
}

type updateUserPreferencesForTokenResponse struct {
	DidUpdate bool `json:"did_update"`
}

func handleUpdateUserPreferencesForToken(body []byte) (interface{}, error) {
	var req updateUserPreferencesForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := parseSubscriptionManagementToken(req.Token, ptr.String(req.EmailAddress))
	if err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, classification := range req.ClassificationsByLanguage {
			update, err := userreadability.UpdateReadingLevelClassificationForUser(tx, *userID, classification.LanguageCode, classification.ReadingLevelClassification)
			if err != nil {
				return err
			}
			didUpdate = didUpdate || update
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return updateUserPreferencesForTokenResponse{
		DidUpdate: didUpdate,
	}, nil
}
