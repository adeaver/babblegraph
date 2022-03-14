package advertising

import (
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

const (
	getAllVendorsQuery  = "SELECT * FROM advertising_vendors"
	insertVendorQuery   = "INSERT INTO advertising_vendors (name, website_url, is_active) VALUES ($1, $2, $3) RETURNING _id"
	editVendorByIDQuery = "UPDATE advertising_vendors SET name = $1, website_url = $2, is_active = $3, last_modified_at = timezone('utc', now()) WHERE _id = $4"

	getAllSourcesQuery  = "SELECT * FROM advertising_sources"
	insertSourceQuery   = "INSERT INTO advertising_sources (name, url, type, is_active) VALUES ($1, $2, $3, $4)"
	editSourceByIDQuery = "UPDATE advertising_sources SET name = $1, url = $2, type = $3, is_active = $4, last_modified_at = timezone('utc', now()) WHERE _id = $5"
)

func GetAllVendors(tx *sqlx.Tx) ([]Vendor, error) {
	var matches []dbVendor
	if err := tx.Select(&matches, getAllVendorsQuery); err != nil {
		return nil, err
	}
	var out []Vendor
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func InsertVendor(tx *sqlx.Tx, name string, u urlparser.ParsedURL, isActive bool) (*VendorID, error) {
	rows, err := tx.Queryx(insertVendorQuery, name, u.URL, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vendorID VendorID
	for rows.Next() {
		if err := rows.Scan(&vendorID); err != nil {
			return nil, err
		}
	}
	return &vendorID, nil
}

func EditVendor(tx *sqlx.Tx, id VendorID, name string, u urlparser.ParsedURL, isActive bool) error {
	if _, err := tx.Exec(editVendorByIDQuery, name, u.URL, isActive, id); err != nil {
		return err
	}
	return nil
}

func GetAllSources(tx *sqlx.Tx) ([]AdvertisementSource, error) {
	var matches []dbAdvertisementSource
	if err := tx.Select(&matches, getAllSourcesQuery); err != nil {
		return nil, err
	}
	var out []AdvertisementSource
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func InsertSource(tx *sqlx.Tx, name string, u urlparser.ParsedURL, sourceType AdvertisementSourceType, isActive bool) (*AdvertisementSourceID, error) {
	rows, err := tx.Queryx(insertSourceQuery, name, u.URL, sourceType, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sourceID AdvertisementSourceID
	for rows.Next() {
		if err := rows.Scan(&sourceID); err != nil {
			return nil, err
		}
	}
	return &sourceID, nil
}

type EditSourceInput struct {
	Name       string
	URL        urlparser.ParsedURL
	SourceType AdvertisementSourceType
	IsActive   bool
}

func EditSource(tx *sqlx.Tx, id AdvertisementSourceID, input EditSourceInput) error {
	if _, err := tx.Exec(editSourceByIDQuery, input.Name, input.URL.URL, input.SourceType, input.IsActive, id); err != nil {
		return err
	}
	return nil
}
