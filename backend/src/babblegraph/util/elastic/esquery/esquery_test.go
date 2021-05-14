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

func TestTerm(t *testing.T) {
	testQuery := Term("text", "abc 123")
	expected := `{"term":{"text":"abc 123"}}`
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
	expected := `{"bool":{"must":[{"match":{"text":"abc 123"}}],"filter":[{"match":{"text":"abc 123"}}]}}`
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

func TestTerms(t *testing.T) {
	testQuery := Terms("_id", []string{"abc", "123"})
	expected := `{"terms":{"_id":["abc","123"]}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}

func TestMatchAll(t *testing.T) {
	testQuery := MatchAll()
	expected := `{"match_all":{}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}

func TestScript(t *testing.T) {
	testQuery := Script("doc['content_topics.keyword'].size() == 2")
	expected := `{"script":{"script":"doc['content_topics.keyword'].size() == 2"}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}

func TestQueryWithSort(t *testing.T) {
	ascendingSort := NewAscendingSortBuilder("field").AsSort()
	testQuery := searchBody{
		Query: MatchAll(),
		Sort: []sort{
			ascendingSort,
		},
	}
	expectedBody := `{"query":{"match_all":{}},"sort":[{"field":{"order":"asc","missing":"_last"}}]}`
	jsonBytes, err := json.Marshal(testQuery)
	if err != nil {
		t.Fatalf("Error on testing query with sort: %s", err.Error())
	}
	if string(jsonBytes) != expectedBody {
		t.Errorf("Error on testing query with sort: expected %s, but got %s", expectedBody, string(jsonBytes))
	}
}

func TestQueryWithoutSort(t *testing.T) {
	testQuery := searchBody{
		Query: MatchAll(),
		Sort:  []sort{},
	}
	expectedBody := `{"query":{"match_all":{}}}`
	jsonBytes, err := json.Marshal(testQuery)
	if err != nil {
		t.Fatalf("Error on testing query without sort: %s", err.Error())
	}
	if string(jsonBytes) != expectedBody {
		t.Errorf("Error on testing query without sort: expected %s, but got %s", expectedBody, string(jsonBytes))
	}
}
