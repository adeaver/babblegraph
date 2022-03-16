package advertising

import (
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getUserAdvertisementQuery = "SELECT * FROM advertising_user_advertisements WHERE _id = $1"

	registerUserAdvertisementClickQuery = "INSERT INTO advertising_advertisement_clicks (user_advertisement_id) VALUES ($1)"
)

func GetURLForUserAdvertisementID(tx *sqlx.Tx, id UserAdvertisementID) (*string, error) {
	var matches []dbUserAdvertisement
	err := tx.Select(&matches, getUserAdvertisementQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) != 1:
		return nil, fmt.Errorf("Expected exactly one match for user advertisement with id %s, but got %d", id, len(matches))
	default:
		campaign, err := GetCampaignByID(tx, matches[0].CampaignID)
		if err != nil {
			return nil, err
		}
		return ptr.String(campaign.URL), nil
	}
}

func RegisterUserAdvertisementClick(tx *sqlx.Tx, id UserAdvertisementID) error {
	_, err := tx.Exec(registerUserAdvertisementClickQuery, id)
	return err
}
