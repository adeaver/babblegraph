package blogposts

import (
	"strings"
	"time"
)

type BlogPostID string

type BlogPost struct {
	ID                 BlogPostID
	Title              string
	Description        string
	Tags               []string
	TrackingTag        string
	URLPath            string
	HeroImageURL       string
	FirstPublishedDate time.Time
	UpdateDate         *time.Time
}

type dbBlogPost struct {
	ID                 BlogPostID `db:"_id"`
	Title              string     `db:"title"`
	Description        string     `db:"description"`
	Tags               string     `db:"tags"`
	TrackingTag        string     `db:"tracking_tag"`
	URLPath            string     `db:"url_path"`
	HeroImageURL       string     `db:"hero_image_url"`
	FirstPublishedDate time.Time  `db:"first_published_date"`
	UpdateDate         *time.Time `db:"updated_date"`
	IsVisible          bool       `db:"is_visible"`
}

func (d dbBlogPost) ToNonDB() BlogPost {
	return BlogPost{
		ID:                 d.ID,
		Title:              d.Title,
		Description:        d.Description,
		Tags:               strings.Split(d.Tags, ","),
		TrackingTag:        d.TrackingTag,
		URLPath:            d.URLPath,
		FirstPublishedDate: d.FirstPublishedDate,
		UpdateDate:         d.UpdateDate,
	}
}
