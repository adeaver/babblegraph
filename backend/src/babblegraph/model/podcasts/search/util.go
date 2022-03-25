package search

import (
	"babblegraph/util/urlparser"
	"fmt"
)

// Listen Notes adds UTM parameters to the URL
// We don't want that.
func MaybeParseURLForListenNotesWebsiteURL(u string) *urlparser.ParsedURL {
	parsedWebsite := urlparser.ParseURL(u)
	if parsedWebsite == nil {
		return nil
	}
	websiteURL := fmt.Sprintf("%s/%s", parsedWebsite.Domain, parsedWebsite.Path)
	if parsedWebsite.Protocol != nil {
		websiteURL = fmt.Sprintf("%s//%s", *parsedWebsite.Protocol, websiteURL)
	}
	return urlparser.ParseURL(websiteURL)
}
