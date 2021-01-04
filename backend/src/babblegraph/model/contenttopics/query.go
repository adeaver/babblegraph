package contenttopics

import (
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ApplyContentTopicsToURL(tx *sqlx.Tx, url string, topics []ContentTopic) error {
	parsedURL := urlparser.ParseURL(url)
	if parsedURL == nil {
		return fmt.Errorf("url is invalid: %s", url)
	}
	queryBuilder, err := database.NewBulkInsertQueryBuilder("content_topic_mappings", "url_identifier", "content_topic")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("DO NOTHING")
	for _, t := range topics {
		queryBuilder.AddValues(parsedURL.URLIdentifier, t)
	}
	return queryBuilder.Execute(tx)
}

func GetTopicsForURL(tx *sqlx.Tx, url string) ([]ContentTopic, error) {
	parsedURL := urlparser.ParseURL(url)
	if parsedURL == nil {
		return nil, fmt.Errorf("url is invalid: %s", url)
	}
	var matches []dbContentTopicMapping
	if err := tx.Select(&matches, "SELECT * FROM content_topic_mappings WHERE url_identifier = $1", parsedURL.URLIdentifier); err != nil {
		return nil, err
	}
	var out []ContentTopic
	for _, m := range matches {
		out = append(out, m.ContentTopic)
	}
	return out, nil
}
