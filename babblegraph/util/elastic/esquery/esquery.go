package esquery

import (
	"babblegraph/util/elastic"
	"encoding/json"
	"fmt"
	"log"
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

type searchBody struct {
	Query query `json:"query"`
}

func ExecuteSearch(index elastic.Index, query query, fn func(source []byte) error) error {
	bodyBytes, err := json.Marshal(searchBody{
		Query: query,
	})
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Sending elasticsearch request %s", string(bodyBytes)))
	req := esapi.SearchRequest{
		Index: []string{index.GetName()},
		Body:  strings.NewReader(string(bodyBytes)),
	}
	return elastic.RunSearchRequest(req, fn)
}
