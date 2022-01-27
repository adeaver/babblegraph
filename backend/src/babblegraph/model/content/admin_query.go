package content

import (
	"babblegraph/util/geo"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"strings"

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
            title,
            url,
            type,
            country,
            ingest_strategy,
            should_use_url_as_seed_url,
            is_active,
            monthly_access_limit
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        ) RETURNING _id`
	updateSourceQuery = `UPDATE
        content_source
    SET
        language_code=$1,
        title=$2,
        url=$3,
        type=$4,
        country=$5,
        ingest_strategy=$6,
        should_use_url_as_seed_url=$7,
        is_active=$8,
        monthly_access_limit=$9
    WHERE
        _id = $10
    `

	getAllSourceSeedsForSourceQuery = "SELECT * FROM content_source_seed WHERE root_id = $1"
	addSourceSeedQuery              = "INSERT INTO content_source_seed (root_id, url, is_active) VALUES ($1, $2, $3) RETURNING _id"
	updateSourceSeedQuery           = "UPDATE content_source_seed SET url=$1, is_active=$2 WHERE _id = $3"

	getAllSourceSeedTopicMappingsQuery = "SELECT * FROM content_source_seed_topic_mapping WHERE source_seed_id IN (?)"
	upsertSourceSeedTopicMapping       = `INSERT INTO
        content_source_seed_topic_mapping (
            source_seed_id, topic_id, is_active
        ) VALUES ($1, $2, $3) ON CONFLICT
        (source_seed_id, topic_id) DO UPDATE
        SET is_active = $3
    `

	getSourceFilterForSourceQuery    = "SELECT * FROM content_source_filter WHERE root_id = $1"
	upsertSourceFilterForSourceQuery = `INSERT INTO
        content_source_filter (
            root_id, is_active, use_ld_json_validation, paywall_classes, paywall_ids
        ) VALUES (
            $1, $2, $3, $4, $5
        ) ON CONFLICT (root_id) DO UPDATE SET
    is_active=$2, use_ld_json_validation=$3, paywall_classes=$4, paywall_ids=$5
    RETURNING _id`
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
	Title                 string
	LanguageCode          wordsmith.LanguageCode
	URL                   string
	Type                  SourceType
	IngestStrategy        IngestStrategy
	Country               geo.CountryCode
	ShouldUseURLAsSeedURL bool
	IsActive              bool
	MonthlyAccessLimit    *int64
}

func InsertSource(tx *sqlx.Tx, input InsertSourceInput) (*SourceID, error) {
	rows, err := tx.Query(insertSourceQuery, input.LanguageCode, input.Title, input.URL, input.Type, input.Country, input.IngestStrategy, input.ShouldUseURLAsSeedURL, input.IsActive, input.MonthlyAccessLimit)
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
	LanguageCode          wordsmith.LanguageCode
	Title                 string
	URL                   string
	Type                  SourceType
	IngestStrategy        IngestStrategy
	IsActive              bool
	MonthlyAccessLimit    *int64
	Country               geo.CountryCode
	ShouldUseURLAsSeedURL bool
}

func UpdateSource(tx *sqlx.Tx, id SourceID, input UpdateSourceInput) error {
	if _, err := tx.Exec(updateSourceQuery, input.LanguageCode, input.Title, input.URL, input.Type, input.Country, input.IngestStrategy, input.ShouldUseURLAsSeedURL, input.IsActive, input.MonthlyAccessLimit, id); err != nil {
		return err
	}
	return nil
}

func GetAllSourceSeedsForSource(tx *sqlx.Tx, sourceID SourceID) ([]SourceSeed, error) {
	var matches []dbSourceSeed
	if err := tx.Select(&matches, getAllSourceSeedsForSourceQuery, sourceID); err != nil {
		return nil, err
	}
	var out []SourceSeed
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func AddSourceSeed(tx *sqlx.Tx, sourceID SourceID, u urlparser.ParsedURL, isActive bool) (*SourceSeedID, error) {
	rows, err := tx.Query(addSourceSeedQuery, sourceID, u.URL, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sourceSeedID SourceSeedID
	for rows.Next() {
		if err := rows.Scan(&sourceSeedID); err != nil {
			return nil, err
		}
	}
	return &sourceSeedID, nil
}

func UpdateSourceSeed(tx *sqlx.Tx, sourceSeedID SourceSeedID, u urlparser.ParsedURL, isActive bool) error {
	if _, err := tx.Exec(updateSourceSeedQuery, u.URL, isActive, sourceSeedID); err != nil {
		return err
	}
	return nil
}

func GetAllSourceSeedTopicMappings(tx *sqlx.Tx, sourceSeedIDs []SourceSeedID) ([]SourceSeedTopicMapping, error) {
	query, args, err := sqlx.In(getAllSourceSeedTopicMappingsQuery, sourceSeedIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	var matches []dbSourceSeedTopicMapping
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	var out []SourceSeedTopicMapping
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func UpsertSourceSeedMapping(tx *sqlx.Tx, sourceSeedID SourceSeedID, topicID TopicID, isActive bool) error {
	if _, err := tx.Exec(upsertSourceSeedTopicMapping, sourceSeedID, topicID, isActive); err != nil {
		return err
	}
	return nil
}

func LookupSourceFilterForSource(tx *sqlx.Tx, sourceID SourceID) (*SourceFilter, error) {
	var matches []dbSourceFilter
	err := tx.Select(&matches, getSourceFilterForSourceQuery, sourceID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		m := matches[0].ToNonDB()
		return &m, nil
	default:
		return nil, fmt.Errorf("Expected at most 1 source filter for source ID %s, but got %d", sourceID, len(matches))
	}
}

type UpsertSourceFilterForSourceInput struct {
	IsActive            bool
	PaywallClasses      []string
	PaywallIDs          []string
	UseLDJSONValidation *bool
}

func UpsertSourceFilterForSource(tx *sqlx.Tx, sourceID SourceID, input UpsertSourceFilterForSourceInput) (*SourceFilterID, error) {
	var paywallClassesStr, paywallIDsStr *string
	if len(input.PaywallClasses) > 0 {
		paywallClassesStr = ptr.String(strings.Join(input.PaywallClasses, paywallFilterDelimiter))
	}
	if len(input.PaywallIDs) > 0 {
		paywallIDsStr = ptr.String(strings.Join(input.PaywallIDs, paywallFilterDelimiter))
	}
	rows, err := tx.Query(upsertSourceFilterForSourceQuery, sourceID, input.IsActive, input.UseLDJSONValidation, paywallClassesStr, paywallIDsStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sourceFilterID SourceFilterID
	for rows.Next() {
		if err := rows.Scan(&sourceFilterID); err != nil {
			return nil, err
		}
	}
	return &sourceFilterID, nil
}
