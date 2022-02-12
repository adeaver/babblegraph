package user

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/usercontenttopics"
	"babblegraph/util/database"
	"babblegraph/util/email"
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
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	userID, err := parseSubscriptionManagementToken(req.Token, ptr.String(formattedEmailAddress))
	if err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var topicMappings []usercontenttopics.ContentTopicWithTopicID
		for _, t := range req.ContentTopics {
			topicID, err := content.GetTopicIDByContentTopic(tx, t)
			if err != nil {
				return err
			}
			topicMappings = append(topicMappings, usercontenttopics.ContentTopicWithTopicID{
				Topic:   t,
				TopicID: *topicID,
			})
		}
		return usercontenttopics.UpdateContentTopicsForUser(tx, *userID, topicMappings)
	}); err != nil {
		return nil, err
	}
	return updateUserContentTopicsForTokenResponse{}, nil
}
