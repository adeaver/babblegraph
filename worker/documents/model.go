package documents

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DocumentID string

type dbDocument struct {
	ID       DocumentID `db:"_id"`
	URL      string     `db:"url"`
	Language *string    `db:"language"`
	Metadata dbMetadata `db:"metadata"`
}

type dbMetadata map[string]string

func (d dbMetadata) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *dbMetadata) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Could not scan into byte array")
	}
	var i interface{}
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}
	*d, ok = i.(map[string]string)
	if !ok {
		return fmt.Errorf("Could not convert type")
	}
	return nil
}

type Document struct {
	ID       DocumentID        `db:"_id"`
	URL      string            `db:"url"`
	Language *string           `db:"language"`
	Metadata map[string]string `db:"metadata"`
}
