package usercontenttopics

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
)

type UserContentTopicMappingID string

type dbUserContentTopicMapping struct {
	ID             UserContentTopicMappingID  `db:"_id"`
	UserID         users.UserID               `db:"user_id"`
	ContentTopic   contenttopics.ContentTopic `db:"content_topic"`
	ContentTopicID *content.TopicID           `db:"content_topic_id"`
	IsActive       bool                       `db:"is_active"`
}
