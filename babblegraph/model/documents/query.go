package documents

import (
	"babblegraph/util/elastic/esquery"
	"babblegraph/wordsmith"
	"encoding/json"
	"strings"
)

type documentsQueryBuilder struct {
	Terms           []string
	Language        *string
	SentDocumentIDs []string
}

func NewDocumentsQueryBuilder() *documentsQueryBuilder {
	return &documentsQueryBuilder{}
}

func (d *documentsQueryBuilder) ContainingTerms(terms []wordsmith.LemmaID) {
	d.Terms = append(d.Terms, termsToString(terms)...)
}

func (d *documentsQueryBuilder) ForLanguage(languageCode wordsmith.LanguageCode) {
	language := languageCode.Str()
	d.Language = &language
}

func (d *documentsQueryBuilder) NotContainingDocumentIDs(docIDs []DocumentID) {
	var ids []string
	for _, docID := range docIDs {
		ids = append(ids, string(docID))
	}
	d.SentDocumentIDs = append(d.SentDocumentIDs, ids...)
}

func termsToString(terms []wordsmith.LemmaID) []string {
	var out []string
	for _, t := range terms {
		out = append(out, string(t))
	}
	return out
}

func (d *documentsQueryBuilder) ExecuteQuery() ([]Document, error) {
	queryBuilder := esquery.NewBoolQueryBuilder()
	if len(d.Terms) != 0 {
		queryBuilder.AddMust(esquery.Match("lemmatized_body", strings.Join(d.Terms, " ")))
	}
	if len(d.SentDocumentIDs) != 0 {
		queryBuilder.AddMustNot(esquery.IDs(d.SentDocumentIDs))
	}
	var docs []Document
	if err := esquery.ExecuteSearch(documentIndex{}, queryBuilder.BuildBoolQuery(), func(source []byte) error {
		var doc Document
		if err := json.Unmarshal(source, &doc); err != nil {
			return err
		}
		docs = append(docs, doc)
		return nil
	}); err != nil {
		return nil, err
	}
	return docs, nil
}
