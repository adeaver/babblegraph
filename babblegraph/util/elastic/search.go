package elastic

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/esapi"
)

// TODO: turn response into a type
func RunSearchRequest(req esapi.SearchRequest) (map[string]interface{}, error) {
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
