package blog

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Required Methods
// - Get all live blog post metadata
// - Get blog post metadata by url path
// - Get blog content by url path

const (
	getHeroImageForBlogPostQuery = "SELECT * FROM blog_post_image_metadata WHERE blog_id = $1 AND is_hero_image = TRUE"
)

func lookupHeroImageForBlogPost(tx *sqlx.Tx, blogID ID) (*dbImageMetadata, error) {
	var matches []dbImageMetadata
	err := tx.Select(&matches, getHeroImageForBlogPostQuery, blogID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most 1 result, but got %d", len(matches))
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		m := matches[0]
		return &m, nil
	default:
		panic("unreachable")
	}
}
