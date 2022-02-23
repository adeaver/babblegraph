package urltopicmapping

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TODO(topic-migration): Remove this
type TopicMappingUnion struct {
	Topic          contenttopics.ContentTopic
	TopicMappingID content.TopicMappingID
}

func ApplyContentTopicsToURL(tx *sqlx.Tx, parsedURL urlparser.ParsedURL, topicUnions []TopicMappingUnion) error {
	queryBuilder, err := database.NewBulkInsertQueryBuilder("content_topic_mappings", "url_identifier", "content_topic", "topic_mapping_id")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("(url_identifier, content_topic) DO UPDATE SET topic_mapping_id = EXCLUDED.topic_mapping_id")
	for _, u := range topicUnions {
		queryBuilder.AddValues(parsedURL.URLIdentifier, u.Topic, u.TopicMappingID)
	}
	return queryBuilder.Execute(tx)
}

func GetTopicsAndMappingIDsForURL(tx *sqlx.Tx, url string) ([]contenttopics.ContentTopic, []content.TopicMappingID, error) {
	parsedURL := urlparser.ParseURL(url)
	if parsedURL == nil {
		return nil, nil, fmt.Errorf("url is invalid: %s", url)
	}
	var matches []dbContentTopicMapping
	if err := tx.Select(&matches, "SELECT * FROM content_topic_mappings WHERE url_identifier = $1", parsedURL.URLIdentifier); err != nil {
		return nil, nil, err
	}
	var topics []contenttopics.ContentTopic
	var topicMappingIDs []content.TopicMappingID
	for _, m := range matches {
		topics = append(topics, m.ContentTopic)
		if m.TopicMappingID != nil {
			topicMappingIDs = append(topicMappingIDs, *m.TopicMappingID)
		}
	}
	return topics, topicMappingIDs, nil
}
