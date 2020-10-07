package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

type searchBody struct {
	Query interface{} `json:"query"`
}

func makeSearchRequest(index Index, body interface{}) (*esapi.SearchRequest, error) {
	bodyBytes, err := json.Marshal(searchBody{
		Query: body,
	})
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Body: %s", string(bodyBytes)))
	return &esapi.SearchRequest{
		Index: []string{index.GetName()},
		Body:  strings.NewReader(string(bodyBytes)),
	}, nil
}
