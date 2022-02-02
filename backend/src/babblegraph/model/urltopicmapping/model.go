package urltopicmapping

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
)

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID             contentTopicMappingID      `db:"_id"`
	URLIdentifier  string                     `db:"url_identifier"`
	ContentTopic   contenttopics.ContentTopic `db:"content_topic"`
	TopicMappingID content.TopicMappingID     `db:"topic_mapping_id"`
}
