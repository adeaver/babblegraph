package userquery

import (
	"babblegraph/model/documents"
	"babblegraph/model/userdocuments"
	"babblegraph/services/email/labels"
	"babblegraph/services/email/userprefs"
	"babblegraph/wordsmith"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const maxDocumentsPerEmail = 5

func GetDocumentsForUser(tx *sqlx.Tx, userInfos []userprefs.UserEmailInfo) (map[string][]documents.Document, error) {
	labelSearchTerms, err := labels.GetLemmaIDsForLabelNames()
	if err != nil {
		return nil, err
	}
	out := make(map[string][]documents.Document)
	for _, userInfo := range userInfos {
		sentDocumentIDs, err := userdocuments.GetDocumentIDsSentToUser(tx, userInfo.UserID)
		if err != nil {
			return nil, err
		}
		terms := getTermsForUser(labelSearchTerms, userInfo)
		docs, err := queryDocsForUser(userInfo, terms, sentDocumentIDs)
		if err != nil {
			return nil, err
		}
		topDocs := pickTopDocuments(docs)
		var docIDs []documents.DocumentID
		for _, doc := range topDocs {
			docIDs = append(docIDs, doc.ID)
		}
		out[userInfo.EmailAddress] = topDocs
		log.Println(fmt.Sprintf("Found doc IDs %+v for user %s", docIDs, userInfo.EmailAddress))
		if err := userdocuments.InsertDocumentIDsForUser(tx, userInfo.UserID, docIDs); err != nil {
			return nil, err
		}
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

func queryDocsForUser(userInfo userprefs.UserEmailInfo, terms []wordsmith.LemmaID, sentDocumentIDs []documents.DocumentID) ([]documents.Document, error) {
	docQueryBuilder := documents.NewDocumentsQueryBuilder()
	/*
		Later this will become a problem
		that I will need to map label -> map[languageCode][]LemmaID
		and likely send one email per language
	*/
	docQueryBuilder.ContainingTerms(terms)
	docQueryBuilder.ForLanguage(userInfo.Languages[0])
	docQueryBuilder.NotContainingDocumentIDs(sentDocumentIDs)
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
