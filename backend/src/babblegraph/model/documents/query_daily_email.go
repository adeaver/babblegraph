package documents

import (
	"babblegraph/model/content"
	"babblegraph/util/elastic/esquery"
	"babblegraph/wordsmith"
	"strings"
	"time"
)

// Each query can only search for a single topic
// to make sure that the documents returned are most
// relevant for that topic - not the union of the two
type dailyEmailDocumentsQueryBuilder struct {
	recencyBias *RecencyBias
	topic       *content.TopicID
	lemmas      []string
}

func NewDailyEmailDocumentsQueryBuilder() *dailyEmailDocumentsQueryBuilder {
	return &dailyEmailDocumentsQueryBuilder{}
}

func (d *dailyEmailDocumentsQueryBuilder) WithRecencyBias(r RecencyBias) {
	d.recencyBias = r.Ptr()
}

func (d *dailyEmailDocumentsQueryBuilder) ForTopic(topic *content.TopicID) {
	d.topic = topic
}

func (d *dailyEmailDocumentsQueryBuilder) ContainingLemmas(lemmaIDs []wordsmith.LemmaID) {
	for _, l := range lemmaIDs {
		d.lemmas = append(d.lemmas, l.Str())
	}
}

func (d *dailyEmailDocumentsQueryBuilder) ExtendBaseQuery(queryBuilder *esquery.BoolQueryBuilder) error {
	if d.topic != nil {
		queryBuilder.AddMust(esquery.MatchPhrase("topic_ids", d.topic.Str()))
		// HACK HACK HACK
		// So there's an issue with lemma mappings
		// that can cause similar hyphenated strings
		// to interfere. I don't want to bother
		// with recreating the index, so I'm doing this instead.
		// The idea here is that we filter out documents that contain
		// similar keywords
		queryBuilder.AddFilter(esquery.Term("topic_ids.keyword", d.topic.Str()))
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
