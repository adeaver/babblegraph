package content

import (
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getSourceForURLQuery           = "SELECT * FROM content_source WHERE url_identifier = $1"
	getSourceSeedForURLQuery       = "SELECT * FROM content_source_seed WHERE url_identifier = $1"
	getSourceSeedTopicMappingQuery = "SELECT * FROM content_source_seed_topic_mapping WHERE topic_id = $1 AND source_seed_id = $2"
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
	err := tx.Select(&matches, getSourceForURLQuery, u.URLIdentifier)
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
	err := tx.Select(&matches, getSourceSeedForURLQuery, u.URLIdentifier)
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
		if err != nil {
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
