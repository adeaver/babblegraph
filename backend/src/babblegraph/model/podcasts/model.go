package podcasts

import (
	"babblegraph/model/content"
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
