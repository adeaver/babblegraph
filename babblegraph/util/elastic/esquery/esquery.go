package esquery

import (
	"babblegraph/util/elastic"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

type queryName string

const (
	queryNameMatch queryName = "match"
	queryNameBool  queryName = "bool"
	queryNameRange queryName = "range"
)

func (q queryName) Str() string {
	return string(q)
}

type query map[string]interface{}

func makeQuery(key string, value interface{}) query {
	return query(map[string]interface{}{
		key: value,
	})
}

func ExecuteSearch(index elastic.Index, query query) (map[string]interface{}, error) {
	bodyBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	req := esapi.SearchRequest{
		Index: []string{index.GetName()},
		Body:  strings.NewReader(string(bodyBytes)),
	}
	return elastic.RunSearchRequest(req)
}
