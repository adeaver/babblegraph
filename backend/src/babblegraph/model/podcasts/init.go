package podcasts

import (
	"babblegraph/util/ctx"
	"babblegraph/util/elastic"
	"babblegraph/util/elastic/esmapping"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
)

const (
	analyzerWithWhitespaceTokenizerName = "analyzer_with_whitespace_tokenizer"
)

func CreatePodcastIndexes(c ctx.LogContext) {
	languageCodes := wordsmith.GetSupportedLanguageCodes()
	for _, code := range languageCodes {
		podcastIndex := getPodcastIndexForLanguageCode(code)
		if err := elastic.CreateIndex(podcastIndex, &elastic.CreateIndexSettings{
			Analysis: elastic.IndexAnalysis{
				Analyzer: elastic.IndexAnalyzer{
					Name: analyzerWithWhitespaceTokenizerName,
					Body: elastic.IndexAnalyzerBody{
						Type:      "custom",
						Tokenizer: elastic.AnalyzerTokenizerWhitespace,
					},
				},
			},
		}); err != nil {
			c.Errorf("Error creating index for language %s: %s", code, err.Error())
		}
	}
}

func CreatePodcastMappings(c ctx.LogContext) {
	languageCodes := wordsmith.GetSupportedLanguageCodes()
	for _, code := range languageCodes {
		podcastIndex := getPodcastIndexForLanguageCode(code)
		if err := esmapping.UpdateMapping(podcastIndex, []esmapping.Mapping{
			esmapping.MakeTextMapping("id", esmapping.MappingOptions{
				Analyzer: ptr.String(analyzerWithWhitespaceTokenizerName),
			}),
			esmapping.MakeTextMapping("title", esmapping.MappingOptions{}),
			esmapping.MakeTextMapping("description", esmapping.MappingOptions{}),
			esmapping.MakeDateMapping("publication_date", esmapping.MappingOptions{}),
			esmapping.MakeTextMapping("episode_type", esmapping.MappingOptions{}),
			esmapping.MakeLongMapping("duration", esmapping.MappingOptions{}),
			esmapping.MakeBooleanMapping("is_explicit", esmapping.MappingOptions{}),
			esmapping.MakeObjectMapping("audio_file", []esmapping.Mapping{
				esmapping.MakeTextMapping("url", esmapping.MappingOptions{
					Analyzer: ptr.String(analyzerWithWhitespaceTokenizerName),
				}),
				esmapping.MakeTextMapping("type", esmapping.MappingOptions{
					Analyzer: ptr.String(analyzerWithWhitespaceTokenizerName),
				}),
			}),
			esmapping.MakeLongMapping("version", esmapping.MappingOptions{}),
			esmapping.MakeTextMapping("source_id", esmapping.MappingOptions{
				Analyzer: ptr.String(analyzerWithWhitespaceTokenizerName),
			}),
			esmapping.MakeTextMapping("topic_ids", esmapping.MappingOptions{
				Analyzer: ptr.String(analyzerWithWhitespaceTokenizerName),
			}),
			esmapping.MakeTextMapping("language_code", esmapping.MappingOptions{}),
		}); err != nil {
			c.Errorf("Error updating mappings for podcast index %s: %s", code, err.Error())
		}
	}
}
