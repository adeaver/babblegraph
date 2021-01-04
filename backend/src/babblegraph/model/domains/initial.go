package domains

import (
	"babblegraph/model/contenttopics"
	"babblegraph/util/urlparser"
)

var domainMap map[string]AllowableDomain

func initializeDomainMap() {
	for _, d := range allowableDomains {
		domainMap[string(d.Domain)] = d
	}
}

func IsURLAllowed(u urlparser.ParsedURL) bool {
	if len(domainMap) == 0 {
		initializeDomainMap()
	}
	_, ok := domainMap[u.Domain]
	return ok
}

func GetSeedDomains() map[string][]contenttopics.ContentTopic {
	out := make(map[string][]contenttopics.ContentTopic)
	for _, d := range allowableDomains {
		out[string(d.Domain)] = d.Topics
	}
	// TODO: add seed urls initial data
	return out
}
