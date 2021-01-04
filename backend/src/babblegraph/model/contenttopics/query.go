package contenttopics

import (
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ApplyContentTopicToURL(tx *sqlx.Tx, url string, topic ContentTopic) error {
	parsedURL := urlparser.ParseURL(url)
	if parsedURL == nil {
		return fmt.Errorf("url is invalid: %s", url)
	}
	_, err := tx.Exec("INSERT INTO content_topic_mappings (url_identifier, content_topic) VALUES ($1, $2)", parsedURL.URLIdentifier, topic)
	return err
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
