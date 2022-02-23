package content

import (
	"babblegraph/util/ctx"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getSourceForURLQuery           = "SELECT * FROM content_source WHERE url_identifier = $1"
	getSourceSeedForURLQuery       = "SELECT * FROM content_source_seed WHERE url_identifier = $1 AND url_params IS NOT DISTINCT FROM $2"
	getSourceSeedTopicMappingQuery = "SELECT * FROM content_source_seed_topic_mapping WHERE topic_id = $1 AND source_seed_id = $2"

	getSourceSeedForSourceQuery                  = "SELECT * FROM content_source_seed WHERE root_id = $1"
	getSourceSeedTopicMappingForSourceSeedsQuery = "SELECT * FROM content_source_seed_topic_mapping WHERE topic_id = '%s' AND source_seed_id IN (?)"
)

func GetSourceIDForParsedURL(tx *sqlx.Tx, u urlparser.ParsedURL) (*SourceID, error) {
	sourceID, err := LookupSourceIDForParsedURL(tx, u)
	switch {
	case err != nil:
		return nil, err
	case sourceID == nil:
		return nil, fmt.Errorf("Expected exactly one source ID, but got none")
	}
	return sourceID, nil
}

func LookupSourceIDForParsedURL(tx *sqlx.Tx, u urlparser.ParsedURL) (*SourceID, error) {
	var matches []dbSource
	err := tx.Select(&matches, getSourceForURLQuery, u.Domain)
	switch {
	case err != nil:
		return nil, err
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one source, but got %d", len(matches))
	case len(matches) == 0:
		return nil, nil
	default:
		return matches[0].ID.Ptr(), nil
	}
}

func lookupSourceSeedIDForParsedURL(tx *sqlx.Tx, u urlparser.ParsedURL) (*SourceSeedID, error) {
	var matches []dbSourceSeed
	err := tx.Select(&matches, getSourceSeedForURLQuery, u.URLIdentifier, u.Params)
	switch {
	case err != nil:
		return nil, err
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one source seed, but got %d", len(matches))
	case len(matches) == 0:
		return nil, nil
	default:
		return matches[0].ID.Ptr(), nil
	}
}

func lookupSourceSeedMappingID(tx *sqlx.Tx, sourceSeedID SourceSeedID, topicID TopicID) (*SourceSeedTopicMappingID, error) {
	var matches []dbSourceSeedTopicMapping
	err := tx.Select(&matches, getSourceSeedTopicMappingQuery, topicID, sourceSeedID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one source seed topic mapping, but got %d", len(matches))
	case len(matches) == 0:
		return nil, nil
	default:
		return matches[0].ID.Ptr(), nil
	}
}

func LookupTopicMappingIDForURL(tx *sqlx.Tx, u urlparser.ParsedURL, topicID TopicID) (*TopicMappingID, error) {
	sourceSeedID, err := lookupSourceSeedIDForParsedURL(tx, u)
	switch {
	case err != nil:
		return nil, err
	case sourceSeedID != nil:
		sourceSeedTopicMappingID, err := lookupSourceSeedMappingID(tx, *sourceSeedID, topicID)
		switch {
		case err != nil:
			return nil, err
		case sourceSeedTopicMappingID == nil:
			return nil, nil
		}
		return MustMakeTopicMappingID(MakeTopicMappingIDInput{
			SourceSeedTopicMappingID: sourceSeedTopicMappingID,
		}).Ptr(), nil
	default:
		// TODO: implement same logic with source mapping
		return nil, nil
	}
}

func LookupTopicMappingIDForSourceAndTopic(c ctx.LogContext, tx *sqlx.Tx, sourceID SourceID, topicID TopicID) (*TopicMappingID, error) {
	var matches []dbSourceSeed
	err := tx.Select(&matches, getSourceSeedForSourceQuery, sourceID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	default:
		var sourceSeedIDs []SourceSeedID
		for _, m := range matches {
			sourceSeedIDs = append(sourceSeedIDs, m.ID)
		}
		query, args, err := sqlx.In(fmt.Sprintf(getSourceSeedTopicMappingForSourceSeedsQuery, topicID), sourceSeedIDs)
		if err != nil {
			return nil, err
		}
		sql := tx.Rebind(query)
		var mappings []dbSourceSeedTopicMapping
		err = tx.Select(&mappings, sql, args...)
		switch {
		case err != nil:
			return nil, err
		case len(mappings) == 0:
			return nil, nil
		case len(mappings) == 1:
			return MustMakeTopicMappingID(MakeTopicMappingIDInput{
				SourceSeedTopicMappingID: mappings[0].ID.Ptr(),
			}).Ptr(), nil
		default:
			c.Infof("Found %d mappings for source ID %s and topic %s: %+v, choosing the first one", len(mappings), sourceID, topicID, mappings)
			return MustMakeTopicMappingID(MakeTopicMappingIDInput{
				SourceSeedTopicMappingID: mappings[0].ID.Ptr(),
			}).Ptr(), nil
		}
	}
}
