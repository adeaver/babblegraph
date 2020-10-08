package esquery

import (
	"encoding/json"
	"testing"
)

func TestMatch(t *testing.T) {
	testQuery := Match("text", "abc 123")
	expected := `{"match":{"text":"abc 123"}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}

func TestBool(t *testing.T) {
	builder := NewBoolQueryBuilder()
	builder.AddMust(Match("text", "abc 123"))
	builder.AddFilter(Match("text", "abc 123"))
	testQuery := builder.BuildBoolQuery()
	expected := `{"bool":{"must":[{"match":{"text":"abc 123"}}],"filters":[{"match":{"text":"abc 123"}}]}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}

func TestRange(t *testing.T) {
	builder := NewRangeQueryBuilderForFieldName("test")
	builder.LessThanInt64(50)
	builder.GreaterThanInt64(10)
	testQuery := builder.BuildRangeQuery()
	expected := `{"range":{"test":{"gt":10,"lt":50}}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}
