package documents

import (
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
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
	log.Println(fmt.Sprintf("response %+v", res))
	hits, ok := res["hits"]
	if !ok {
		log.Println("no hits")
		return nil, nil
	}
	hitsMap, isMap := hits.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("hits is not a map")
	}
	hitResults, ok := hitsMap["hits"]
	if !ok {
		log.Println("no hit results")
		return nil, nil
	}
	hitList, isList := hitResults.([]interface{})
	if !isList {
		return nil, fmt.Errorf("results is not a list")
	}
	var out []Document
	for _, h := range hitList {
		m, isMap := h.(map[string]interface{})
		if !isMap {
			return nil, fmt.Errorf("not a map")
		}
		_source, ok := m["_source"]
		if !ok {
			continue
		}
		source, isMap := _source.(map[string]string)
		if !isMap {
			return nil, fmt.Errorf("source is not a map")
		}
		out = append(out, Document{})
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
