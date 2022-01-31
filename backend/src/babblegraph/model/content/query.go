package content

import (
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getSourceForURLQuery = "SELECT * FROM content_source WHERE title = $1"
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
		return nil, fmt.Errorf("Expected at most one topic, but got %d", len(matches))
	case len(matches) == 0:
		return nil, nil
	default:
		return matches[0].ID.Ptr(), nil
	}
}
