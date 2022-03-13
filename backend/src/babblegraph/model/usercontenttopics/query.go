package usercontenttopics

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

const (
	getContentTopicsForUserQuery              = "SELECT * FROM user_content_topic_mappings WHERE user_id = $1 AND is_active = TRUE"
	setAllContentTopicsToInactiveForUserQuery = "UPDATE user_content_topic_mappings SET is_active = FALSE WHERE user_id = $1"
)

func GetTopicIDsForUser(tx *sqlx.Tx, userID users.UserID) ([]content.TopicID, error) {
	var matches []dbUserContentTopicMapping
	if err := tx.Select(&matches, getContentTopicsForUserQuery, userID); err != nil {
		return nil, err
	}
	var out []content.TopicID
	for _, m := range matches {
		if m.ContentTopicID == nil {
			continue
		}
		out = append(out, *m.ContentTopicID)
	}
	return out, nil
}

// TODO: Deprecate
func GetContentTopicsForUser(tx *sqlx.Tx, userID users.UserID) ([]contenttopics.ContentTopic, error) {
	var matches []dbUserContentTopicMapping
	if err := tx.Select(&matches, getContentTopicsForUserQuery, userID); err != nil {
		return nil, err
	}
	var out []contenttopics.ContentTopic
	for _, m := range matches {
		out = append(out, m.ContentTopic)
	}
	return out, nil
}

type ContentTopicWithTopicID struct {
	Topic   contenttopics.ContentTopic
	TopicID content.TopicID
}

func UpdateContentTopicsForUser(tx *sqlx.Tx, userID users.UserID, contentTopics []ContentTopicWithTopicID) error {
	if _, err := tx.Exec(setAllContentTopicsToInactiveForUserQuery, userID); err != nil {
		return err
	}
	queryBuilder, err := database.NewBulkInsertQueryBuilder("user_content_topic_mappings", "user_id", "content_topic", "content_topic_id", "is_active")
	if err != nil {
		return err
	}
	queryBuilder.AddConflictResolution("(user_id, content_topic) DO UPDATE SET is_active = TRUE")
	for _, t := range contentTopics {
		if err := queryBuilder.AddValues(userID, t.Topic, t.TopicID, true); err != nil {
			return err
		}
	}
	return queryBuilder.Execute(tx)
}
