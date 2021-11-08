package elastic

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
)

func RunUpdateRequest(req esapi.UpdateRequest) error {
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(res)
	return nil
}
