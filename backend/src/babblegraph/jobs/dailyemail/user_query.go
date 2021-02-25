package dailyemail

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"math/rand"
	"sort"
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

func queryDocsForUser(userInfo userEmailInfo) ([]documents.DocumentWithScore, error) {
	var trackingLemmas []wordsmith.LemmaID
	for _, lemmaMapping := range userInfo.TrackingLemmas {
		if lemmaMapping.IsActive {
			trackingLemmas = append(trackingLemmas, lemmaMapping.LemmaID)
		}
	}
	docQueryBuilder := documents.NewDocumentsQueryBuilderForLanguage(userInfo.Languages[0])
	readingLevelLowerBound := ptr.Int64(userInfo.ReadingLevel.LowerBound)
	readingLevelUpperBound := ptr.Int64(userInfo.ReadingLevel.UpperBound)

	docQueryBuilder.NotContainingDocuments(userInfo.SentDocuments)
	docQueryBuilder.ForVersionRange(documents.Version2.Ptr(), documents.Version4.Ptr())
	docQueryBuilder.ForReadingLevelRange(readingLevelLowerBound, readingLevelUpperBound)
	docQueryBuilder.ContainingLemmas(trackingLemmas)

	genericDocuments, err := docQueryBuilder.ExecuteQuery()
	if err != nil {
		return nil, err
	}
	if len(userInfo.Topics) == 0 {
		return genericDocuments, nil
	}
	var outDocuments []documents.DocumentWithScore
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
		outDocuments = append(outDocuments, documents...)
	}
	if len(outDocuments) == 0 {
		return genericDocuments, nil
	}
	return outDocuments, nil
}

func pickTopDocuments(docs []documents.DocumentWithScore) []documents.Document {
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Score.GreaterThan(docs[j].Score)
	})
	var out []documents.Document
	docIDHash := make(map[documents.DocumentID]bool)
	for i := 0; i < len(docs) && len(out) < maxDocumentsPerEmail; i++ {
		docID := docs[i].Document.ID
		if ok, _ := docIDHash[docID]; !ok {
			out = append(out, docs[i].Document)
			docIDHash[docID] = true
		}
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
