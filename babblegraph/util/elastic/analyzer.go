package elastic

import (
	"encoding/json"
	"fmt"
)

type AnalyzerTokenizer string
type AnalyzerCharacterFilter string
type AnalyzerFilter string

const (
	AnalyzerTokenizerStandard   AnalyzerTokenizer = "standard"
	AnalyzerTokenizerWhitespace AnalyzerTokenizer = "whitespace"
)

type IndexAnalysis struct {
	Analyer IndexAnalyzer `json:"analyzer"`
}

type IndexAnalyzer struct {
	Body IndexAnalyzerBody
	Name string
}

type IndexAnalyzerBody struct {
	Type            string                    `json:"type"`
	CharacterFilter []AnalyzerCharacterFilter `json:"char_filter"`
	Tokenizer       AnalyzerTokenizer         `json:"tokenizer"`
	Filter          []AnalyzerFilter          `json:"filter"`
}

func (i IndexAnalyzer) MarshalJSON() ([]byte, error) {
	analyzerJSONMap := make(map[string]IndexAnalyzerBody)
	analyzerJSONMap[i.Name] = i.Body
	return json.Marshal(analyzerJSONMap)
}

func (i *IndexAnalyzer) UnmarshalJSON(data []byte) error {
	var analyzerMap map[string]IndexAnalyzerBody
	if err := json.Unmarshal(data, &analyzerMap); err != nil {
		return err
	}
	count := 0
	for name, body := range analyzerMap {
		if count > 0 {
			return fmt.Errorf("Got multiple analyzers")
		}
		i.Name = name
		i.Body = body
		count++
	}
	return nil
}
