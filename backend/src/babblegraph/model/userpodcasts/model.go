package userpodcasts

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/users"
	"babblegraph/util/env"
	"fmt"
	"time"
)

type ID string

func (i ID) Str() string {
	return string(i)
}

func (i ID) GetListenURL() string {
	return env.GetAbsoluteURLForEnvironment(fmt.Sprintf("podcast/%s", id))
}

type dbUserPodcast struct {
	ID             ID                 `db:"_id"`
	CreatedAt      time.Time          `db:"created_at"`
	LastModifiedAt time.Time          `db:"last_modified_at"`
	EpisodeID      podcasts.EpisodeID `db:"episode_id"`
	SourceID       content.SourceID   `db:"source_id"`
	UserID         users.UserID       `db:"user_id"`
	EmailRecordID  email.ID           `db:"email_record_id"`
	FirstOpenedAt  *time.Time         `db:"first_opened_at"`
}

type UserPodcast struct {
	ID        ID                 `json:"id"`
	EpisodeID podcasts.EpisodeID `json:"episode_id"`
	UserID    users.UserID       `json:"user_id"`
	SourceID  content.SourceID   `json:"source_id"`
}

func (d dbUserPodcast) ToNonDB() UserPodcast {
	return UserPodcast{
		ID:        d.ID,
		EpisodeID: d.EpisodeID,
		SourceID:  d.SourceID,
		UserID:    d.UserID,
	}
}
