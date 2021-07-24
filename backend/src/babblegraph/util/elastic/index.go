package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

type Index interface {
	GetName() string
	ValidateDocument(document interface{}) error
	GenerateIDForDocument(document interface{}) (*string, error)
}

type CreateIndexSettings struct {
	Analysis IndexAnalysis `json:"analysis"`
}

type settingsBody struct {
	Settings CreateIndexSettings `json:"settings"`
}

func CreateIndex(index Index, settings *CreateIndexSettings) error {
	createIndexRequest := esapi.IndicesCreateRequest{
		Index: index.GetName(),
	}
	if settings != nil {
		bodyBytes, err := json.Marshal(&settingsBody{Settings: *settings})
		if err != nil {
			return err
		}
		createIndexRequest.Body = strings.NewReader(string(bodyBytes))
	}
	res, err := createIndexRequest.Do(context.Background(), esClient)
	if err != nil {
		log.Println(fmt.Sprintf("Caught error creating: %s", err.Error()))
		return nil
	}
	defer res.Body.Close()
	log.Println(res)
	migrationRes, err := createIndexRequest.Do(context.Background(), migrationClient)
	if err != nil {
		handleMigrationError(fmt.Errorf("Caught error creating index for migration stack: %s", err.Error()))
		return nil
	}
	defer migrationRes.Body.Close()
	log.Println(migrationRes)
	return nil
}

func IndexDocument(index Index, document interface{}) error {
	if err := index.ValidateDocument(document); err != nil {
		return fmt.Errorf("Document validation error for index %s: %s", index.GetName(), err.Error())
	}
	docID, err := index.GenerateIDForDocument(document)
	if err != nil {
		return err
	}
	documentAsJSON, err := json.Marshal(&document)
	if err != nil {
		return fmt.Errorf("Marshalling error for document %+v: %s", document, err.Error())
	}
	indexRequest := esapi.IndexRequest{
		Index:      index.GetName(),
		Body:       strings.NewReader(string(documentAsJSON)),
		DocumentID: *docID,
		Refresh:    "true",
	}
	res, err := indexRequest.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(res)
	indexRequest = esapi.IndexRequest{
		Index:      index.GetName(),
		Body:       strings.NewReader(string(documentAsJSON)),
		DocumentID: *docID,
		Refresh:    "true",
	}
	migrationRes, err := indexRequest.Do(context.Background(), migrationClient)
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
