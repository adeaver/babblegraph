package userpodcasts

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/users"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getByIDQuery     = "SELECT * FROM user_podcasts WHERE _id = $1"
	getByUserIDQuery = "SELECT * FROM user_podcasts WHERE user_id = $1"

	insertUserPodcastsQuery = "INSERT INTO user_podcasts (user_id, episode_id, source_id, email_record_id) VALUES ($1, $2, $3, $4) RETURNING _id"
)

func GetByID(tx *sqlx.Tx, id ID) (*UserPodcast, error) {
	var matches []dbUserPodcast
	err := tx.Select(&matches, getByIDQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected exactly one user podcast for ID %s, but got %d", id, len(matches))
	default:
		out := matches[0].ToNonDB()
		return &out, nil
	}
}

func GetByUserID(tx *sqlx.Tx, userID users.UserID) ([]UserPodcast, error) {
	var matches []dbUserPodcast
	if err := tx.Select(&matches, getByUserIDQuery, userID); err != nil {
		return nil, err
	}
	var out []UserPodcast
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func InsertUserPodcastAndReturnID(tx *sqlx.Tx, episodeID podcasts.EpisodeID, userID users.UserID, sourceID content.SourceID, emailRecordID email.ID) (*ID, error) {
	rows, err := tx.Query(insertUserPodcastsQuery, userID, episodeID, sourceID, emailRecordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id ID
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
	}
	return &id, nil
}

func RegisterOpenedPodcast(tx *sqlx.Tx, id ID) error {
	return nil
}
