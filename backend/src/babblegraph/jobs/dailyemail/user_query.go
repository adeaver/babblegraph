package dailyemail

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/model/userlemma"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	defaultNumberOfArticlesPerEmail = 12
	defaultNumberOfTopicsPerEmail   = 4

	minimumDaysSinceLastSpotlight = 3
)

func getSpotlightDocumentForUser(tx *sqlx.Tx, userInfo userEmailInfo, userScheduleForDay *usernewsletterschedule.UserNewsletterScheduleDayMetadata) (*documents.DocumentWithScore, error) {
	potentialSpotlightLemmas, err := getPotentialSpotlightLemmasForUser(tx, userInfo)
	if err != nil {
		return nil, err
	}

	for _, spotlightLemma := range potentialSpotlightLemmas {
		log.Println(spotlightLemma)
	}
	return nil, nil
}

func getPotentialSpotlightLemmasForUser(tx *sqlx.Tx, userInfo userEmailInfo) ([]wordsmith.LemmaID, error) {
	/*
	   Select all spotlighted words for a user ordered by the date that they were sent.
	   Select all active lemmas from a user's tracking list
	   Collect a list that is ordered by:
	       - First group: has never been sent
	       - Second group: has been sent, ordered by number of days since last sent (at least 3 days ago)
	*/
	lemmaReinforcementSpotlightRecords, err := userlemma.GetLemmaReinforcementRecordsForUserOrderedBySentOn(tx, userInfo.UserID)
	if err != nil {
		return nil, err
	}
	lemmaSpotlightRecordSentOnTimeByID := make(map[wordsmith.LemmaID]time.Time)
	for _, spotlightRecord := range lemmaReinforcementSpotlightRecords {
		lemmaSpotlightRecordSentOnTimeByID[spotlightRecord.LemmaID] = spotlightRecord.LastSentOn
	}
	now := time.Now()
	var lemmasNotSent, sentLemmas []wordsmith.LemmaID
	for _, lemmaMapping := range userInfo.TrackingLemmas {
		if !lemmaMapping.IsActive {
			continue
		}
		lastSent, ok := lemmaSpotlightRecordSentOnTimeByID[lemmaMapping.LemmaID]
		if !ok {
			lemmasNotSent = append(lemmasNotSent, lemmaMapping.LemmaID)
		} else {
			if !now.Add(-1 * minimumDaysSinceLastSpotlight * 24 * time.Hour).Before(lastSent) {
				sentLemmas = append(sentLemmas, lemmaMapping.LemmaID)
			}
		}
	}
	sort.Slice(sentLemmas, func(i, j int) bool {
		iSentOn, _ := lemmaSpotlightRecordSentOnTimeByID[sentLemmas[i]]
		jSentOn, _ := lemmaSpotlightRecordSentOnTimeByID[sentLemmas[j]]
		return iSentOn.Before(jSentOn)
	})
	return append(lemmasNotSent, sentLemmas...), nil
}

func getDocumentsForUser(tx *sqlx.Tx, userInfo userEmailInfo, userScheduleForDay *usernewsletterschedule.UserNewsletterScheduleDayMetadata) ([]email_actions.CategorizedDocuments, error) {
	docs, genericDocs, err := queryDocsForUser(userInfo, userScheduleForDay)
	if err != nil {
		return nil, err
	}
	numberOfDocumentsPerEmail := defaultNumberOfArticlesPerEmail
	if userScheduleForDay != nil {
		numberOfDocumentsPerEmail = userScheduleForDay.NumberOfArticles
	}
	return pickTopDocuments(docs, genericDocs, numberOfDocumentsPerEmail), nil
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
			if metadata.NumberOfMonthlyFreeArticles != nil && countForDomain >= *metadata.NumberOfMonthlyFreeArticles {
				continue
			}
		}
		out = append(out, d)
	}
	return out, nil
}

