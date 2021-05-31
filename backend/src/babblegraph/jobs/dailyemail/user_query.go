package dailyemail

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"math/rand"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	maxDocumentsPerEmail = 12
	maxTopicsPerEmail    = 4
)

func getDocumentsForUser(tx *sqlx.Tx, userInfo userEmailInfo) ([]email_actions.CategorizedDocuments, error) {
	docs, genericDocs, err := queryDocsForUser(userInfo)
	if err != nil {
		return nil, err
	}
	return pickTopDocuments(docs, genericDocs), nil
}

type documentsWithTopic struct {
	topic     *contenttopics.ContentTopic
	documents []documents.DocumentWithScore
}

func getAllowedDomains(userDomainCounts map[string]int64) ([]string, error) {
	var out []string
	for _, d := range domains.GetDomains() {
		countForDomain, ok := userDomainCounts[d]
		if ok {
			metadata, err := domains.GetDomainMetadata(d)
			if err != nil {
				return nil, err
			}
			if metadata.NumberOfMonthlyFreeArticles != nil && countForDomain > *metadata.NumberOfMonthlyFreeArticles {
				continue
			}
		}
		out = append(out, d)
	}
	return out, nil
}

func queryDocsForUser(userInfo userEmailInfo) (_categorizedDocument []documentsWithTopic, _genericDocuments []documents.DocumentWithScore, _err error) {
	var trackingLemmas []wordsmith.LemmaID
	for _, lemmaMapping := range userInfo.TrackingLemmas {
		if lemmaMapping.IsActive {
			trackingLemmas = append(trackingLemmas, lemmaMapping.LemmaID)
		}
	}
	userDomainCounts := make(map[string]int64)
	for _, domainCount := range userInfo.UserDomainCounts {
		userDomainCounts[domainCount.Domain] = domainCount.Count
	}
	allowableDomains, err := getAllowedDomains(userDomainCounts)
	if err != nil {
		return nil, nil, err
	}

	docQueryBuilder := documents.NewDocumentsQueryBuilderForLanguage(userInfo.Languages[0])
	readingLevelLowerBound := ptr.Int64(userInfo.ReadingLevel.LowerBound)
	readingLevelUpperBound := ptr.Int64(userInfo.ReadingLevel.UpperBound)

	docQueryBuilder.WithValidDomains(allowableDomains)
	docQueryBuilder.NotContainingDocuments(userInfo.SentDocuments)
	docQueryBuilder.ForVersionRange(documents.Version3.Ptr(), documents.Version6.Ptr())
	docQueryBuilder.ForReadingLevelRange(readingLevelLowerBound, readingLevelUpperBound)
	docQueryBuilder.ContainingLemmas(trackingLemmas)

	genericDocuments, err := docQueryBuilder.ExecuteQuery()
	if err != nil {
		return nil, nil, err
	}
	if len(userInfo.Topics) == 0 {
		return nil, genericDocuments, nil
	}
	var outDocuments []documentsWithTopic
	topics := pickTopics(userInfo.Topics)
	for _, topic := range topics {
		// This is a bit of a hack.
		// We iteratre through the topics and clobber the topic
		// And rerun the query.
		docQueryBuilder.ForTopic(topic.Ptr())
		documents, err := docQueryBuilder.ExecuteQuery()
		if err != nil {
			return nil, nil, err
		}
		if len(documents) > 0 {
			outDocuments = append(outDocuments, documentsWithTopic{
				topic:     topic.Ptr(),
				documents: documents,
			})
		}
	}
	return outDocuments, genericDocuments, nil
}

func pickTopDocuments(docsWithTopic []documentsWithTopic, genericDocuments []documents.DocumentWithScore) []email_actions.CategorizedDocuments {
	sort.Slice(docsWithTopic, func(i, j int) bool {
		return docsWithTopic[i].documents[0].Score.GreaterThan(docsWithTopic[i].documents[0].Score)
	})
	var categorizedDocuments []email_actions.CategorizedDocuments
	// HACK: Using URL Identifier here instead of document ID
	// because of an issue with urlparser means that any document < Version5
	// may appear multiple times in the same email
	documentsInEmailByURLIdentifier := make(map[string]bool)
	if len(docsWithTopic) > 0 {
		documentsPerTopic := maxDocumentsPerEmail / len(docsWithTopic)
		for _, docs := range docsWithTopic {
			documentCounter := 0
			var documents []documents.Document
			for i := 0; i < len(docs.documents) && documentCounter < documentsPerTopic; i++ {
				doc := docs.documents[i].Document
				u := urlparser.MustParseURL(doc.URL)
				if _, ok := documentsInEmailByURLIdentifier[u.URLIdentifier]; !ok {
					documents = append(documents, doc)
					documentCounter++
					documentsInEmailByURLIdentifier[u.URLIdentifier] = true
				}
			}
			categorizedDocuments = append(categorizedDocuments, email_actions.CategorizedDocuments{
				Topic:     docs.topic,
				Documents: documents,
			})
		}
	}
	if len(documentsInEmailByURLIdentifier) < maxDocumentsPerEmail {
		var selectedGenericDocuments []documents.Document
		maxGenericDocuments := maxDocumentsPerEmail - len(documentsInEmailByURLIdentifier)
		documentCounter := 0
		for i := 0; i < len(genericDocuments) && documentCounter < maxGenericDocuments; i++ {
			doc := genericDocuments[i].Document
			u := urlparser.MustParseURL(doc.URL)
			if _, ok := documentsInEmailByURLIdentifier[u.URLIdentifier]; !ok {
				selectedGenericDocuments = append(selectedGenericDocuments, doc)
				documentCounter++
				documentsInEmailByURLIdentifier[u.URLIdentifier] = true
			}
		}
		categorizedDocuments = append(categorizedDocuments, email_actions.CategorizedDocuments{
			Documents: selectedGenericDocuments,
		})
	}
	return categorizedDocuments
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
