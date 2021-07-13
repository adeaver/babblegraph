package esquery

import (
	"babblegraph/util/elastic"
	"babblegraph/util/math/decimal"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

type queryName string

const (
	queryNameMatch       queryName = "match"
	queryNameMatchPhrase queryName = "match_phrase"
	queryNameBool        queryName = "bool"
	queryNameRange       queryName = "range"
	queryNameTerms       queryName = "terms"
	queryNameTerm        queryName = "term"
	queryNameIDs         queryName = "ids"
	queryNameScript      queryName = "script"
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
	Query query  `json:"query"`
	Sort  []sort `json:"sort,omitempty"`
}

func ExecuteSearch(index elastic.Index, query query, orderedSort *orderedSort, fn func(source []byte, relevance decimal.Number) error) error {
	var sorts []sort
	if orderedSort != nil {
		sorts = orderedSort.sorts
	}
	bodyBytes, err := json.Marshal(searchBody{
		Query: query,
		Sort:  sorts,
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
