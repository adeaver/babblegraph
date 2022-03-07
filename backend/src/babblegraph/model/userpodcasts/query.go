package userpodcasts

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getByIDQuery = "SELECT * FROM user_podcasts WHERE _id = $1"
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
