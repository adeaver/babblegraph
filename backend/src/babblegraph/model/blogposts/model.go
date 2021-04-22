package blogposts

import (
	"strings"
	"time"
)

type BlogPostID string

type BlogPost struct {
	ID                 BlogPostID `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	Tags               []string   `json:"tags"`
	TrackingTag        string     `json:"tracking_tag"`
	URLPath            string     `json:"url_path"`
	HeroImageURL       string     `json:"hero_image_url"`
	HeroImageAltText   string     `json:"hero_image_alt_text"`
	ContentURL         string     `json:"content_url"`
	FirstPublishedDate time.Time  `json:"first_published_date"`
	UpdateDate         *time.Time `json:"updated_date"`
}

type dbBlogPost struct {
	ID                 BlogPostID `db:"_id"`
	Title              string     `db:"title"`
	Description        string     `db:"description"`
	Tags               string     `db:"tags"`
	TrackingTag        string     `db:"tracking_tag"`
	URLPath            string     `db:"url_path"`
	HeroImageURL       string     `db:"hero_image_url"`
	HeroImageAltText   string     `db:"hero_image_alt_text"`
	ContentURL         string     `db:"content_url"`
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
		HeroImageURL:       d.HeroImageURL,
		HeroImageAltText:   d.HeroImageAltText,
		ContentURL:         d.ContentURL,
		FirstPublishedDate: d.FirstPublishedDate,
		UpdateDate:         d.UpdateDate,
	}
}
