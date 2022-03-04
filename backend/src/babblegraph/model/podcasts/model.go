package podcasts

import (
	"babblegraph/model/content"
	"babblegraph/wordsmith"
	"time"
)

type ID string

type dbPodcastMetadata struct {
	CreatedAt      time.Time        `db:"created_at"`
	LastModifiedAt time.Time        `db:"last_modified_at"`
	ID             ID               `db:"_id"`
	ContentID      content.SourceID `db:"content_id"`
	ImageURL       *string          `db:"image_url"`
}

type Version int

const (
	Version1 Version = 1

	CurrentVersion = Version1
)

type EpisodeID string

type Episode struct {
	ID              EpisodeID     `json:"id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	PublicationDate time.Time     `json:"publication_date"`
	EpisodeType     string        `json:"episode_type"`
	Duration        time.Duration `json:"duration"`
	IsExplicit      bool          `json:"is_explicit"`
	AudioFile       AudioFile     `json:"audio_file"`

	// This is only guaranteed to be unique across
	// episodes of the same podcasts, not universally across
	// podcasts
	GUID string `json:"guid"`

	Version Version `json:"verison"`

	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	TopicIDs     []content.TopicID      `json:"topic_ids"`
	SourceID     content.SourceID       `json:"source_id"`
}

type AudioFile struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}
