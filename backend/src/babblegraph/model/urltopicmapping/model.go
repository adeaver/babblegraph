package urltopicmapping

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID             contentTopicMappingID      `db:"_id"`
	URLIdentifier  string                     `db:"url_identifier"`
	ContentTopic   contenttopics.ContentTopic `db:"content_topic"`
	TopicMappingID content.TopicMappingID     `db:"topic_mapping_id"`
}

// TODO(topic-migration): Remove this
type TopicMappingUnion struct {
	Topic          contenttopics.ContentTopic
	TopicMappingID content.TopicMappingID
}

func ApplyContentTopicsToURL(tx *sqlx.Tx, url string, topicUnions []TopicMappingUnion) error {
	parsedURL := urlparser.ParseURL(url)
	if parsedURL == nil {
		return fmt.Errorf("url is invalid: %s", url)
	}
	queryBuilder, err := database.NewBulkInsertQueryBuilder("content_topic_mappings", "url_identifier", "content_topic", "topic_mapping_id")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("DO NOTHING")
	for _, u := range topicUnions {
		queryBuilder.AddValues(parsedURL.URLIdentifier, u.Topic, u.TopicMappingID)
	}
	return queryBuilder.Execute(tx)
}

func GetTopicsForURL(tx *sqlx.Tx, url string) ([]contenttopics.ContentTopic, error) {
	parsedURL := urlparser.ParseURL(url)
	if parsedURL == nil {
		return nil, fmt.Errorf("url is invalid: %s", url)
	}
	var matches []dbContentTopicMapping
	if err := tx.Select(&matches, "SELECT * FROM content_topic_mappings WHERE url_identifier = $1", parsedURL.URLIdentifier); err != nil {
		return nil, err
	}
	var out []contenttopics.ContentTopic
	for _, m := range matches {
		out = append(out, m.ContentTopic)
	}
	return out, nil
}
