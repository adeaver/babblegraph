package user

import (
	"babblegraph/model/content"
	"babblegraph/model/routes"
	"babblegraph/model/usercontenttopics"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

type getUserContentTopicsForTokenRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type getUserContentTopicsForTokenResponse struct {
	TopicIDs []content.TopicID  `json:"topics"`
	Error    *clienterror.Error `json:"error"`
}

func handleGetUserContentTopicsForToken(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getUserContentTopicsForTokenRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, clientErr, err := routermiddleware.ValidateUserAuthWithToken(userAuth, routermiddleware.ValidateUserAuthWithTokenInput{
		Token:   req.SubscriptionManagementToken,
		KeyType: routes.SubscriptionManagementRouteEncryptionKey,
	})
	switch {
	case err != nil:
		return nil, err
	case clientErr != nil:
		return getUserContentTopicsForTokenResponse{
			Error: clientErr,
		}, nil
	default:
		// no-op
	}
	var topicIDs []content.TopicID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topicIDs, err = usercontenttopics.GetTopicIDsForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserContentTopicsForTokenResponse{
		TopicIDs: topicIDs,
	}, nil
}

type updateUserContentTopicsForTokenRequest struct {
	SubscriptionManagementToken string            `json:"subscription_management_token"`
	EmailAddress                *string           `json:"email_address,omitempty"`
	ActiveTopicIDs              []content.TopicID `json:"active_topic_ids"`
}

type updateUserContentTopicsForTokenResponse struct {
	Error   *clienterror.Error `json:"error,omitempty"`
	Success bool               `json:"success"`
}

func handleUpdateUserContentTopicsForToken(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req updateUserContentTopicsForTokenRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, clientErr, err := routermiddleware.ValidateUserAuthWithToken(userAuth, routermiddleware.ValidateUserAuthWithTokenInput{
		EmailAddress:        req.EmailAddress,
		RequireEmailAddress: true,
		Token:               req.SubscriptionManagementToken,
		KeyType:             routes.SubscriptionManagementRouteEncryptionKey,
	})
	switch {
	case err != nil:
		return nil, err
	case clientErr != nil:
		return updateUserContentTopicsForTokenResponse{
			Error: clientErr,
		}, nil
	default:
		// no-op
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var topicMappings []usercontenttopics.ContentTopicWithTopicID
		for _, t := range req.ActiveTopicIDs {
			topic, err := content.GetContentTopicForTopicID(tx, t)
			if err != nil {
				return err
			}
			topicMappings = append(topicMappings, usercontenttopics.ContentTopicWithTopicID{
				Topic:   *topic,
				TopicID: t,
			})
		}
		return usercontenttopics.UpdateContentTopicsForUser(tx, *userID, topicMappings)
	}); err != nil {
		return nil, err
	}
	return updateUserContentTopicsForTokenResponse{
		Success: true,
	}, nil
}
