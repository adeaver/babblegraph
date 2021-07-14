package documents

import (
	"babblegraph/model/contenttopics"
	"babblegraph/util/elastic/esquery"
	"babblegraph/wordsmith"
	"strings"
)

// Each query can only search for a single topic
// to make sure that the documents returned are most
// relevant for that topic - not the union of the two
type dailyEmailDocumentsQueryBuilder struct {
	topic  *contenttopics.ContentTopic
	lemmas []string
}

func NewDailyEmailDocumentsQueryBuilder() *dailyEmailDocumentsQueryBuilder {
	return &dailyEmailDocumentsQueryBuilder{}
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
	if len(d.lemmas) > 0 {
		queryBuilder.AddShould(esquery.Match("lemmatized_description", strings.Join(d.lemmas, " ")))
	}
	return nil
}
