package content

import (
	"babblegraph/model/admin"
	"babblegraph/model/content"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "content",
	Routes: []router.Route{
		{
			Path: "get_all_topics_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				getAllContentTopics,
			),
		}, {
			Path: "add_topic_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				addContentTopic,
			),
		}, {
			Path: "get_topic_by_id_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				getTopicByID,
			),
		}, {
			Path: "update_is_topic_active_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				updateIsContentTopicActive,
			),
		},
	},
}

type getAllContentTopicsRequest struct{}

type getAllContentTopicsResponse struct {
	Topics []content.Topic `json:"topics"`
}

func getAllContentTopics(adminID admin.ID, r *router.Request) (interface{}, error) {
	var topics []content.Topic
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topics, err = content.GetAllTopics(tx)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllContentTopicsResponse{
		Topics: topics,
	}, nil
}

type getTopicByIDRequest struct {
	ID content.TopicID `json:"id"`
}

type getTopicByIDResponse struct {
	Topic content.Topic `json:"topic"`
}

func getTopicByID(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getTopicByIDRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var topic *content.Topic
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topic, err = content.GetTopic(tx, req.ID)
		return err
	}); err != nil {
		return nil, err
	}
	return getTopicByIDResponse{
		Topic: *topic,
	}, nil
}

type addContentTopicRequest struct {
	Label string `json:"label"`
}

type addContentTopicResponse struct {
	TopicID content.TopicID `json:"topic_id"`
}

func addContentTopic(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addContentTopicRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var topicID *content.TopicID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topicID, err = content.AddTopic(tx, req.Label, false)
		return err
	}); err != nil {
		return nil, err
	}
	return addContentTopicResponse{
		TopicID: *topicID,
	}, nil
}

type updateIsContentTopicActiveRequest struct {
	ID       content.TopicID `json:"id"`
	IsActive bool            `json:"is_active"`
}

type updateIsContentTopicActiveResponse struct {
	Success bool `json:"success"`
}

func updateIsContentTopicActive(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateIsContentTopicActiveRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return content.ToggleTopicIsActive(tx, req.ID, req.IsActive)
	}); err != nil {
		return nil, err
	}
	return updateIsContentTopicActiveResponse{
		Success: true,
	}, nil
}
