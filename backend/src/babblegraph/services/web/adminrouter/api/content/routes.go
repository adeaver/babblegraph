package content

import (
	"babblegraph/model/admin"
	"babblegraph/model/content"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"strings"

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
		}, {
			Path: "get_all_topic_display_names_for_topic_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				getAllTopicDisplayNamesForTopic,
			),
		}, {
			Path: "add_topic_display_name_for_topic_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				addTopicDisplayNameForTopic,
			),
		}, {
			Path: "update_topic_display_name_label_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				updateTopicDisplayNameLabel,
			),
		}, {
			Path: "toggle_topic_display_name_is_active_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentTopics,
				toggleTopicDisplayNameIsActive,
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
		topicID, err = content.AddTopic(tx, strings.ToLower(req.Label), false)
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

type getAllTopicDisplayNamesForTopicRequest struct {
	TopicID content.TopicID `json:"topic_id"`
}

type getAllTopicDisplayNamesForTopicResponse struct {
	TopicDisplayNames []content.TopicDisplayName `json:"topic_display_names"`
}

func getAllTopicDisplayNamesForTopic(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getAllTopicDisplayNamesForTopicRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var displayNames []content.TopicDisplayName
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		displayNames, err = content.GetAllTopicDipslayNamesForTopic(tx, req.TopicID)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllTopicDisplayNamesForTopicResponse{
		TopicDisplayNames: displayNames,
	}, nil
}

type addTopicDisplayNameForTopicRequest struct {
	TopicID      content.TopicID `json:"topic_id"`
	Label        string          `json:"label"`
	LanguageCode string          `json:"language_code"`
}

type addTopicDisplayNameForTopicResponse struct {
	TopicDisplayNameID content.TopicDisplayNameID `json:"topic_display_name_id"`
}

func addTopicDisplayNameForTopic(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addTopicDisplayNameForTopicRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	var displayNameID *content.TopicDisplayNameID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		displayNameID, err = content.AddTopicDisplayName(tx, req.TopicID, *languageCode, req.Label, false)
		return err
	}); err != nil {
		return nil, err
	}
	return addTopicDisplayNameForTopicResponse{
		TopicDisplayNameID: *displayNameID,
	}, nil
}

type updateTopicDisplayNameLabelRequest struct {
	TopicDisplayNameID content.TopicDisplayNameID `json:"topic_display_name_id"`
	Label              string                     `json:"label"`
}

type updateTopicDisplayNameLabelResponse struct {
	Success bool `json:"success"`
}

func updateTopicDisplayNameLabel(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateTopicDisplayNameLabelRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return content.UpdateTopicDisplayNameLabel(tx, req.TopicDisplayNameID, strings.ToLower(req.Label))
	}); err != nil {
		return nil, err
	}
	return updateTopicDisplayNameLabelResponse{
		Success: true,
	}, nil
}

type toggleTopicDisplayNameIsActiveRequest struct {
	TopicDisplayNameID content.TopicDisplayNameID `json:"topic_display_name_id"`
	IsActive           bool                       `json:"is_active"`
}

type toggleTopicDisplayNameIsActiveResponse struct {
	Success bool `json:"success"`
}

func toggleTopicDisplayNameIsActive(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req toggleTopicDisplayNameIsActiveRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return content.ToggleTopicDisplayNameIsActive(tx, req.TopicDisplayNameID, req.IsActive)
	}); err != nil {
		return nil, err
	}
	return toggleTopicDisplayNameIsActiveResponse{
		Success: true,
	}, nil
}
