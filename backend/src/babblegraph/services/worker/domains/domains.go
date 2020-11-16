package domains

import "babblegraph/util/urlparser"

func IsURLAllowed(u urlparser.ParsedURL) bool {
	_, ok := allowableDomains[u.Domain]
	return ok
}

func GetSeedURLs() []string {
	var out []string
	for url, _ := range allowableDomains {
		out = append(out, url)
	}
	return out
}