func queryDocsForUser(userInfo userEmailInfo, userScheduleForDay *usernewsletterschedule.UserNewsletterScheduleDayMetadata) (_categorizedDocument []documentsWithTopic, _genericDocuments []documents.DocumentWithScore, _err error) {
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

	dailyEmailDocQueryBuilder := documents.NewDailyEmailDocumentsQueryBuilder()
	dailyEmailDocQueryBuilder.ContainingLemmas(trackingLemmas)

	genericDocuments, err := documents.ExecuteDocumentQuery(dailyEmailDocQueryBuilder, documents.ExecuteDocumentQueryInput{
		LanguageCode:        userInfo.Languages[0],
		ValidDomains:        allowableDomains,
		ExcludedDocumentIDs: userInfo.SentDocuments,
		MinimumReadingLevel: ptr.Int64(userInfo.ReadingLevel.LowerBound),
		MaximumReadingLevel: ptr.Int64(userInfo.ReadingLevel.UpperBound),
	})
	if err != nil {
		return nil, nil, err
	}
	if len(userInfo.Topics) == 0 {
		return nil, genericDocuments, nil
	}
	var outDocuments []documentsWithTopic
	var topics []contenttopics.ContentTopic
	switch {
	case userScheduleForDay != nil && len(userScheduleForDay.ContentTopics) != 0:
		topics = userScheduleForDay.ContentTopics
		if len(topics) < defaultNumberOfTopicsPerEmail {
			var remainingTopics []contenttopics.ContentTopic
			for _, t := range userInfo.Topics {
				for _, selectedTopic := range topics {
					if t == selectedTopic {
						continue
					}
				}
				remainingTopics = append(remainingTopics, t)
			}
			topics = append(topics, pickTopics(remainingTopics, defaultNumberOfTopicsPerEmail-len(topics))...)
		}
	case userScheduleForDay != nil && len(userScheduleForDay.ContentTopics) == 0,
		userScheduleForDay == nil:
		topics = pickTopics(userInfo.Topics, defaultNumberOfTopicsPerEmail)
	}
	for _, topic := range topics {
		// This is a bit of a hack.
		// We iteratre through the topics and clobber the topic
		// And rerun the query.
		dailyEmailDocQueryBuilder.ForTopic(topic.Ptr())
		documents, err := documents.ExecuteDocumentQuery(dailyEmailDocQueryBuilder, documents.ExecuteDocumentQueryInput{
			LanguageCode:        userInfo.Languages[0],
			ValidDomains:        allowableDomains,
			ExcludedDocumentIDs: userInfo.SentDocuments,
			MinimumReadingLevel: ptr.Int64(userInfo.ReadingLevel.LowerBound),
			MaximumReadingLevel: ptr.Int64(userInfo.ReadingLevel.UpperBound),
		})

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

func pickTopDocuments(docsWithTopic []documentsWithTopic, genericDocuments []documents.DocumentWithScore, numberOfDocuments int) []email_actions.CategorizedDocuments {
	sort.Slice(docsWithTopic, func(i, j int) bool {
		return docsWithTopic[i].documents[0].Score.GreaterThan(docsWithTopic[i].documents[0].Score)
	})
	var categorizedDocuments []email_actions.CategorizedDocuments
	// HACK: Using URL Identifier here instead of document ID
	// because of an issue with urlparser means that any document < Version5
	// may appear multiple times in the same email
	documentsInEmailByURLIdentifier := make(map[string]bool)
	if len(docsWithTopic) > 0 {
		for idx, docs := range docsWithTopic {
			documentsPerTopic := (numberOfDocuments - len(documentsInEmailByURLIdentifier)) / (len(docsWithTopic) - idx)
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
			if documentCounter > 0 {
				categorizedDocuments = append(categorizedDocuments, email_actions.CategorizedDocuments{
					Topic:     docs.topic,
					Documents: documents,
				})
			}
		}
	}
	if len(documentsInEmailByURLIdentifier) < numberOfDocuments {
		var selectedGenericDocuments []documents.Document
		maxGenericDocuments := numberOfDocuments - len(documentsInEmailByURLIdentifier)
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

func pickTopics(topics []contenttopics.ContentTopic, numberOfTopics int) []contenttopics.ContentTopic {
	stopIdx := numberOfTopics
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
