package blog

import "time"

type dbBlogPostMetadata struct {
	CreatedAt      time.Time  `db:"created_at"`
	LastModifiedAt time.Time  `db:"last_modified_at"`
	PublishedAt    time.Time  `db:"published_at"`
	HeroImagePath  *string    `db:"hero_image_path"`
	Title          string     `db:"title"`
	AuthorName     string     `db:"author_name"`
	Description    string     `db:"description"`
	URLPath        string     `db:"url_path"`
	Status         PostStatus `db:"status"`
}

func (d dbBlogPostMetadata) ToNonDB() BlogPostMetadata {
	return BlogPostMetadata{
		PublishedAt:   d.PublishedAt,
		HeroImagePath: d.HeroImagePath,
		Title:         d.Title,
		Description:   d.Description,
		URLPath:       d.URLPath,
		Status:        d.Status,
		AuthorName:    d.AuthorName,
	}
}

type BlogPostMetadata struct {
	PublishedAt   time.Time  `json:"published_at"`
	HeroImagePath *string    `json:"hero_image_path,omitempty"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	URLPath       string     `json:"url_path"`
	Status        PostStatus `json:"status"`
	AuthorName    string     `json:"author_name"`
}

type PostStatus string

const (
	PostStatusDraft   PostStatus = "draft"
	PostStatusLive    PostStatus = "live"
	PostStatusHidden  PostStatus = "hidden"
	PostStatusDeleted PostStatus = "deleted"
)
