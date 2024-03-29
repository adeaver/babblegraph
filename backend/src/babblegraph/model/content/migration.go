package content

import (
	"babblegraph/model/contenttopics"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	getTopicByLabelQuery = "SELECT * FROM content_topic WHERE label = $1"
)

func GetTopicIDByContentTopic(tx *sqlx.Tx, t contenttopics.ContentTopic) (*TopicID, error) {
	var matches []dbTopic
	err := tx.Select(&matches, getTopicByLabelQuery, t)
	switch {
	case err != nil:
		return nil, err
	case len(matches) > 1,
		len(matches) == 0:
		return nil, fmt.Errorf("Expected exactly one topic for topic %s, but got %d", t, len(matches))
	default:
		return matches[0].ID.Ptr(), nil
	}
}

func GetContentTopicForTopicID(tx *sqlx.Tx, topicID TopicID) (*contenttopics.ContentTopic, error) {
	t, err := GetTopic(tx, topicID)
	if err != nil {
		return nil, err
	}
	return contenttopics.GetContentTopicForString(strings.ToLower(t.Label))
}
