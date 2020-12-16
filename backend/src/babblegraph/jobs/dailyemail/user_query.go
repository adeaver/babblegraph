package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/util/ptr"

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
	return topDocs, nil
}

func queryDocsForUser(userInfo userEmailInfo) ([]documents.Document, error) {
	docQueryBuilder := documents.NewDocumentsQueryBuilderForLanguage(userInfo.Languages[0])
	readingLevelLowerBound := ptr.Int64(userInfo.ReadingLevel.LowerBound)
	readingLevelUpperBound := ptr.Int64(userInfo.ReadingLevel.UpperBound)

	docQueryBuilder.NotContainingDocuments(userInfo.SentDocuments)
	docQueryBuilder.ForVersionRange(documents.Version2.Ptr(), documents.Version2.Ptr())
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
