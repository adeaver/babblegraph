package blog

import (
	"babblegraph/model/utm"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getAllBlogPostMetadataQuery       = "SELECT * FROM blog_post_metadata ORDER BY created_at DESC"
	getBlogPostMetadataByURLPathQuery = "SELECT * FROM blog_post_metadata WHERE url_path = $1"
	addBlogPostMetadataQuery          = "INSERT INTO blog_post_metadata (title, description, author_name, url_path, status) VALUES ($1, $2, $3, $4, $5)"
	updateBlogPostMetadataQuery       = "UPDATE blog_post_metadata SET title=$1, description=$2, author_name=$3, last_modified_at = timezone('utc', now()) WHERE url_path = $4"
	updateBlogPostStatusQuery         = "UPDATE blog_post_metadata SET status=$1, last_modified_at=timezone('utc', now()) WHERE url_path = $2"
	updateBlogPostPublishedTimeQuery  = "UPDATE blog_post_metadata SET published_at=timezone('utc', now()), last_modified_at=timezone('utc', now()) WHERE url_path=$1"

	insertBlogImageQuery = "INSERT INTO blog_post_image_metadata (blog_id, path, file_name, alt_text, caption, is_hero_image) VALUES ($1, $2, $3, $4, $5, $6)"

	lookupBlogPostViewsForBlogQuery = "SELECT * FROM blog_post_view WHERE blog_id = $1"
)

func GetAllBlogPostMetadata(tx *sqlx.Tx) ([]BlogPostMetadata, error) {
	var matches []dbBlogPostMetadata
	if err := tx.Select(&matches, getAllBlogPostMetadataQuery); err != nil {
		return nil, err
	}
	var out []BlogPostMetadata
	for _, m := range matches {
		blogPost, err := m.ToNonDB(tx)
		if err != nil {
			return nil, err
		}
		out = append(out, *blogPost)
	}
	return out, nil
}

func getBlogPostMetadataByURLPath(tx *sqlx.Tx, urlPath string) (*dbBlogPostMetadata, error) {
	var matches []dbBlogPostMetadata
	if err := tx.Select(&matches, getBlogPostMetadataByURLPathQuery, urlPath); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, fmt.Errorf("Blog post %s not found", urlPath)
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected only one blog post, but found %d for blog post %s", len(matches), urlPath)
	default:
		m := matches[0]
		return &m, nil
	}
}

func GetBlogPostMetadataByURLPath(tx *sqlx.Tx, urlPath string) (*BlogPostMetadata, error) {
	metadata, err := getBlogPostMetadataByURLPath(tx, urlPath)
	if err != nil {
		return nil, err
	}
	return metadata.ToNonDB(tx)
}

type AddBlogPostMetadataInput struct {
	Title       string
	Description string
	AuthorName  string
	URLPath     string
}

func AddBlogPostMetadata(tx *sqlx.Tx, input AddBlogPostMetadataInput) error {
	if _, err := tx.Exec(addBlogPostMetadataQuery, input.Title, input.Description, input.AuthorName, input.URLPath, PostStatusDraft); err != nil {
		return err
	}
	return nil
}

type UpdateBlogPostMetadataInput struct {
	Title       string
	Description string
	AuthorName  string
	URLPath     string
}

func UpdateBlogPostMetadata(tx *sqlx.Tx, input UpdateBlogPostMetadataInput) error {
	if _, err := tx.Exec(updateBlogPostMetadataQuery, input.Title, input.Description, input.AuthorName, input.URLPath); err != nil {
		return err
	}
	return nil
}

func UpdateBlogPostStatus(tx *sqlx.Tx, urlPath string, status PostStatus) error {
	switch status {
	case PostStatusHidden,
		PostStatusDeleted,
		PostStatusDraft:
		if _, err := tx.Exec(updateBlogPostStatusQuery, status, urlPath); err != nil {
			return err
		}
	case PostStatusLive:
		if _, err := tx.Exec(updateBlogPostPublishedTimeQuery, urlPath); err != nil {
			return err
		}
		if _, err := tx.Exec(updateBlogPostStatusQuery, status, urlPath); err != nil {
			return err
		}
	}
	return nil
}

type InsertBlogImageMetadataInput struct {
	BlogID      ID
	Path        string
	FileName    string
	AltText     string
	Caption     *string
	IsHeroImage bool
}

func InsertBlogImageMetadata(tx *sqlx.Tx, input InsertBlogImageMetadataInput) error {
	if _, err := tx.Exec(insertBlogImageQuery, input.BlogID, input.Path, input.FileName, input.AltText, input.Caption, input.IsHeroImage); err != nil {
		return err
	}
	return nil
}

type BlogPostViewMetrics struct {
	TotalViews           int64 `json:"total_views"`
	UniqueViews          int64 `json:"unique_views"`
	LastMonthTotalViews  int64 `json:"last_month_total_views"`
	LastMonthUniqueViews int64 `json:"last_month_unique_views"`
}

func GetBlogPostViewMetrics(tx *sqlx.Tx, urlPath string) (*BlogPostViewMetrics, error) {
	blogPostMetadata, err := getBlogPostMetadataByURLPath(tx, urlPath)
	if err != nil {
		return nil, err
	}
	var matches []dbBlogPostView
	if err := tx.Select(&matches, lookupBlogPostViewsForBlogQuery, blogPostMetadata.ID); err != nil {
		return nil, err
	}
	uniqueViews := make(map[utm.TrackingID]bool)
	uniqueViewsInLastMonth := make(map[utm.TrackingID]bool)
	var totalViews, totalViewsInLastMonth int64
	lastMonth := time.Now().Add(-31 * 24 * time.Hour)
	for _, m := range matches {
		totalViews++
		isInLastMonth := m.CreatedAt.After(lastMonth)
		switch {
		case m.TrackingID == nil && isInLastMonth:
			totalViewsInLastMonth++
		case m.TrackingID != nil && !isInLastMonth:
			uniqueViews[*m.TrackingID] = true
		case m.TrackingID != nil && isInLastMonth:
			totalViewsInLastMonth++
			uniqueViewsInLastMonth[*m.TrackingID] = true
			uniqueViews[*m.TrackingID] = true
		default:
			// no-op
		}
	}
	return &BlogPostViewMetrics{
		TotalViews:           totalViews,
		UniqueViews:          int64(len(uniqueViews)),
		LastMonthTotalViews:  totalViewsInLastMonth,
		LastMonthUniqueViews: int64(len(uniqueViewsInLastMonth)),
	}, nil
}
