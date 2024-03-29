package esquery

import (
	"encoding/json"
	"testing"
)

func TestDescendingSort(t *testing.T) {
	s := NewDescendingSortBuilder("field")
	jsonBytes, err := json.Marshal(s.AsSort())
	if err != nil {
		t.Fatalf("Error on descending sort test: %s", err.Error())
	}
	expected := `{"field":{"order":"desc"}}`
	if string(jsonBytes) != expected {
		t.Errorf("Error on descending sort test. Expected %s, but got %s", expected, string(jsonBytes))
	}
}

func TestAscendingSort(t *testing.T) {
	s := NewAscendingSortBuilder("field")
	jsonBytes, err := json.Marshal(s.AsSort())
	if err != nil {
		t.Fatalf("Error on ascending sort test: %s", err.Error())
	}
	expected := `{"field":{"order":"asc"}}`
	if string(jsonBytes) != expected {
		t.Errorf("Error on ascending sort test. Expected %s, but got %s", expected, string(jsonBytes))
	}
}

func TestDescendingSortWithMissingValues(t *testing.T) {
	s := NewDescendingSortBuilder("field")
	s.WithMissingValuesLast()
	jsonBytes, err := json.Marshal(s.AsSort())
	if err != nil {
		t.Fatalf("Error on descending sort test: %s", err.Error())
	}
	expected := `{"field":{"order":"desc","missing":"_last"}}`
	if string(jsonBytes) != expected {
		t.Errorf("Error on descending sort test. Expected %s, but got %s", expected, string(jsonBytes))
	}
}

func TestAscendingSortWithMissingValues(t *testing.T) {
	s := NewAscendingSortBuilder("field")
	s.WithMissingValuesFirst()
	jsonBytes, err := json.Marshal(s.AsSort())
	if err != nil {
		t.Fatalf("Error on ascending sort test: %s", err.Error())
	}
	expected := `{"field":{"order":"asc","missing":"_first"}}`
	if string(jsonBytes) != expected {
		t.Errorf("Error on ascending sort test. Expected %s, but got %s", expected, string(jsonBytes))
	}
}

func TestDescendingSortWithLong(t *testing.T) {
	s := NewDescendingSortBuilder("field")
	s.AsUnmappedTypeLong()
	jsonBytes, err := json.Marshal(s.AsSort())
	if err != nil {
		t.Fatalf("Error on descending sort test: %s", err.Error())
	}
	expected := `{"field":{"order":"desc","unmapped_type":"long"}}`
	if string(jsonBytes) != expected {
		t.Errorf("Error on descending sort test. Expected %s, but got %s", expected, string(jsonBytes))
	}
}
