package blog

import (
	"babblegraph/model/utm"
	"time"

	"github.com/jmoiron/sqlx"
)

type ID string

type dbBlogPostMetadata struct {
	ID             ID         `db:"_id"`
	CreatedAt      time.Time  `db:"created_at"`
	LastModifiedAt time.Time  `db:"last_modified_at"`
	PublishedAt    *time.Time `db:"published_at"`
	Title          string     `db:"title"`
	AuthorName     string     `db:"author_name"`
	Description    string     `db:"description"`
	URLPath        string     `db:"url_path"`
	Status         PostStatus `db:"status"`
}

func (d dbBlogPostMetadata) ToNonDB(tx *sqlx.Tx) (*BlogPostMetadata, error) {
	heroImage, err := lookupHeroImageForBlogPost(tx, d.ID)
	if err != nil {
		return nil, err
	}
	var i *Image
	if heroImage != nil {
		temp := heroImage.ToNonDB()
		i = &temp
	}
	return &BlogPostMetadata{
		ID:          d.ID,
		PublishedAt: d.PublishedAt,
		Title:       d.Title,
		Description: d.Description,
		URLPath:     d.URLPath,
		Status:      d.Status,
		AuthorName:  d.AuthorName,
		HeroImage:   i,
	}, nil
}

type BlogPostMetadata struct {
	ID          ID         `json:"id"`
	PublishedAt *time.Time `json:"published_at"`
	HeroImage   *Image     `json:"hero_image,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URLPath     string     `json:"url_path"`
	Status      PostStatus `json:"status"`
	AuthorName  string     `json:"author_name"`
}

type PostStatus string

const (
	PostStatusDraft   PostStatus = "draft"
	PostStatusLive    PostStatus = "live"
	PostStatusHidden  PostStatus = "hidden"
	PostStatusDeleted PostStatus = "deleted"
)

type imageID string

type dbImageMetadata struct {
	CreatedAt      time.Time `db:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at"`
	ID             imageID   `db:"_id"`
	Path           string    `db:"path"`
	BlogID         ID        `db:"blog_id"`
	FileName       string    `db:"file_name"`
	AltText        string    `db:"alt_text"`
	Caption        *string   `db:"caption"`
	IsHeroImage    bool      `db:"is_hero_image"`
}

func (d dbImageMetadata) ToNonDB() Image {
	return Image{
		Path:    d.Path,
		AltText: d.AltText,
		Caption: d.Caption,
	}
}

type blogPostViewID string

type dbBlogPostView struct {
	CreatedAt      time.Time       `db:"created_at"`
	LastModifiedAt time.Time       `db:"last_modified_at"`
	ID             blogPostViewID  `db:"_id"`
	BlogID         ID              `db:"blog_id"`
	TrackingID     *utm.TrackingID `db:"tracking_id"`
}
