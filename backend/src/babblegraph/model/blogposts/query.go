package blogposts

import "github.com/jmoiron/sqlx"

const (
	getBlogByURLPathQuery    = "SELECT * FROM blog_posts WHERE url_path = $1 AND is_visible = TRUE"
	captureBlogPostViewQuery = "INSERT INTO blog_post_view (blog_post_id) VALUES ($1)"
	getAllBlogsPaginated     = "SELECT * FROM blog_posts WHERE is_visible = TRUE ORDER BY first_published_date DESC OFFSET $1"

	paginatedBlogsSize = 5
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

func GetAllBlogsPaginated(tx *sqlx.Tx, pageIndex int) ([]BlogPost, error) {
	var matches []dbBlogPost
	if err := tx.Select(&matches, getAllBlogsPaginated, paginatedBlogsSize*pageIndex); err != nil {
		return nil, err
	}
	var out []BlogPost
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}
