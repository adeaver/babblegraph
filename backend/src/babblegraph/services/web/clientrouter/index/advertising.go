package index

import (
	"babblegraph/config"
	"babblegraph/model/advertising"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

func handleAdClick(r *router.Request) (interface{}, error) {
	var useAdditionalLink bool
	if additionalLinkParam := r.GetQueryParam(config.AdvertisingUseAdditionalLinkQueryParamName); additionalLinkParam != nil && *additionalLinkParam == config.AdvertisingUseAdditionalLinkQueryParamValue {
		useAdditionalLink = true
	}
	userAdvertisementIDStr, err := r.GetRouteVar("token")
	if err != nil {
		return nil, err
	}
	userAdvertisementID := advertising.UserAdvertisementID(*userAdvertisementIDStr)
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return advertising.RegisterUserAdvertisementClick(tx, userAdvertisementID)
	}); err != nil {
		r.Warnf("Error inserting ad click for ID %s: %s", userAdvertisementID, err.Error())
	}
	var url *string
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		url, err = advertising.GetURLForUserAdvertisementID(tx, userAdvertisementID, useAdditionalLink)
		return err
	}); err != nil {
		return nil, err
	}
	return urlparser.EnsureProtocol(*url)
}
