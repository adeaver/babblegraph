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
	expected := `{"document_name":{"analyzer":"my-cool-analyzer","enabled":true,"type":"text"}}`
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
	expected := `{"document":{"properties":{"author":{"properties":{"name":{"analyzer":"my-cool-analyzer","type":"text"}},"type":"object"},"name":{"type":"text"}},"type":"object"}}`
	if expected != string(mappingBytes) {
		t.Errorf("Expected text mapping %s but got %s", expected, string(mappingBytes))
	}
}

func TestMappingWithFields(t *testing.T) {
	withFields := MappingWithFields(MakeTextMapping("document_name", MappingOptions{}), []Mapping{
		MakeTextMapping("keyword", MappingOptions{}),
	})
	mappingBytes, err := json.Marshal(withFields)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := `{"document_name":{"fields":{"keyword":{"type":"text"}},"type":"text"}}`
	if expected != string(mappingBytes) {
		t.Errorf("Expected text mapping %s but got %s", expected, string(mappingBytes))
	}
}
