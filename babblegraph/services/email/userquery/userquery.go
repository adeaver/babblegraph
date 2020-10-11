package userquery

import (
	"babblegraph/model/documents"
	"babblegraph/services/email/labels"
	"babblegraph/services/email/userprefs"
	"babblegraph/wordsmith"
	"fmt"
	"log"
)

const maxDocumentsPerEmail = 5

func GetDocumentsForUser(userInfos []userprefs.UserEmailInfo) (map[string][]documents.Document, error) {
	labelSearchTerms, err := labels.GetLemmaIDsForLabelNames()
	if err != nil {
		return nil, err
	}
	out := make(map[string][]documents.Document)
	for _, userInfo := range userInfos {
		terms := getTermsForUser(labelSearchTerms, userInfo)
		docs, err := queryDocsForUser(userInfo, terms)
		if err != nil {
			return nil, err
		}
		topDocs := pickTopDocuments(docs)
		log.Println(fmt.Sprintf("Found docs %+v for %s", topDocs, userInfo.EmailAddress))
		out[userInfo.EmailAddress] = topDocs
	}
	return out, nil
}

func getTermsForUser(lemmaIDsForLabelName map[labels.LabelName][]wordsmith.LemmaID, userInfo userprefs.UserEmailInfo) []wordsmith.LemmaID {
	var out []wordsmith.LemmaID
	for _, label := range userInfo.InterestLabels {
		lemmaIDs, ok := lemmaIDsForLabelName[label]
		if !ok {
			log.Println(fmt.Sprintf("No lemmas found label %s", label))
			continue
		}
		out = append(out, lemmaIDs...)
	}
	return out
}

func queryDocsForUser(userInfo userprefs.UserEmailInfo, terms []wordsmith.LemmaID) ([]documents.Document, error) {
	docQueryBuilder := documents.NewDocumentsQueryBuilder()
	/*
		Later this will become a problem
		that I will need to map label -> map[languageCode][]LemmaID
		and likely send one email per language
	*/
	docQueryBuilder.ContainingTerms(terms)
	docQueryBuilder.ForLanguage(userInfo.Languages[0])
	// TODO: add reading level and docs to be removed
	return docQueryBuilder.ExecuteQuery()
}

func pickTopDocuments(docs []documents.Document) []documents.Document {
	stopIdx := maxDocumentsPerEmail
	if len(docs) < stopIdx {
		stopIdx = len(docs)
	}
	var out []documents.Document
	for i := 0; i < stopIdx; i++ {
		out = append(out, docs[i])
	}
	return out
}
