package routes

import "fmt"

const assetsBaseURL = "https://static.babblegraph.com/assets/"

func GetStaticAssetURLForResourceName(resourceName string) string {
	return fmt.Sprintf("%s%s", assetsBaseURL, resourceName)
}
