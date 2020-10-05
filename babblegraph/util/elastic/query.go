package elastic

import (
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

type searchBody struct {
	Query interface{} `json:"query"`
}

func makeSearchRequest(index Index, body interface{}) (*esapi.SearchRequest, error) {
	body, err := json.Marshal(searchBody{
		Query: body,
	})
	if err != nil {
		return nil, err
	}
	return &esapi.SearchRequest{
		Index: []string{index.GetName()},
		Body:  strings.NewReader(string(body)),
	}, nil
}
