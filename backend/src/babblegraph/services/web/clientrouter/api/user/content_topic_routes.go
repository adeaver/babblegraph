package user

import (
	"babblegraph/model/content"
	"babblegraph/model/routes"
	"babblegraph/model/usercontenttopics"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type getUserContentTopicsRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type getUserContentTopicsResponse struct {
	Error    *clienterror.Error `json:"error,omitempty"`
	TopicIDs []content.TopicID  `json:"topic_ids,omitempty"`
}

func getUserContentTopics(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getUserContentTopicsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return getUserContentTopicsResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
		}, nil
	}
	cErr, err := routermiddleware.ValidateUserAuth(userAuth, routermiddleware.ValidateUserAuthInput{
		DecodedUserID: *userID,
	})
	switch {
	case err != nil:
		return nil, err
	case cErr != nil:
		return getUserContentTopicsResponse{
			Error: cErr,
		}, nil
	}
	var topicIDs []content.TopicID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topicIDs, err = usercontenttopics.GetTopicIDsForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserContentTopicsResponse{
		TopicIDs: topicIDs,
	}, nil
}

type upsertUserContentTopicsRequest struct {
	EmailAddress                *string           `json:"email_address,omitempty"`
	SubscriptionManagementToken string            `json:"subscription_management_token"`
	TopicIDs                    []content.TopicID `json:"topic_ids"`
	LanguageCode                string            `json:"language_code"`
}

type upsertUserContentTopicsResponse struct {
	Success bool               `json:"success"`
	Error   *clienterror.Error `json:"error,omitempty"`
}

func upsertUserContentTopics(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req upsertUserContentTopicsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return upsertUserContentTopicsResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
		}, nil
	}
	cErr, err := routermiddleware.ValidateUserAuth(userAuth, routermiddleware.ValidateUserAuthInput{
		DecodedUserID: *userID,
		EmailAddress:  req.EmailAddress,
	})
	switch {
	case err != nil:
		return nil, err
	case cErr != nil:
		return upsertUserContentTopicsResponse{
			Error: cErr,
		}, nil
	}
	// TODO(multiple-languages): eventually, I'll need to update this
	// so that it is language aware, since not all topics will apply to all languages
	if _, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode); err != nil {
		return upsertUserContentTopicsResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var mappings []usercontenttopics.ContentTopicWithTopicID
		for _, t := range req.TopicIDs {
			topic, err := content.GetContentTopicForTopicID(tx, t)
			if err != nil {
				return err
			}
			mappings = append(mappings, usercontenttopics.ContentTopicWithTopicID{
				Topic:   *topic,
				TopicID: t,
			})
		}
		return usercontenttopics.UpdateContentTopicsForUser(tx, *userID, mappings)
	}); err != nil {
		return nil, err
	}
	return upsertUserContentTopicsResponse{
		Success: true,
	}, nil
}
