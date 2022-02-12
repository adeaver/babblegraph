package content

import (
	"babblegraph/model/contenttopics"
	"fmt"

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
		return nil, fmt.Errorf("Expected at most one topic, but got %d", len(matches))
	default:
		return matches[0].ID.Ptr(), nil
	}
}
