package documents

import (
	"babblegraph/util/elastic"
	"babblegraph/util/elastic/esmapping"
	"babblegraph/util/ptr"
)

func CreateDocumentIndex() error {
	return elastic.CreateIndex(documentIndex{}, &elastic.CreateIndexSettings{
		Analysis: elastic.IndexAnalysis{
			Analyzer: elastic.IndexAnalyzer{
				Name: "custom_analyzer",
				Body: elastic.IndexAnalyzerBody{
					Type:      "custom",
					Tokenizer: elastic.AnalyzerTokenizerWhitespace,
				},
			},
		},
	})
}

func makeDefaultTextWithKeywordField(fieldName string) esmapping.Mapping {
	return esmapping.MappingWithFields(
		esmapping.MakeTextMapping(fieldName, esmapping.MappingOptions{}),
		[]esmapping.Mapping{
			esmapping.MakeKeywordMapping("keyword", esmapping.MappingOptions{
				IgnoreAbove: ptr.Int64(256),
			}),
		},
	)
}

func CreateDocumentMappings() error {
	return esmapping.UpdateMapping(documentIndex{}, []esmapping.Mapping{
		makeDefaultTextWithKeywordField("content_topics"),
		makeDefaultTextWithKeywordField("document_type"),
		makeDefaultTextWithKeywordField("domain"),
		esmapping.MakeBooleanMapping("has_paywall", esmapping.MappingOptions{}),
		makeDefaultTextWithKeywordField("id"),
		makeDefaultTextWithKeywordField("language_code"),
		makeDefaultTextWithKeywordField("lemmatized_body"),
		makeDefaultTextWithKeywordField("lemmatized_description"),
		esmapping.MakeObjectMapping("metadata", []esmapping.Mapping{
			makeDefaultTextWithKeywordField("description"),
			makeDefaultTextWithKeywordField("image"),
			makeDefaultTextWithKeywordField("title"),
			makeDefaultTextWithKeywordField("url"),
			esmapping.MakeDateMapping("publication_time_utc", esmapping.MappingOptions{}),
		}),
		makeDefaultTextWithKeywordField("url"),
		esmapping.MakeLongMapping("lemmatized_description_index_mappings", esmapping.MappingOptions{}),
		esmapping.MakeLongMapping("readability_score", esmapping.MappingOptions{}),
		esmapping.MakeLongMapping("seed_job_ingest_timestamp", esmapping.MappingOptions{}),
		esmapping.MakeLongMapping("version", esmapping.MappingOptions{}),
	})
}
