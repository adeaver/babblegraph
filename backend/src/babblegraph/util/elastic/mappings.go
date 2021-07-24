package elastic

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
)

func RunUpdateMappingsRequest(req, migrationReq esapi.IndicesPutMappingRequest) error {
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(res)
	migrationRes, err := migrationReq.Do(context.Background(), migrationClient)
	if err != nil {
		handleMigrationError(fmt.Errorf("Caught error indexing for migration stack: %s", err.Error()))
		return nil
	}
	defer migrationRes.Body.Close()
	if migrationRes.StatusCode >= 300 {
		handleMigrationError(fmt.Errorf("Got status code %d for migration: %+v", migrationRes.StatusCode, migrationRes))
		return nil
	}
	log.Println(migrationRes)
	return nil
}
