package advertising

import "github.com/jmoiron/sqlx"

const (
	getAllVendorsQuery  = "SELECT * FROM advertising_vendor"
	getVendorByIDQuery  = "SELECT * FROM advertising_vendor WHERE _id = $1"
	insertVendorQuery   = "INSERT INTO advertising_vendor (name, website_url, is_active) VALUES ($1, $2, $3) RETURNING _id"
	editVendorByIDQuery = "UPDATE advertising_vendor SET name = $1, website_url = $2, is_active = $3 WEHRE _id = $4"
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
