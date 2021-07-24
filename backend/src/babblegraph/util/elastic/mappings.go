package elastic

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/getsentry/sentry-go"
)

func RunUpdateMappingsRequest(req esapi.IndicesPutMappingRequest) error {
	updateMappingsReq := req
	res, err := updateMappingsReq.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(res)
	migrationUpdateMappingsReq := req
	migrationRes, err := migrationUpdateMappingsReq.Do(context.Background(), migrationClient)
	if err != nil {
		sentry.CaptureException(fmt.Errorf("Caught error indexing for migration stack: %s", err.Error()))
		return nil
	}
	defer migrationRes.Body.Close()
	if migrationRes.StatusCode >= 300 {
		sentry.CaptureException(fmt.Errorf("Got status code %d for migration: %+v", migrationRes.StatusCode, migrationRes))
		return nil
	}
	log.Println(migrationRes)
	return nil
}
