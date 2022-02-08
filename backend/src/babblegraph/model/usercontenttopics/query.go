package usercontenttopics

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

const (
	getContentTopicsForUserQuery              = "SELECT * FROM user_content_topic_mappings WHERE user_id = $1 AND is_active = TRUE"
	setAllContentTopicsToInactiveForUserQuery = "UPDATE user_content_topic_mappings SET is_active = FALSE WHERE user_id = $1"
)

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

// TODO(content-migration): Remove this
func BackfillUserContentTopicMappings(c ctx.LogContext, tx *sqlx.Tx) error {
	rows, err := tx.Query("SELECT * FROM user_content_topic_mappings WHERE content_topic_id IS NULL")
	if err != nil {
		return err
	}
	defer rows.Close()
	var count int64
	c.Infof("Starting user topic mappings update")
	for rows.Next() {
		var match dbUserContentTopicMapping
		if err := rows.Scan(&match); err != nil {
			return err
		}
		topicID, err := content.GetTopicIDByContentTopic(tx, match.ContentTopic)
		if err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE user_content_topic_mappings SET content_topic = $1", *topicID); err != nil {
			return err
		}
		count++
		if count%1000 == 0 {
			c.Infof("Successfully completed %d mapping updates", count)
		}
	}
	c.Infof("Finished user topic mappings update")
	return nil
}
