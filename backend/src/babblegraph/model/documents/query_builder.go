package documents

import (
	"babblegraph/util/elastic/esquery"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"
)

type IntRange struct {
	LowerBound *int64
	UpperBound *int64
}

type documentsQueryBuilder struct {
	languageCode      wordsmith.LanguageCode
	sentDocumentIDs   []string
	readingLevelRange *IntRange
	version           *IntRange
}

func NewDocumentsQueryBuilderForLanguage(languageCode wordsmith.LanguageCode) *documentsQueryBuilder {
	return &documentsQueryBuilder{languageCode: languageCode}
}

func (d *documentsQueryBuilder) ExecuteQuery() ([]Document, error) {
	queryBuilder := esquery.NewBoolQueryBuilder()
	queryBuilder.AddMust(esquery.Match("language_code", d.languageCode.Str()))
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
	var docs []Document
	if err := esquery.ExecuteSearch(documentIndex{}, queryBuilder.BuildBoolQuery(), func(source []byte) error {
		log.Println(fmt.Sprintf("Document search got body %s", string(source)))
		var doc Document
		if err := json.Unmarshal(source, &doc); err != nil {
			return err
		}
		docs = append(docs, doc)
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
