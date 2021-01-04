package contenttopics

type ContentTopic string

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID            contentTopicMappingID `db:"_id"`
	URLIdentifier string                `db:"url_identifier"`
	ContentTopic  ContentTopic          `db:"content_topic"`
}
