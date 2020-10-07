package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type QueryAction string

const (
	QueryActionTerm  QueryAction = "term"
	QueryActionTerms QueryAction = "terms"
	QueryActionMatch QueryAction = "match"
)

func (q QueryAction) Str() string {
	return string(q)
}

type InQuery struct {
	FieldName string
	Values    []string
}

func (i InQuery) SearchIndex(index Index) (map[string]interface{}, error) {
	req, err := makeSearchRequest(index, i)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Request: %+v", req))
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

func (i InQuery) MarshalJSON() ([]byte, error) {
	bodyMap := make(map[string]string)
	bodyMap[i.FieldName] = strings.Join(i.Values, " ")
	queryMap := make(map[QueryAction]map[string]string)
	queryMap[QueryActionMatch] = bodyMap
	return json.Marshal(queryMap)
}

func (i *InQuery) UnmarshalJSON(data []byte) error {
	var queryMap map[string]map[string][]string
	if err := json.Unmarshal(data, &queryMap); err != nil {
		return err
	}
	bodyMap, ok := queryMap[QueryActionTerms.Str()]
	if !ok {
		return fmt.Errorf("malformatted terms query")
	}
	count := 0
	for fieldName, values := range bodyMap {
		if count > 0 {
			return fmt.Errorf("malformatted terms query: too many fields")
		}
		i.FieldName = fieldName
		i.Values = values
		count++
	}
	return nil

}
