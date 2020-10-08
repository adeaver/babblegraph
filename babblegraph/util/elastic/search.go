package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

func RunSearchRequest(req esapi.SearchRequest, fn func(sourceBytes []byte) error) error {
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var r searchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Got response from elasticsearch %+v", r))
	hits := r.Hits.Hits
	log.Println(fmt.Sprintf("Hits in response %+v", hits))
	for _, h := range hits {
		sourceBytes, err := json.Marshal(h.Source)
		if err != nil {
			return err
		}
		if err := fn(sourceBytes); err != nil {
			return err
		}
	}
	return nil
}
