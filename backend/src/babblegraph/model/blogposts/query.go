package blogposts

import "github.com/jmoiron/sqlx"

const (
	getBlogByURLPathQuery    = "SELECT * FROM blog_posts WHERE url_path = $1 AND is_visible = TRUE"
	captureBlogPostViewQuery = "INSERT INTO blog_post_views (blog_post_id) VALUES ($1)"
)

func LookupBlogByURLPath(tx *sqlx.Tx, urlPath string) (*BlogPost, error) {
	var matches []dbBlogPost
	if err := tx.Select(&matches, getBlogByURLPathQuery, urlPath); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, nil
	}
	blogPost := matches[0].ToNonDB()
	return &blogPost, nil
}

func CaptureBlogPostView(tx *sqlx.Tx, blogPostID BlogPostID) error {
	if _, err := tx.Exec(captureBlogPostViewQuery, blogPostID); err != nil {
		return err
	}
	return nil
}
