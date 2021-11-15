package user

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/usercontenttopics"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

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

type updateUserContentTopicsForTokenResponse struct{}

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
	return updateUserContentTopicsForTokenResponse{}, nil
}
