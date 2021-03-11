package documents

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/domains"
	"babblegraph/util/elastic/esquery"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type DocumentWithScore struct {
	Document Document
	Score    decimal.Number
}

type IntRange struct {
	LowerBound *int64
	UpperBound *int64
}

type documentsQueryBuilder struct {
	languageCode      wordsmith.LanguageCode
	sentDocumentIDs   []string
	readingLevelRange *IntRange
	version           *IntRange

	// Each query can only search for a single topic
	// to make sure that the documents returned are most
	// relevant for that topic - not the union of the two
	topic  *contenttopics.ContentTopic
	lemmas []string
}

func NewDocumentsQueryBuilderForLanguage(languageCode wordsmith.LanguageCode) *documentsQueryBuilder {
	return &documentsQueryBuilder{languageCode: languageCode}
}

func (d *documentsQueryBuilder) ExecuteQuery() ([]DocumentWithScore, error) {
	queryBuilder := esquery.NewBoolQueryBuilder()
	queryBuilder.AddMust(esquery.Match("language_code", d.languageCode.Str()))
	queryBuilder.AddMust(esquery.Terms("domain.keyword", domains.GetDomains()))
	/*
	   We want to filter out articles that are assigned to all keywords
	   because they are irrelevant. This is necessary since those articles
	   may have many keywords from a user's tracking list.

	   Less than 5 was chosen by viewing the number of articles that have
	   multiple content topics attached to them. There is a huge spike at 7,
	   which is obviously not useful. By inspection, 4 seems reasonable.
	*/
	queryBuilder.AddFilter(esquery.Script("doc['content_topics.keyword'].size() < 5"))
	if filteredWords, ok := filteredWordsForLanguageCode[d.languageCode]; ok {
		queryBuilder.AddMustNot(esquery.Terms("metadata.title", filteredWords))
	}
	if len(d.sentDocumentIDs) != 0 {
		queryBuilder.AddMustNot(esquery.Terms("id.keyword", d.sentDocumentIDs))
	}
	if d.readingLevelRange != nil {
		readingLevelRangeQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("readability_score")
		if minReadingLevel := d.readingLevelRange.LowerBound; minReadingLevel != nil {
			readingLevelRangeQueryBuilder.GreaterThanOrEqualToInt64(*minReadingLevel)
		}
		if maxReadingLevel := d.readingLevelRange.UpperBound; maxReadingLevel != nil {
			readingLevelRangeQueryBuilder.LessThanOrEqualToInt64(*maxReadingLevel)
		}
		queryBuilder.AddMust(readingLevelRangeQueryBuilder.BuildRangeQuery())
	}
	if d.topic != nil {
		queryBuilder.AddMust(esquery.Match("content_topics", d.topic.Str()))
        // HACK HACK HACK
        // So there's an issue with lemma mappings
        // that can cause similar hyphenated strings
        // to interfere. I don't want to bother
        // with recreating the index, so I'm doing this instead.
        // The idea here is that we filter out documents that contain
        // similar keywords
	    queryBuilder.AddFilter(esquery.Term("content_topics.keyword", d.topic.Str())
	}
	if d.version != nil {
		versionRangeQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("version")
		if minVersion := d.version.LowerBound; minVersion != nil {
			versionRangeQueryBuilder.GreaterThanOrEqualToInt64(*minVersion)
		}
		if maxVersion := d.version.UpperBound; maxVersion != nil {
			versionRangeQueryBuilder.LessThanOrEqualToInt64(*maxVersion)
		}
		queryBuilder.AddMust(versionRangeQueryBuilder.BuildRangeQuery())
	}
	if len(d.lemmas) > 0 {
		queryBuilder.AddShould(esquery.Match("lemmatized_description", strings.Join(d.lemmas, " ")))
	}
	var docs []DocumentWithScore
	if err := esquery.ExecuteSearch(documentIndex{}, queryBuilder.BuildBoolQuery(), func(source []byte, score decimal.Number) error {
		log.Println(fmt.Sprintf("Document search got body %s", string(source)))
		var doc Document
		if err := json.Unmarshal(source, &doc); err != nil {
			return err
		}
		docs = append(docs, DocumentWithScore{
			Document: doc,
			Score:    score,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return docs, nil
}

func (d *documentsQueryBuilder) NotContainingDocuments(docIDs []DocumentID) {
	for _, doc := range docIDs {
		d.sentDocumentIDs = append(d.sentDocumentIDs, string(doc))
	}
}

func (d *documentsQueryBuilder) ForTopic(topic *contenttopics.ContentTopic) {
	d.topic = topic
}

func (d *documentsQueryBuilder) ForReadingLevelRange(lowerBound, upperBound *int64) {
	d.readingLevelRange = &IntRange{
		LowerBound: lowerBound,
		UpperBound: upperBound,
	}
}

func (d *documentsQueryBuilder) ForVersionRange(lowerBound, upperBound *Version) {
	var minVersion, maxVersion *int64
	if lowerBound != nil {
		minVersion = ptr.Int64(int64(*lowerBound))
	}
	if upperBound != nil {
		maxVersion = ptr.Int64(int64(*upperBound))
	}
	d.version = &IntRange{
		LowerBound: minVersion,
		UpperBound: maxVersion,
	}
}

func (d *documentsQueryBuilder) ContainingLemmas(lemmaIDs []wordsmith.LemmaID) {
	for _, l := range lemmaIDs {
		d.lemmas = append(d.lemmas, l.Str())
	}
}
