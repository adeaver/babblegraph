package userquery

import (
	"babblegraph/model/documents"
	"babblegraph/model/userdocuments"
	"babblegraph/services/email/userprefs"
	"babblegraph/util/ptr"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const maxDocumentsPerEmail = 5

func GetDocumentsForUser(tx *sqlx.Tx, userInfos []userprefs.UserEmailInfo) (map[string][]documents.Document, error) {
	out := make(map[string][]documents.Document)
	for _, userInfo := range userInfos {
		sentDocumentIDs, err := userdocuments.GetDocumentIDsSentToUser(tx, userInfo.UserID)
		if err != nil {
			return nil, err
		}
		docs, err := queryDocsForUser(userInfo, sentDocumentIDs)
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

func queryDocsForUser(userInfo userprefs.UserEmailInfo, sentDocumentIDs []documents.DocumentID) ([]documents.Document, error) {
	docQueryBuilder := documents.NewDocumentsQueryBuilderForLanguage(userInfo.Languages[0])
	readingLevelLowerBound := ptr.Int64(userInfo.ReadingLevel.LowerBound)
	readingLevelUpperBound := ptr.Int64(userInfo.ReadingLevel.UpperBound)

	docQueryBuilder.NotContainingDocuments(sentDocumentIDs)
	docQueryBuilder.ForVersionRange(documents.Version2.Ptr(), documents.Version4.Ptr())
	docQueryBuilder.ForReadingLevelRange(readingLevelLowerBound, readingLevelUpperBound)

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
