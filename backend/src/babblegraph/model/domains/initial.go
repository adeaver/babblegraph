package domains

import (
	"babblegraph/model/contenttopics"
	"babblegraph/util/urlparser"
	"fmt"
	"sync"
)

var (
	domainMap        map[string]AllowableDomain
	initializerMutex sync.Mutex
)

func initializeDomainMap() {
	initializerMutex.Lock()
	defer initializerMutex.Unlock()
	if len(domainMap) != 0 {
		return
	}
	domainMap = make(map[string]AllowableDomain)
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

func GetDomainMetadata(d string) (*AllowableDomain, error) {
	if len(domainMap) == 0 {
		initializeDomainMap()
	}
	metadata, ok := domainMap[d]
	if !ok {
		return nil, fmt.Errorf("Invalid domain")
	}
	return &metadata, nil
}

func GetDomains() []string {
	var out []string
	for _, d := range allowableDomains {
		out = append(out, string(d.Domain))
	}
	return out
}

func GetSeedURLs() map[string][]contenttopics.ContentTopic {
	out := make(map[string][]contenttopics.ContentTopic)
	for _, d := range allowableDomains {
		out[string(d.Domain)] = d.Topics
	}
	for _, d := range seedURLs {
		out[d.URL] = d.Topics
	}
	return out
}
