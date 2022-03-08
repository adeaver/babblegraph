package documents

import (
	"babblegraph/model/content"
	"babblegraph/util/elastic/esquery"
	"babblegraph/wordsmith"
	"time"
)

type lemmaSpotlightQueryBuilder struct {
	lemmaID     wordsmith.LemmaID
	topics      []content.TopicID
	recencyBias *RecencyBias
}

func NewLemmaSpotlightQueryBuilder(lemmaID wordsmith.LemmaID) *lemmaSpotlightQueryBuilder {
	return &lemmaSpotlightQueryBuilder{
		lemmaID: lemmaID,
	}
}

func (l *lemmaSpotlightQueryBuilder) AddTopics(topics []content.TopicID) {
	l.topics = append(l.topics, topics...)
}

func (l *lemmaSpotlightQueryBuilder) WithRecencyBias(r RecencyBias) {
	l.recencyBias = r.Ptr()
}

func (l *lemmaSpotlightQueryBuilder) ExtendBaseQuery(queryBuilder *esquery.BoolQueryBuilder) error {
	// Note, this is a bit of a hack. We're using match phrase because
	// of a bug with how ElasticSearch is setup. Currently, the analyzer
	// splits hyphens, so match phrase will guarantee that our lemma is matched exactly.
	queryBuilder.AddFilter(esquery.MatchPhrase("lemmatized_description", string(l.lemmaID)))
	if l.recencyBias != nil {
		seedJobIngestTimestampRangeQuery := esquery.NewRangeQueryBuilderForFieldName("seed_job_ingest_timestamp")
		recencyBoundary := time.Now().Add(RecencyBiasBoundary).Unix()
		switch {
		case *l.recencyBias == RecencyBiasMostRecent:
			seedJobIngestTimestampRangeQuery.GreaterThanOrEqualToInt64(recencyBoundary)
			queryBuilder.AddMust(seedJobIngestTimestampRangeQuery.BuildRangeQuery())
		case *l.recencyBias == RecencyBiasNotRecent:
			// The minus one here is to ensure that documents don't appear twice
			seedJobIngestTimestampRangeQuery.GreaterThanOrEqualToInt64(recencyBoundary)
			queryBuilder.AddMustNot(seedJobIngestTimestampRangeQuery.BuildRangeQuery())
		}
	}
	if len(l.topics) > 0 {
		var topicsQueryString []string
		for _, t := range l.topics {
			topicsQueryString = append(topicsQueryString, t.Str())
		}
		queryBuilder.AddShould(esquery.Terms("topic_ids", topicsQueryString))
	}
	return nil
}
