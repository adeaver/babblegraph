package elastic

import (
	"babblegraph/util/math/decimal"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/esapi"
)

type searchResponse struct {
	ScrollID string    `json:"_scroll_id"`
	Took     int64     `json:"took"`
	TimedOut bool      `json:"timed_out"`
	Shards   shardInfo `json:"_shards"`
	Hits     hitsInfo  `json:"hits"`
}

type shardInfo struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

type hitsInfo struct {
	Total    hitsTotal `json:"total"`
	MaxScore float64   `json:"max_score"`
	Hits     []hit     `json:"hits"`
}

type hit struct {
	Index  string      `json:"_index"`
	Type   *string     `json:"type,omitempty"`
	ID     string      `json:"_id"`
	Score  float64     `json:"_score"`
	Source interface{} `json:"_source"`
	Fields interface{} `json:"fields"`
}

type hitsTotal struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

func RunSearchRequest(req esapi.SearchRequest, fn func(sourceBytes []byte, relevance decimal.Number) error) error {
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var r searchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	hits := r.Hits.Hits
	for _, h := range hits {
		sourceBytes, err := json.Marshal(h.Source)
		if err != nil {
			return err
		}
		if err := fn(sourceBytes, decimal.FromFloat64(h.Score)); err != nil {
			return err
		}
	}
	return nil
}
