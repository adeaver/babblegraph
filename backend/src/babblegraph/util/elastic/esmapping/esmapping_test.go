package esmapping

import (
	"babblegraph/util/ptr"
	"encoding/json"
	"testing"
)

func TestMakeTextMapping(t *testing.T) {
	textMapping := MakeTextMapping("document_name", MappingOptions{
		Analyzer: ptr.String("my-cool-analyzer"),
		Enabled:  ptr.Bool(true),
	})
	mappingBytes, err := json.Marshal(textMapping)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := `{"document_name":{"type":"text","analyzer":"my-cool-analyzer","enabled":true}}`
	if expected != string(mappingBytes) {
		t.Errorf("Expected text mapping %s but got %s", expected, string(mappingBytes))
	}
}

func TestMakeObjectMapping(t *testing.T) {
	objectMapping := MakeObjectMapping("document", []Mapping{
		MakeTextMapping("name", MappingOptions{}),
		MakeObjectMapping("author", []Mapping{
			MakeTextMapping("name", MappingOptions{
				Analyzer: ptr.String("my-cool-analyzer"),
			}),
		}),
	})
	mappingBytes, err := json.Marshal(objectMapping)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := `{"document":{"type":"object","properties":{"name":{"type":"text"},"author":{"type":"object","properties":{"name":{"type":"text","analyzer":"my-cool-analyzer"}}}}}}`
	if expected != string(mappingBytes) {
		t.Errorf("Expected text mapping %s but got %s", expected, string(mappingBytes))
	}
}
