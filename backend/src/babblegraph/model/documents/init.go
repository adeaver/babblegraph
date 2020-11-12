package documents

import "babblegraph/util/elastic"

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
