package documents

import (
	"babblegraph/model/contenttopics"
	"babblegraph/util/elastic/esquery"
	"babblegraph/wordsmith"
	"strings"
	"time"
)

type RecencyBias string

const (
	RecencyBiasMostRecent RecencyBias = "most_recent"
	RecencyBiasNotRecent  RecencyBias = "not_recent"

	RecencyBiasBoundary = -7 * 24 * time.Hour // one week
)

func (r RecencyBias) Ptr() *RecencyBias {
	return &r
}

// Each query can only search for a single topic
// to make sure that the documents returned are most
// relevant for that topic - not the union of the two
type dailyEmailDocumentsQueryBuilder struct {
	recencyBias *RecencyBias
	topic       *contenttopics.ContentTopic
	lemmas      []string
}

func NewDailyEmailDocumentsQueryBuilder() *dailyEmailDocumentsQueryBuilder {
	return &dailyEmailDocumentsQueryBuilder{}
}

func (d *dailyEmailDocumentsQueryBuilder) WithRecencyBias(r RecencyBias) {
	d.recencyBias = r.Ptr()
}

func (d *dailyEmailDocumentsQueryBuilder) ForTopic(topic *contenttopics.ContentTopic) {
	d.topic = topic
}

func (d *dailyEmailDocumentsQueryBuilder) ContainingLemmas(lemmaIDs []wordsmith.LemmaID) {
	for _, l := range lemmaIDs {
		d.lemmas = append(d.lemmas, l.Str())
	}
}

func (d *dailyEmailDocumentsQueryBuilder) ExtendBaseQuery(queryBuilder *esquery.BoolQueryBuilder) error {
	if d.topic != nil {
		queryBuilder.AddMust(esquery.Match("content_topics", d.topic.Str()))
		// HACK HACK HACK
		// So there's an issue with lemma mappings
		// that can cause similar hyphenated strings
		// to interfere. I don't want to bother
		// with recreating the index, so I'm doing this instead.
		// The idea here is that we filter out documents that contain
		// similar keywords
		queryBuilder.AddFilter(esquery.Term("content_topics.keyword", d.topic.Str()))
	}
	if d.recencyBias != nil {
		seedJobIngestTimestampRangeQuery := esquery.NewRangeQueryBuilderForFieldName("seed_job_ingest_timestamp")
		recencyBoundary := time.Now().Add(RecencyBiasBoundary).Unix()
		switch {
		case *d.recencyBias == RecencyBiasMostRecent:
			seedJobIngestTimestampRangeQuery.GreaterThanOrEqualToInt64(recencyBoundary)
			queryBuilder.AddMust(seedJobIngestTimestampRangeQuery.BuildRangeQuery())
		case *d.recencyBias == RecencyBiasNotRecent:
			// The minus one here is to ensure that documents don't appear twice
			seedJobIngestTimestampRangeQuery.GreaterThanOrEqualToInt64(recencyBoundary)
			queryBuilder.AddMustNot(seedJobIngestTimestampRangeQuery.BuildRangeQuery())
		}
	}
	if len(d.lemmas) > 0 {
		queryBuilder.AddShould(esquery.Match("lemmatized_description", strings.Join(d.lemmas, " ")))
	}
	return nil
}
