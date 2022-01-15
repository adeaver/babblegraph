package content

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getAllTopicsQuery      = "SELECT * FROM content_topic"
	getTopicQuery          = "SELECT * FROM content_topic WHERE _id = $1"
	insertTopicQuery       = "INSERT INTO content_topic (label, is_active) VALUES ($1, $2)"
	toggleTopicActiveQuery = "UPDATE content_topic SET is_active = $1 WHERE _id = $2"
	updateTopicLabelQuery  = "UPDATE content_topic SET label = $1 WHERE _id = $2"
)

func GetAllTopics(tx *sqlx.Tx) ([]Topic, error) {
	var matches []dbTopic
	if err := tx.Select(&matches, getAllTopicsQuery); err != nil {
		return nil, err
	}
	var out []Topic
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func GetTopic(tx *sqlx.Tx, id TopicID) (*Topic, error) {
	var matches []dbTopic
	err := tx.Select(&matches, getTopicQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected 1 topic for ID %s, but got %d", id, len(matches))
	default:
		m := matches[0].ToNonDB()
		return &m, nil
	}
}

func AddTopic(tx *sqlx.Tx, label string, isActive bool) (*TopicID, error) {
	rows, err := tx.Query(insertTopicQuery, label, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var topicID TopicID
	for rows.Next() {
		if err := rows.Scan(&topicID); err != nil {
			return nil, err
		}
	}
	return &topicID, nil
}

func ToggleTopicIsActive(tx *sqlx.Tx, id TopicID, isActive bool) error {
	if _, err := tx.Exec(toggleTopicActiveQuery, isActive, id); err != nil {
		return err
	}
	return nil
}

func UpdateTopicLabel(tx *sqlx.Tx, id TopicID, label string) error {
	if _, err := tx.Exec(updateTopicLabelQuery, label, id); err != nil {
		return err
	}
	return nil
}
