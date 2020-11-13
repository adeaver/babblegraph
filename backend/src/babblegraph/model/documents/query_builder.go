package documents

import (
	"babblegraph/util/elastic/esquery"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"encoding/json"
)

type IntRange struct {
	LowerBound *int64
	UpperBound *int64
}

type documentsQueryBuilder struct {
	language          string
	sentDocumentIDs   []string
	readingLevelRange *IntRange
	version           *IntRange
}

func NewDocumentsQueryBuilderForLanguage(languageCode wordsmith.LanguageCode) *documentsQueryBuilder {
	return &documentsQueryBuilder{language: languageCode.Str()}
}

func (d *documentsQueryBuilder) ExecuteQuery() ([]Document, error) {
	queryBuilder := esquery.NewBoolQueryBuilder()
	queryBuilder.AddMust(esquery.Match("language_code", d.language))
	if len(d.sentDocumentIDs) != 0 {
		queryBuilder.AddMustNot(esquery.Terms("id", d.sentDocumentIDs))
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
	}
	var docs []Document
	if err := esquery.ExecuteSearch(documentIndex{}, queryBuilder.BuildBoolQuery(), func(source []byte) error {
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
