package content

import (
	"babblegraph/util/geo"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getAllTopicsQuery      = "SELECT * FROM content_topic"
	getTopicQuery          = "SELECT * FROM content_topic WHERE _id = $1"
	insertTopicQuery       = "INSERT INTO content_topic (label, is_active) VALUES ($1, $2) RETURNING _id"
	toggleTopicActiveQuery = "UPDATE content_topic SET is_active = $1 WHERE _id = $2"
	updateTopicLabelQuery  = "UPDATE content_topic SET label = $1 WHERE _id = $2"

	getAllTopicDipslayNamesForTopicQuery = "SELECT * FROM content_topic_display_name WHERE topic_id = $1"
	insertTopicDisplayNameForTopicQuery  = "INSERT INTO content_topic_display_name (topic_id, language_code, label, is_active) VALUES ($1, $2, $3, $4) RETURNING _id"
	updateTopicDisplayNameLabelQuery     = "UPDATE content_topic_display_name SET label = $1 WHERE _id = $2"
	updateTopicDisplayNameIsActiveQuery  = "UPDATE content_topic_display_name SET is_active = $1 WHERE _id = $2"

	getAllSourcesQuery = "SELECT * FROM content_source"
	getSourceQuery     = "SELECT * FROM content_source WHERE _id = $1"
	insertSourceQuery  = `INSERT INTO
        content_source (
            language_code,
            url,
            type,
            country,
            ingest_strategy,
            is_active,
            monthly_access_limit
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        ) RETURNING _id`
	updateSourceQuery = `UPDATE
        content_source
    SET
        url=$1,
        type=$2,
        country=$3,
        ingest_strategy=$4,
        is_active=$5,
        monthly_access_limit=$6
    WHERE
        _id = $7
    `
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
	return topicID.Ptr(), nil
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

func GetAllTopicDipslayNamesForTopic(tx *sqlx.Tx, topicID TopicID) ([]TopicDisplayName, error) {
	var matches []dbTopicDisplayName
	if err := tx.Select(&matches, getAllTopicDipslayNamesForTopicQuery, topicID); err != nil {
		return nil, err
	}
	var out []TopicDisplayName
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func AddTopicDisplayName(tx *sqlx.Tx, topicID TopicID, languageCode wordsmith.LanguageCode, label string, isActive bool) (*TopicDisplayNameID, error) {
	rows, err := tx.Query(insertTopicDisplayNameForTopicQuery, topicID, languageCode, label, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var topicDisplayNameID TopicDisplayNameID
	for rows.Next() {
		if err := rows.Scan(&topicDisplayNameID); err != nil {
			return nil, err
		}
	}
	return &topicDisplayNameID, nil
}

func UpdateTopicDisplayNameLabel(tx *sqlx.Tx, topicDisplayNameID TopicDisplayNameID, label string) error {
	if _, err := tx.Exec(updateTopicDisplayNameLabelQuery, label, topicDisplayNameID); err != nil {
		return err
	}
	return nil
}

func ToggleTopicDisplayNameIsActive(tx *sqlx.Tx, topicDisplayNameID TopicDisplayNameID, isActive bool) error {
	if _, err := tx.Exec(updateTopicDisplayNameIsActiveQuery, isActive, topicDisplayNameID); err != nil {
		return err
	}
	return nil
}

func GetAllSources(tx *sqlx.Tx) ([]Source, error) {
	var matches []dbSource
	if err := tx.Select(&matches, getAllSourcesQuery); err != nil {
		return nil, err
	}
	var out []Source
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func GetSource(tx *sqlx.Tx, id SourceID) (*Source, error) {
	var matches []dbSource
	err := tx.Select(&matches, getSourceQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected 1 match, but got %d for source ID %s", len(matches), id)
	default:
		m := matches[0].ToNonDB()
		return &m, nil
	}
}

type InsertSourceInput struct {
	LanguageCode       wordsmith.LanguageCode
	URL                string
	Type               SourceType
	IngestStrategy     IngestStrategy
	Country            geo.CountryCode
	MonthlyAccessLimit *int64
}

func InsertSource(tx *sqlx.Tx, input InsertSourceInput) (*SourceID, error) {
	rows, err := tx.Query(insertSourceQuery, input.LanguageCode, input.URL, input.Type, input.Country, input.IngestStrategy, false, input.MonthlyAccessLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sourceID SourceID
	for rows.Next() {
		if err := rows.Scan(&sourceID); err != nil {
			return nil, err
		}
	}
	return &sourceID, nil
}

type UpdateSourceInput struct {
	URL                string
	Type               SourceType
	IngestStrategy     IngestStrategy
	IsActive           bool
	MonthlyAccessLimit *int64
	Country            geo.CountryCode
}

func UpdateSource(tx *sqlx.Tx, id SourceID, input UpdateSourceInput) error {
	if _, err := tx.Exec(updateSourceQuery, input.URL, input.Type, input.Country, input.IngestStrategy, input.IsActive, input.MonthlyAccessLimit, id); err != nil {
		return err
	}
	return nil
}
