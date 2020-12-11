package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "user",
		Routes: []router.Route{
			{
				Path:    "unsubscribe_user_1",
				Handler: handleUnsubscribeUser,
			}, {
				Path:    "get_user_preferences_for_token_1",
				Handler: handleGetUserPreferencesForToken,
			},
		},
	})
}

type unsubscribeUserRequest struct {
	Token        string `json:"token"`
	EmailAddress string `json:"email_address"`
}

type unsubscribeUserResponse struct {
	Success bool `json:"success"`
}

func handleUnsubscribeUser(body []byte) (interface{}, error) {
	var r unsubscribeUserRequest
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := encrypt.WithDecodedToken(r.Token, func(t encrypt.TokenPair) error {
		if t.Key != "unsubscribe" {
			return fmt.Errorf("incorrect key")
		}
		userID, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			didUpdate, err = users.UnsubscribeUserForIDAndEmail(tx, users.UserID(userID), r.EmailAddress)
			return err
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return unsubscribeUserResponse{
		Success: didUpdate,
	}, nil
}

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
	var readingLevelsByLanguageCode map[wordsmith.LanguageCode]userreadability.ReadingLevelClassification
	if err := encrypt.WithDecodedToken(req.Token, func(t encrypt.TokenPair) error {
		if t.Key != routes.SubscriptionManagementRouteEncryptionKey.Str() {
			return fmt.Errorf("Invalid key")
		}
		userID, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("Invalid type")
		}
		return database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			readingLevelsByLanguageCode, err = userreadability.GetReadingLevelClassificationsForUser(tx, users.UserID(userID))
			return err
		})
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
