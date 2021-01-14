package usercontenttopics

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
)

type UserContentTopicMappingID string

type dbUserContentTopicMapping struct {
	ID           UserContentTopicMappingID  `db:"_id"`
	UserID       users.UserID               `db:"user_id"`
	ContentTopic contenttopics.ContentTopic `db:"content_topic"`
	IsActive     bool                       `db:"is_active"`
}
