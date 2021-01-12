package user

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/routes"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/ptr"
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
			}, {
				Path:    "update_user_preferences_for_token_1",
				Handler: handleUpdateUserPreferencesForToken,
			}, {
				Path:    "get_user_content_topics_for_token_1",
				Handler: handleGetUserContentTopicsForToken,
			}, {
				Path:    "update_user_content_topics_for_token_1",
				Handler: handleUpdateUserContentTopicsForToken,
			},
		},
	})
}

func parseSubscriptionManagementToken(token string, emailAddress *string) (*users.UserID, error) {
	var userID users.UserID
	if err := encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != routes.SubscriptionManagementRouteEncryptionKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		userIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		userID = users.UserID(userIDStr)
		return nil
	}); err != nil {
		return nil, err
	}
	if emailAddress != nil {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			user, err := users.LookupUserForIDAndEmail(tx, *userID, req.EmailAddress)
			if err != nil {
				return err
			}
			if user == nil {
				return fmt.Errorf("Invalid email address for token")
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return &userID, nil
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
	userID, err := parseSubscriptionManagementToken(r.Token, ptr.String(r.EmailAddress))
	if err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		didUpdate, err = users.UnsubscribeUserForIDAndEmail(tx, *userID, r.EmailAddress)
		return err
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

type getUserContentTopicsForTokenRequest struct {
	Token string `json:"token"`
}

type getUserContentTopicsForTokenResponse struct {
	ContentTopics []contenttopics.ContentTopic `json:"content_topics"`
}

func handleGetUserContentTopicsForToken(body []byte) (interface{}, error) {
	var req getUserContentTopicsForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := parseSubscriptionManagementToken(req.Token, nil)
	if err != nil {
		return nil, err
	}
	var contentTopics []contenttopics.ContentTopic
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		contentTopics, err = usercontenttopics.GetContentTopicsForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserContentTopicsForTokenResponse{
		ContentTopics: contentTopics,
	}, nil
}

type updateUserContentTopicsForTokenRequest struct {
	Token         string                       `json:"token"`
	EmailAddress  string                       `json:"email_address"`
	ContentTopics []contenttopics.ContentTopic `json:"content_topics"`
}

func handleUpdateUserContentTopicsForToken(body []byte) (interface{}, error) {
	var req updateUserContentTopicsForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := parseSubscriptionManagementToken(req.Token, ptr.String(req.EmailAddress))
	if err != nil {
		return nil, err
	}
    if err := database.WithTx(func(tx *sqlx.Tx) error {
        return usercontenttopics.UpdateContentTopicsForUser(tx, *userID, req.ContentTopics)
    }); err != nil {
        return nil, err
    }
	return struct{}, nil
}
