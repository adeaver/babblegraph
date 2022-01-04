package blog

import (
	"babblegraph/model/utm"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getHeroImageForBlogPostQuery             = "SELECT * FROM blog_post_image_metadata WHERE blog_id = $1 AND is_hero_image = TRUE"
	getClientBlogPostMetadataForURLPathQuery = "SELECT * FROM blog_post_metadata WHERE url_path = $1 AND status = $2"
	registerBlogViewQuery                    = "INSERT INTO blog_post_view (blog_id, tracking_id) VALUES ($1, $2)"
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

func GetClientBlogPostMetadataForURLPath(tx *sqlx.Tx, blogURLPath string) (*BlogPostMetadata, error) {
	var matches []dbBlogPostMetadata
	err := tx.Select(&matches, getClientBlogPostMetadataForURLPathQuery, blogURLPath, PostStatusLive)
	switch {
	case err != nil:
		return nil, err
	case len(matches) != 1:
		return nil, fmt.Errorf("Expected exactly 1 result, but got %d", len(matches))
	case len(matches) == 1:
		return matches[0].ToNonDB(tx)
	default:
		panic("unreachable")
	}
}

func RegisterBlogView(tx *sqlx.Tx, blogID ID, trackingID *utm.TrackingID) error {
	if _, err := tx.Exec(registerBlogViewQuery, blogID, trackingID); err != nil {
		return err
	}
	return nil
}
