package documents

import (
	"babblegraph/model/contenttopics"
	"babblegraph/util/elastic/esquery"
	"babblegraph/wordsmith"
)

type lemmaSpotlightQueryBuilder struct {
	lemmaID wordsmith.LemmaID
	topics  []contenttopics.ContentTopic
}

func NewLemmaSpotlightQueryBuilder(lemmaID wordsmith.LemmaID) *lemmaSpotlightQueryBuilder {
	return &lemmaSpotlightQueryBuilder{
		lemmaID: lemmaID,
	}
}

func (l *lemmaSpotlightQueryBuilder) AddTopics(topics []contenttopics.ContentTopic) {
	l.topics = append(l.topics, topics...)
}

func (l *lemmaSpotlightQueryBuilder) ExtendBaseQuery(queryBuilder *esquery.BoolQueryBuilder) error {
	// Note, this is a bit of a hack. We're using match phrase because
	// of a bug with how ElasticSearch is setup. Currently, the analyzer
	// splits hyphens, so match phrase will guarantee that our lemma is matched exactly.
	queryBuilder.AddMust(esquery.MatchPhrase("lemmatized_description", string(l.lemmaID)))
	if len(l.topics) > 0 {
		var topicsQueryString []string
		for _, t := range l.topics {
			topicsQueryString = append(topicsQueryString, string(t))
		}
		queryBuilder.AddShould(esquery.Terms("content_topics", topicsQueryString))
	}
	return nil
}
