package htmlpages

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type HTMLPageID string

type HTMLPage struct {
	ID            HTMLPageID
	URL           string
	Language      *string
	Metadata      map[string]string
	OpenGraphType *string
}

type dbHTMLPage struct {
	ID            HTMLPageID `db:"_id"`
	URL           string     `db:"url"`
	Language      *string    `db:"language"`
	Metadata      dbMetadata `db:"metadata"`
	OpenGraphType *string    `db:"og_type"`
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
