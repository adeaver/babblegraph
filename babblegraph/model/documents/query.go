package documents

import (
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"
)

func FindDocumentsContainingTerms(terms []wordsmith.LemmaID) ([]Document, error) {
	res, err := elastic.InQuery{
		FieldName: "lemmatized_body",
		Values:    termsToString(terms),
	}.SearchIndex(documentIndex{})
	if err != nil {
		return nil, err
	}
	return extractDocuments(res)
}

// TODO: this should live in elastic package
func extractDocuments(res map[string]interface{}) ([]Document, error) {
	hits, ok := res["hits"]
	if !ok {
		log.Println("no hits")
		return nil, nil
	}
	hitResults, ok := hits["hits"]
	if !ok {
		log.Println("no hit results")
		return nil, nil
	}
	hitList, isList := hitResults.([]map[string]interface{})
	if !isList {
		return fmt.Errorf("results is not a list")
	}
	var out []Document
	for _, h := range hitList {
		source, ok := h["_source"]
		if !ok {
			continue
		}
		var doc Document
		if err := json.Unmarshal(&doc); err != nil {
			return nil, err
		}
		out = append(out, doc)
	}
	return out, nil
}

func termsToString(terms []wordsmith.LemmaID) []string {
	var out []string
	for _, t := range terms {
		out = append(out, string(t))
	}
	return out
}
