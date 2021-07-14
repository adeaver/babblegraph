package routes

import (
	"babblegraph/util/geo"
	"fmt"
)

const assetsBaseURL = "https://static.babblegraph.com/assets/"

func GetStaticAssetURLForResourceName(resourceName string) string {
	return fmt.Sprintf("%s%s", assetsBaseURL, resourceName)
}

func GetFlagAssetForCountryCode(code geo.CountryCode) string {
	return GetStaticAssetURLForResourceName(fmt.Sprintf("geoflags/%s.png", string(code)))
}
