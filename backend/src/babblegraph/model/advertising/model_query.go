package advertising

import (
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

const (
	getAllVendorsQuery  = "SELECT * FROM advertising_vendors"
	insertVendorQuery   = "INSERT INTO advertising_vendors (name, website_url, is_active) VALUES ($1, $2, $3) RETURNING _id"
	editVendorByIDQuery = "UPDATE advertising_vendors SET name = $1, website_url = $2, is_active = $3 WHERE _id = $4"
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
