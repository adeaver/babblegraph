package domains

import (
	"babblegraph/model/contenttopics"
	"babblegraph/util/geo"
	"babblegraph/wordsmith"
)

type Domain string

type AllowableDomain struct {
	Domain       Domain
	Country      geo.CountryCode
	LanguageCode wordsmith.LanguageCode

	// This is to be used if an entire domain maps to a specific topic.
	// i.e. Motortrend magazine is all about cars
	// If a domain has multiple topics, this can and should be empty
	Topics []contenttopics.ContentTopic
}

// A seed url is a URL from which to start pulling content
// It can be a url on an allowable domain. This might be something like
// elmundo.es/deportes for the sports page.
type SeedURL struct {
	URL    string
	Topics []contenttopics.ContentTopic
}
