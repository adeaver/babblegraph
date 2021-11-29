package documents

import (
	"babblegraph/util/elastic/esquery"
	"babblegraph/util/math/decimal"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"
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

/*
   We want to filter out articles that are assigned to all keywords
   because they are irrelevant. This is necessary since those articles
   may have many keywords from a user's tracking list.

   Less than 5 was chosen by viewing the number of articles that have
   multiple content topics attached to them. There is a huge spike at 7,
   which is obviously not useful. By inspection, 4 seems reasonable.
*/
const maximumNumberOfTopicsPerDocument int64 = 4

var validVersionsForLanguageCode = map[wordsmith.LanguageCode][]Version{
	wordsmith.LanguageCodeSpanish: {
		Version7,
		Version7,
	},
}

type executableQuery interface {
	ExtendBaseQuery(b *esquery.BoolQueryBuilder) error
}

type ExecuteDocumentQueryInput struct {
	LanguageCode        wordsmith.LanguageCode
	ExcludedDocumentIDs []DocumentID
	ValidDomains        []string
	MinimumReadingLevel *int64
	MaximumReadingLevel *int64
}

type DocumentWithScore struct {
	Document Document
	Score    decimal.Number
}

func ExecuteDocumentQuery(query executableQuery, input ExecuteDocumentQueryInput) ([]DocumentWithScore, error) {
	queryBuilder := esquery.NewBoolQueryBuilder()
	queryBuilder.AddMust(esquery.Match("language_code", input.LanguageCode.Str()))
	queryBuilder.AddMust(esquery.Terms("domain.keyword", input.ValidDomains))
	queryBuilder.AddMustNot(esquery.Match("has_paywall", true))
	if len(input.ExcludedDocumentIDs) != 0 {
		var excludedDocumentIDsQueryString []string
		for _, docID := range input.ExcludedDocumentIDs {
			excludedDocumentIDsQueryString = append(excludedDocumentIDsQueryString, string(docID))
		}
		queryBuilder.AddMustNot(esquery.Terms("id.keyword", excludedDocumentIDsQueryString))
	}
	topicsLengthRangeQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("topics_length")
	topicsLengthRangeQueryBuilder.LessThanOrEqualToInt64(maximumNumberOfTopicsPerDocument)
	queryBuilder.AddFilter(topicsLengthRangeQueryBuilder.BuildRangeQuery())
	if filteredWords, ok := filteredWordsForLanguageCode[input.LanguageCode]; ok {
		queryBuilder.AddMustNot(esquery.Terms("metadata.title", filteredWords))
	}
	if input.MinimumReadingLevel != nil || input.MaximumReadingLevel != nil {
		readingLevelRangeQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("readability_score")
		if input.MinimumReadingLevel != nil {
			readingLevelRangeQueryBuilder.GreaterThanOrEqualToInt64(*input.MinimumReadingLevel)
		}
		if input.MaximumReadingLevel != nil {
			readingLevelRangeQueryBuilder.LessThanOrEqualToInt64(*input.MaximumReadingLevel)
		}
		queryBuilder.AddMust(readingLevelRangeQueryBuilder.BuildRangeQuery())
	}
	versions, ok := validVersionsForLanguageCode[input.LanguageCode]
	if ok && len(versions) == 2 {
		versionRangeQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("version")
		versionRangeQueryBuilder.GreaterThanOrEqualToInt64(int64(versions[0]))
		versionRangeQueryBuilder.LessThanOrEqualToInt64(int64(versions[1]))
		queryBuilder.AddMust(versionRangeQueryBuilder.BuildRangeQuery())
	} else {
		log.Println(fmt.Sprintf("No valid document versions found for language code: %s", input.LanguageCode))
	}
	if err := query.ExtendBaseQuery(queryBuilder); err != nil {
		return nil, err
	}
	// This will sort by score and everything with the same score will be sorted by timestamp
	scoreSort := esquery.NewDescendingSortBuilder("_score").AsSort()
	timestampSortBuilder := esquery.NewDescendingSortBuilder("seed_job_ingest_timestamp")
	timestampSortBuilder.WithMissingValuesLast()
	timestampSortBuilder.AsUnmappedTypeLong()
	orderedSort := esquery.NewOrderedSort(scoreSort, timestampSortBuilder.AsSort())
	var docs []DocumentWithScore
	if err := esquery.ExecuteSearch(documentIndex{}, queryBuilder.BuildBoolQuery(), orderedSort, func(source []byte, score decimal.Number) error {
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

type UpdateDocumentInput struct {
	Version      Version `json:"version"`
	TopicsLength *int64  `json:"topics_length,omitempty"`
}

func UpdateDocumentForURL(u urlparser.ParsedURL, input UpdateDocumentInput) error {
	documentID := makeDocumentIndexForURL(u)
	return esquery.ExecuteUpdate(documentIndex{}, documentID.Str(), input)
}
