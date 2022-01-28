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
	var matches []dbSource
	err := tx.Select(&matches, getSourceForURLQuery, u.Domain)
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
