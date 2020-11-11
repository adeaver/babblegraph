package linkprocessing

import (
	"babblegraph/util/urlparser"
	"strings"
)

var filteredPhrases = []string{
	"google",
	"facebook",
	"pinterest",
	"instagram",
	"amazon",
}

func shouldFilterOutURL(parsedURL urlparser.ParsedURL) bool {
	for _, phrase := range filteredPhrases {
		if strings.Contains(parsedURL.Domain, phrase) {
			return true
		}
	}
	return false
}
