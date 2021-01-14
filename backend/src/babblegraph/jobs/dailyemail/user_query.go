package dailyemail

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/util/ptr"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	maxDocumentsPerEmail = 8
	maxTopicsPerEmail    = 4
)

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

	genericDocuments, err := docQueryBuilder.ExecuteQuery()
	if err != nil {
		return nil, err
	}
	if len(userInfo.Topics) == 0 {
		return genericDocuments, nil
	}
	userDocuments := make(map[documents.DocumentID]documents.Document)
	topics := pickTopics(userInfo.Topics)
	for _, topic := range topics {
		// This is a bit of a hack.
		// We iteratre through the topics and clobber the topic
		// And rerun the query.
		docQueryBuilder.ForTopic(topic.Ptr())
		documents, err := docQueryBuilder.ExecuteQuery()
		if err != nil {
			return nil, err
		}
		for _, doc := range documents {
			userDocuments[doc.ID] = doc
		}
	}
	if len(userDocuments) == 0 {
		return genericDocuments, nil
	}
	var out []documents.Document
	for _, doc := range userDocuments {
		out = append(out, doc)
	}
	return out, nil
}

func pickTopDocuments(docs []documents.Document) []documents.Document {
	stopIdx := maxDocumentsPerEmail
	if len(docs) < stopIdx {
		stopIdx = len(docs)
	}
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	var out []documents.Document
	for i := 0; i < stopIdx; i++ {
		idx := generator.Intn(int(len(docs)))
		out = append(out, docs[idx])
		docs = append(docs[:idx], docs[idx+1:]...)
	}
	return out
}

func pickTopics(topics []contenttopics.ContentTopic) []contenttopics.ContentTopic {
	stopIdx := maxTopicsPerEmail
	if len(topics) < stopIdx {
		stopIdx = len(topics)
	}
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	var out []contenttopics.ContentTopic
	for i := 0; i < stopIdx; i++ {
		idx := generator.Intn(int(len(topics)))
		out = append(out, topics[idx])
		topics = append(topics[:idx], topics[idx+1:]...)
	}
	return out

}
