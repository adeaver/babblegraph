package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/useraccounts"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"sort"
)

type getDocumentCategoriesInput struct {
	emailRecordID                 email.ID
	languageCode                  wordsmith.LanguageCode
	userAccessor                  userPreferencesAccessor
	docsAccessor                  documentAccessor
	contentAccessor               contentAccessor
	podcastAccessor               podcastAccessor
	numberOfDocumentsInNewsletter *int
}

func getDocumentCategories(c ctx.LogContext, input getDocumentCategoriesInput) ([]Category, error) {
	topics := getTopicsForNewsletter(input.userAccessor)
	c.Debugf("Topics %+v", topics)
	allowableSourceIDs := input.userAccessor.getAllowableSources()
	genericDocuments, err := input.docsAccessor.GetDocumentsForUser(c, getDocumentsForUserInput{
		getDocumentsBaseInput: getDocumentsBaseInput{
			LanguageCode:        input.languageCode,
			ExcludedDocumentIDs: input.userAccessor.getSentDocumentIDs(),
			ValidSourceIDs:      allowableSourceIDs,
			MinimumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().LowerBound),
			MaximumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().UpperBound),
		},
		Lemmas: input.userAccessor.getTrackingLemmas(),
	})
	if err != nil {
		return nil, err
	}
	documentsByTopic := make(map[content.TopicID][]documents.DocumentWithScore)
	for _, t := range topics {
		documentsForTopic, err := input.docsAccessor.GetDocumentsForUser(c, getDocumentsForUserInput{
			getDocumentsBaseInput: getDocumentsBaseInput{
				LanguageCode:        input.languageCode,
				ExcludedDocumentIDs: input.userAccessor.getSentDocumentIDs(),
				ValidSourceIDs:      allowableSourceIDs,
				MinimumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().LowerBound),
				MaximumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().UpperBound),
			},
			Lemmas: input.userAccessor.getTrackingLemmas(),
			Topic:  t.Ptr(),
		})
		switch {
		case err != nil:
			return nil, err
		case len(documentsForTopic.RecentDocuments)+len(documentsForTopic.NonRecentDocuments) == 0:
			c.Infof("No documents for topic %s", t.Str())
		default:
			documentsByTopic[t] = append(documentsForTopic.RecentDocuments, documentsForTopic.NonRecentDocuments...)
			c.Infof("Documents for topic %s: %+v", t.Str(), documentsForTopic)
		}
	}
	var podcastEpisodesByTopic map[content.TopicID][]podcasts.Episode
	if input.userAccessor.getUserSubscriptionLevel() != nil {
		podcastEpisodesByTopic, err = input.podcastAccessor.LookupPodcastEpisodesForTopics(topics)
	}
	return joinDocumentsIntoCategories(c, joinDocumentsIntoCategoriesInput{
		emailRecordID:                 input.emailRecordID,
		userAccessor:                  input.userAccessor,
		contentAccessor:               input.contentAccessor,
		podcastAccessor:               input.podcastAccessor,
		languageCode:                  input.languageCode,
		numberOfDocumentsInNewsletter: deref.Int(input.numberOfDocumentsInNewsletter, DefaultNumberOfArticlesPerEmail),
		documentsByTopic:              documentsByTopic,
		podcastEpisodesByTopic:        podcastEpisodesByTopic,
		genericDocuments:              append(genericDocuments.RecentDocuments, genericDocuments.NonRecentDocuments...),
	})
}

type joinDocumentsIntoCategoriesInput struct {
	emailRecordID                 email.ID
	userAccessor                  userPreferencesAccessor
	contentAccessor               contentAccessor
	podcastAccessor               podcastAccessor
	languageCode                  wordsmith.LanguageCode
	numberOfDocumentsInNewsletter int
	documentsByTopic              map[content.TopicID][]documents.DocumentWithScore
	podcastEpisodesByTopic        map[content.TopicID][]podcasts.Episode
	genericDocuments              []documents.DocumentWithScore
}

