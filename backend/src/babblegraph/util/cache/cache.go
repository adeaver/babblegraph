package cache

import (
	"babblegraph/util/database"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

func WithCache(key string, v interface{}, ttl time.Duration, fn func() (interface{}, error)) error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		item, err := lookupItemForCacheKey(tx, key)
		switch {
		case err != nil:
			return err
		case item != nil:
			return json.Unmarshal(item, v)
		case item == nil:
			newItem, err := fn()
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(newItem)
			if err != nil {
				return err
			}
			if err := insertItemIntoCache(tx, key, time.Now().Add(ttl), bytes); err != nil {
				return err
			}
			return json.Unmarshal(bytes, v)
		}
		return nil
	})
}
