package elastic

import (
	"encoding/json"
	"testing"
)

func TestSerialization(t *testing.T) {
	testSettings := CreateIndexSettings{
		Analysis: IndexAnalysis{
			Analyzer: IndexAnalyzer{
				Name: "custom_analyzer",
				Body: IndexAnalyzerBody{
					Type:      "custom",
					Tokenizer: AnalyzerTokenizerWhitespace,
				},
			},
		},
	}
	out, err := json.Marshal(&testSettings)
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := `{"analysis":{"analyzer":{"custom_analyzer":{"type":"custom","tokenizer":"whitespace"}}}}`
	if string(out) != expected {
		t.Errorf("expcted %s, but got %s", expected, string(out))
	}
}