func joinDocumentsIntoCategories(c ctx.LogContext, input joinDocumentsIntoCategoriesInput) ([]Category, error) {
	type scoredDocumentsWithTopic struct {
		documentsWithScore []documents.DocumentWithScore
		topic              content.TopicID
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
	podcastsInEmailByID := make(map[podcasts.EpisodeID]bool)
	var categories []Category
	for idx, documentGroup := range docsWithTopic {
		documentsPerTopic := (input.numberOfDocumentsInNewsletter - len(documentsInEmailByURLIdentifier)) / (len(docsWithTopic) - idx)
		documentCounter := 0
		var links []Link
		for i := 0; i < len(documentGroup.documentsWithScore) && documentCounter < documentsPerTopic; i++ {
			doc := documentGroup.documentsWithScore[i].Document
			u := urlparser.MustParseURL(doc.URL)
			if _, ok := documentsInEmailByURLIdentifier[u.URLIdentifier]; !ok {
				link, err := makeLinkFromDocument(c, makeLinkFromDocumentInput{
					emailRecordID:   input.emailRecordID,
					userAccessor:    input.userAccessor,
					contentAccessor: input.contentAccessor,
					document:        doc,
				})
				switch {
				case err != nil:
					return nil, err
				case link == nil:
					continue
				}
				documentCounter++
				documentsInEmailByURLIdentifier[u.URLIdentifier] = true
				links = append(links, *link)
			}
		}
		var podcastLinks []PodcastLink
		if len(links) > 0 {
			switch {
			case input.podcastEpisodesByTopic == nil,
				len(input.podcastEpisodesByTopic[documentGroup.topic]) == 0:
				// no-op
			default:
				for _, episode := range input.podcastEpisodesByTopic[documentGroup.topic] {
					if _, ok := podcastsInEmailByID[episode.ID]; ok {
						continue
					}
					podcastLink, err := makeLinkFromPodcast(c, input.podcastAccessor, input.contentAccessor, episode, input.emailRecordID)
					switch {
					case err != nil:
						return nil, err
					case podcastLink == nil:
						// no-op
					}
					podcastsInEmailByID[episode.ID] = true
					podcastLinks = append(podcastLinks, *podcastLink)
					if len(podcastLinks) >= maxPodcastsPerTopic {
						break
					}
				}
			}
			if len(links) > podcastArticleRemovalBreakpoint {
				links = append([]Link{}, links[:len(links)-len(podcastLinks)]...)
			}
			var categoryName *string
			displayName, err := input.contentAccessor.GetDisplayNameByTopicID(documentGroup.topic)
			if err != nil {
				c.Errorf("Error generating display name: %s", err.Error())
			} else {
				categoryName = ptr.String(text.ToTitleCaseForLanguage(*displayName, input.languageCode))
			}
			categories = append(categories, Category{
				topicID:      documentGroup.topic.Ptr(),
				Name:         categoryName,
				Links:        links,
				PodcastLinks: podcastLinks,
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
				link, err := makeLinkFromDocument(c, makeLinkFromDocumentInput{
					emailRecordID:   input.emailRecordID,
					userAccessor:    input.userAccessor,
					contentAccessor: input.contentAccessor,
					document:        doc,
				})
				switch {
				case err != nil:
					return nil, err
				case link == nil:
					continue
				}
				documentCounter++
				documentsInEmailByURLIdentifier[u.URLIdentifier] = true
				links = append(links, *link)
			}
		}
		var podcastLinks []PodcastLink
		for _, podcasts := range input.podcastEpisodesByTopic {
			switch {
			case input.podcastEpisodesByTopic == nil,
				len(podcasts) == 0:
				// no-op
			default:
				for _, episode := range podcasts {
					if _, ok := podcastsInEmailByID[episode.ID]; ok {
						continue
					}
					podcastLink, err := makeLinkFromPodcast(c, input.podcastAccessor, input.contentAccessor, episode, input.emailRecordID)
					switch {
					case err != nil:
						return nil, err
					case podcastLink == nil:
						// no-op
					}
					podcastsInEmailByID[episode.ID] = true
					podcastLinks = append(podcastLinks, *podcastLink)
					if len(podcastLinks) >= maxPodcastsPerTopic {
						break
					}
				}
			}
		}
		if len(links) > podcastArticleRemovalBreakpoint {
			links = append([]Link{}, links[:len(links)-len(podcastLinks)]...)
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

func getTopicsForNewsletter(accessor userPreferencesAccessor) []content.TopicID {
	userSubscriptionLevel := accessor.getUserSubscriptionLevel()
	allUserTopics := accessor.getUserTopics()
	var topics []content.TopicID
	for _, idx := range pickUpToNRandomIndices(int(len(allUserTopics)), defaultNumberOfTopicsPerEmail) {
		topics = append(topics, allUserTopics[idx])
	}
	switch {
	case userSubscriptionLevel == nil:
		return topics
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		topics := accessor.getUserNewsletterSchedule().GetContentTopicsForDay()
		if len(topics) < defaultNumberOfTopicsPerEmail {
			for _, idx := range pickUpToNRandomIndices(int(len(allUserTopics)), defaultNumberOfTopicsPerEmail-len(topics)) {
				topics = append(topics, allUserTopics[idx])
			}
		}
		return topics
	default:
		panic("Unreachable")
	}
}
