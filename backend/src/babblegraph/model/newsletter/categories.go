package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/useraccounts"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"sort"
)

func getDocumentCategories(languageCode wordsmith.LanguageCode, userAccessor userPreferencesAccessor, docsAccessor documentAccessor) ([]Category, error) {
	topics := getTopicsForNewsletter(userAccessor)
	allowableDomains, err := getAllowableDomains(userAccessor)
	if err != nil {
		return nil, err
	}
	genericDocuments, err := docsAccessor.GetDocumentsForUser(getDocumentsForUserInput{
		getDocumentsBaseInput: getDocumentsBaseInput{
			LanguageCode:        languageCode,
			ExcludedDocumentIDs: userAccessor.getSentDocumentIDs(),
			ValidDomains:        allowableDomains,
			MinimumReadingLevel: ptr.Int64(userAccessor.getReadingLevel().LowerBound),
			MaximumReadingLevel: ptr.Int64(userAccessor.getReadingLevel().UpperBound),
		},
		Lemmas: userAccessor.getTrackingLemmas(),
	})
	if err != nil {
		return nil, err
	}
	documentsByTopic := make(map[contenttopics.ContentTopic][]documents.DocumentWithScore)
	for _, t := range topics {
		documentsForTopic, err := docsAccessor.GetDocumentsForUser(getDocumentsForUserInput{
			getDocumentsBaseInput: getDocumentsBaseInput{
				LanguageCode:        languageCode,
				ExcludedDocumentIDs: userAccessor.getSentDocumentIDs(),
				ValidDomains:        allowableDomains,
				MinimumReadingLevel: ptr.Int64(userAccessor.getReadingLevel().LowerBound),
				MaximumReadingLevel: ptr.Int64(userAccessor.getReadingLevel().UpperBound),
			},
			Topic:  t.Ptr(),
			Lemmas: userAccessor.getTrackingLemmas(),
		})
		switch {
		case err != nil:
			return nil, err
		case len(documentsForTopic) == 0:
			// no-op
		default:
			documentsByTopic[t] = documentsForTopic
		}
	}
	// TODO: figure out number of documents
	return joinDocumentsIntoCategories(joinDocumentsIntoCategoriesInput{
		userAccessor: userAccessor,
		languageCode: languageCode,
		// numberOfDocumentsInNewsletter:
		documentsByTopic: documentsByTopic,
		genericDocuments: genericDocuments,
	})
}

type joinDocumentsIntoCategoriesInput struct {
	userAccessor                  userPreferencesAccessor
	languageCode                  wordsmith.LanguageCode
	numberOfDocumentsInNewsletter int
	documentsByTopic              map[contenttopics.ContentTopic][]documents.DocumentWithScore
	genericDocuments              []documents.DocumentWithScore
}

func joinDocumentsIntoCategories(input joinDocumentsIntoCategoriesInput) ([]Category, error) {
	type scoredDocumentsWithTopic struct {
		documentsWithScore []documents.DocumentWithScore
		topic              contenttopics.ContentTopic
	}
	var docsWithTopic []scoredDocumentsWithTopic
	for topic, scoredDocuments := range input.documentsByTopic {
		docsWithTopic = append(docsWithTopic, scoredDocumentsWithTopic{
			documentsWithScore: scoredDocuments,
			topic:              topic,
		})
	}
	sort.Slice(docsWithTopic, func(i, j int) bool {
		return docsWithTopic[i].documentsWithScore[0].Score.GreaterThan(docsWithTopic[j].documentsWithScore[0].Score)
	})
	// HACK: Using URL Identifier here instead of document ID
	// because of an issue with urlparser means that any document < Version5
	// may appear multiple times in the same email
	documentsInEmailByURLIdentifier := make(map[string]bool)
	var categories []Category
	for idx, documentGroup := range docsWithTopic {
		documentsPerTopic := (input.numberOfDocumentsInNewsletter - len(documentsInEmailByURLIdentifier)) / (len(docsWithTopic) - idx)
		documentCounter := 0
		var links []Link
		for i := 0; i < len(documentGroup.documentsWithScore) && documentCounter < documentsPerTopic; i++ {
			doc := documentGroup.documentsWithScore[i].Document
			u := urlparser.MustParseURL(doc.URL)
			if _, ok := documentsInEmailByURLIdentifier[u.URLIdentifier]; !ok {
				link, err := makeLinkFromDocument(input.userAccessor, doc)
				if err != nil {
					return nil, err
				}
				documentCounter++
				documentsInEmailByURLIdentifier[u.URLIdentifier] = true
				links = append(links, *link)
			}
		}
		if len(links) > 0 {
			var categoryName *string
			displayName, err := contenttopics.ContentTopicNameToDisplayName(documentGroup.topic)
			if err != nil {
				log.Println(fmt.Sprintf("Error generating display name: %s", err.Error()))
			} else {
				categoryName = ptr.String(text.ToTitleCaseForLanguage(displayName.Str(), input.languageCode))
			}
			categories = append(categories, Category{
				Name:  categoryName,
				Links: links,
			})
		}
	}
	if len(documentsInEmailByURLIdentifier) < input.numberOfDocumentsInNewsletter {
		maxGenericDocuments := input.numberOfDocumentsInNewsletter - len(documentsInEmailByURLIdentifier)
		documentCounter := 0
		var links []Link
		for i := 0; i < len(input.genericDocuments) && documentCounter < maxGenericDocuments; i++ {
			doc := input.genericDocuments[i].Document
			u := urlparser.MustParseURL(doc.URL)
			if _, ok := documentsInEmailByURLIdentifier[u.URLIdentifier]; !ok {
				link, err := makeLinkFromDocument(input.userAccessor, doc)
				if err != nil {
					return nil, err
				}
				documentCounter++
				documentsInEmailByURLIdentifier[u.URLIdentifier] = true
				links = append(links, *link)
			}
		}
		if len(links) > 0 {
			var categoryName *string
			if len(categories) > 0 {
				displayName := contenttopics.GenericCategoryNameForLanguage(input.languageCode)
				categoryName = ptr.String(text.ToTitleCaseForLanguage(displayName.Str(), input.languageCode))
			}
			categories = append(categories, Category{
				Name:  categoryName,
				Links: links,
			})
		}
	}
	return categories, nil
}

func getTopicsForNewsletter(accessor userPreferencesAccessor) []contenttopics.ContentTopic {
	userSubscriptionLevel := accessor.getUserSubscriptionLevel()
	allUserTopics := accessor.getUserTopics()
	var topics []contenttopics.ContentTopic
	for _, idx := range pickUpToNRandomIndices(int(len(allUserTopics)), defaultNumberOfTopicsPerEmail) {
		topics = append(topics, allUserTopics[idx])
	}
	switch {
	case userSubscriptionLevel == nil:
		return topics
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		if userScheduleForDay := accessor.getUserScheduleForDay(); userScheduleForDay != nil && len(userScheduleForDay.ContentTopics) != 0 {
			topics := userScheduleForDay.ContentTopics
			if len(topics) < defaultNumberOfTopicsPerEmail {
				for _, idx := range pickUpToNRandomIndices(int(len(allUserTopics)), defaultNumberOfTopicsPerEmail-len(topics)) {
					topics = append(topics, allUserTopics[idx])
				}
			}
		}
		return topics
	default:
		panic("Unreachable")
	}
}
