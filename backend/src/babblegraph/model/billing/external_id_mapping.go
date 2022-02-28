package billing

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getExternalIDMappingByIDQuery = "SELECT * FROM billing_external_id_mapping WHERE _id = $1"
	insertExternalIDMappingQuery  = "INSERT INTO billing_external_id_mapping (id_type, external_id) VALUES ($1, $2) RETURNING _id"

	getExternalIDMappingByExternalIDQuery = "SELECT * FROM billing_external_id_mapping WHERE id_type = $1 AND external_id = $2"
)

func getExternalIDMapping(tx *sqlx.Tx, id externalIDMappingID) (*dbExternalIDMapping, error) {
	var matches []dbExternalIDMapping
	err := tx.Select(&matches, getExternalIDMappingByIDQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected exactly one external id mapping for id %s, but got %d", id, len(matches))
	default:
		return &matches[0], nil
	}
}

func lookupExternalIDMappingByExternalID(tx *sqlx.Tx, idType externalIDType, externalID string) (*dbExternalIDMapping, error) {
	var matches []dbExternalIDMapping
	err := tx.Select(&matches, getExternalIDMappingByExternalIDQuery, idType, externalID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one external id mapping for id %s (type %s), but got %d", externalID, idType, len(matches))
	default:
		return &matches[0], nil
	}
}

func insertExternalIDMapping(tx *sqlx.Tx, externalID string) (*externalIDMappingID, error) {
	rows, err := tx.Query(insertExternalIDMappingQuery, externalIDTypeStripe, externalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id externalIDMappingID
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
	}
	return &id, nil
}
