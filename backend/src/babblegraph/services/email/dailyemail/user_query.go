package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/model/userdocuments"
	"babblegraph/util/ptr"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const maxDocumentsPerEmail = 5

func getDocumentsForUser(tx *sqlx.Tx, userInfo userEmailInfo) ([]documents.Document, error) {
	docs, err := queryDocsForUser(userInfo)
	if err != nil {
		return nil, err
	}
	topDocs := pickTopDocuments(docs)
	var docIDs []documents.DocumentID
	for _, doc := range topDocs {
		docIDs = append(docIDs, doc.ID)
	}
	log.Println(fmt.Sprintf("Found doc IDs %+v for user %s", docIDs, userInfo.EmailAddress))
	if err := userdocuments.InsertDocumentIDsForUser(tx, userInfo.UserID, docIDs); err != nil {
		return nil, err
	}
	return topDocs, nil
}

func queryDocsForUser(userInfo userEmailInfo) ([]documents.Document, error) {
	docQueryBuilder := documents.NewDocumentsQueryBuilderForLanguage(userInfo.Languages[0])
	readingLevelLowerBound := ptr.Int64(userInfo.ReadingLevel.LowerBound)
	readingLevelUpperBound := ptr.Int64(userInfo.ReadingLevel.UpperBound)

	docQueryBuilder.NotContainingDocuments(userInfo.SentDocuments)
	docQueryBuilder.ForVersionRange(nil, documents.Version2.Ptr())
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
