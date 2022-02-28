package cache

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	lookupFromCacheQuery = "SELECT * FROM cache WHERE key = $1"
	deleteFromCacheQuery = "DELETE FROM cache WHERE key = $1"
	insertIntoCacheQuery = "INSERT INTO cache (key, expires_at, item) VALUES ($1, $2, $3)"
)

type dbCachedItem struct {
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	Key       string    `db:"key"`
	Item      []byte    `db:"item"`
}

func lookupItemForCacheKey(tx *sqlx.Tx, cacheKey string) ([]byte, error) {
	var matches []dbCachedItem
	err := tx.Select(&matches, lookupFromCacheQuery, cacheKey)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one result for cache key %s, but got %d", cacheKey, len(matches))
	case len(matches) == 1:
		m := matches[0]
		if time.Now().After(m.ExpiresAt) {
			if _, err := tx.Exec(deleteFromCacheQuery, cacheKey); err != nil {
				return nil, err
			}
			return nil, nil
		}
		return m.Item, nil
	default:
		panic("unreachable")
	}
}

func insertItemIntoCache(tx *sqlx.Tx, cacheKey string, expirationTime time.Time, item []byte) error {
	if _, err := tx.Exec(cacheKey, insertIntoCacheQuery, expirationTime, item); err != nil {
		return err
	}
	return nil
}
