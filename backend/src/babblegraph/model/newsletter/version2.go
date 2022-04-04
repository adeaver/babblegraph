package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/math/int2"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// MODEL

type NewsletterVersion2 struct {
	UserID        users.UserID           `json:"user_id"`
	EmailRecordID email.ID               `json:"email_record_id"`
	LanguageCode  wordsmith.LanguageCode `json:"language_code"`
	Body          NewsletterVersion2Body `json:"body"`
}

type NewsletterVersion2Body struct {
	Sections              []Section              `json:"sections"`
	AdvertisingDisclaimer *AdvertisingDisclaimer `json:"advertising_disclaimer,omitempty"`
}

type AdvertisingDisclaimer struct {
	Text                  string                `json:"text"`
	AdvertisingPolicyLink AdvertisingPolicyLink `json:"advertising_policy_link"`
}

type AdvertisingPolicyLink struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type Section struct {
	Title           string               `json:"title"`
	FocusContent    *SectionFocusContent `json:"focus_content,omitempty"`
	OtherLinks      []SectionLink        `json:"other_links,omitempty"`
	OtherLinksTitle *string              `json:"other_links_title"`
}

type SectionFocusContent struct {
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type SectionLink struct {
	Title    string  `json:"text"`
	URL      string  `json:"url"`
	BodyText *string `json:"body_text"`
}

// Query

type CreateNewsletterVersion2Input struct {
	WordsmithAccessor     wordsmithAccessor
	EmailAccessor         emailAccessor
	UserAccessor          userPreferencesAccessor
	DocsAccessor          documentAccessor
	PodcastAccessor       podcastAccessor
	ContentAccessor       contentAccessor
	AdvertisementAccessor advertisementAccessor
}

func CreateNewsletterVersion2(c ctx.LogContext, dateOfSendMidnightUTC time.Time, input CreateNewsletterVersion2Input) (*NewsletterVersion2, error) {
	emailRecordID := email.NewEmailRecordID()
	if err := input.EmailAccessor.InsertEmailRecord(emailRecordID, input.UserAccessor.getUserID()); err != nil {
		return nil, err
	}
	if !input.UserAccessor.getUserNewsletterSchedule().IsSendRequested(dateOfSendMidnightUTC.Weekday()) {
		return nil, nil
	}
	numberOfDocumentsInNewsletter := input.UserAccessor.getUserNewsletterSchedule().GetNumberOfDocuments()
	documentSections, documentIDs, err := getDocumentSections(c, numberOfDocumentsInNewsletter, getDocumentSectionsInput{
		emailRecordID:   emailRecordID,
		userAccessor:    input.UserAccessor,
		docsAccessor:    input.DocsAccessor,
		contentAccessor: input.ContentAccessor,
	})
	if err != nil {
		return nil, err
	}
	spotlightRecord, err := getSpotlightLemmaForNewsletter(c, getSpotlightLemmaForNewsletterInput{
		emailRecordID:           emailRecordID,
		documentIDsInNewsletter: documentIDs,
		userAccessor:            input.UserAccessor,
		docsAccessor:            input.DocsAccessor,
		contentAccessor:         input.ContentAccessor,
		wordsmithAccessor:       input.WordsmithAccessor,
	})
	if err != nil {
		return nil, err
	}
	var advertisingDisclaimer *AdvertisingDisclaimer
	var out []Section
	out = append(out, documentSections[0])
	documentSections = append([]Section{}, documentSections[1:]...)
	if spotlightRecord != nil {
		out = append(out, Section{
			// TODO: make dynamic
			Title: fmt.Sprintf("Tu vocabulario en las noticias: %s", spotlightRecord.LemmaText),
			FocusContent: &SectionFocusContent{
				Title:       *spotlightRecord.Document.Title,
				ImageURL:    *spotlightRecord.Document.ImageURL,
				Description: *spotlightRecord.Document.Description,
				URL:         spotlightRecord.Document.URL,
			},
		})
	}
	userSubscriptionLevel := input.UserAccessor.getUserSubscriptionLevel()
	switch {
	case userSubscriptionLevel == nil:
		advertisement, err := lookupAdvertisement(c, emailRecordID, input.UserAccessor, input.AdvertisementAccessor)
		if err != nil {
			return nil, err
		}
		if advertisement != nil {
			var otherLinks []SectionLink
			if advertisement.AdditionalAdvertisementLink != nil {
				otherLinks = append(otherLinks, SectionLink{
					Title: *advertisement.AdditionalAdvertisementLink.LinkText,
					URL:   advertisement.AdditionalAdvertisementLink.URL,
				})
			}
			// TODO: make this all dynamic
			otherLinks = append(otherLinks, SectionLink{
				Title:    "Si no quieres ver más anuncios, inscribete a Babblegraph Premium",
				BodyText: ptr.String("Con Babblegraph Premium, no verás anuncios como esto. También, tendrás acceso a herramientas exclusivas: como recibir podcasts en el email."),
				URL:      advertisement.PremiumLink,
			})
			out = append(out, Section{
				Title: "Algo que nos gusta",
				FocusContent: &SectionFocusContent{
					Title:       fmt.Sprintf("%s*", advertisement.Title),
					ImageURL:    advertisement.ImageURL,
					URL:         advertisement.URL,
					Description: advertisement.Description,
				},
				OtherLinks: otherLinks,
			})
			advertisingDisclaimer = &AdvertisingDisclaimer{
				Text: "* Asociarnos con excelentes productos y marcas permite que Babblegraph siga funcionando. Podemos ganar una comisión si compra algo a través de uno de estos enlaces.",
				AdvertisingPolicyLink: AdvertisingPolicyLink{
					Text: "Puedes obtener más información sobre anuncios como estos aquí",
					URL:  advertisement.AdvertisementPolicyLink,
				},
			}
		}
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		podcastSection, err := getPodcastSectionForUser(c, getPodcastSectionForUserInput{
			emailRecordID:   emailRecordID,
			userAccessor:    input.UserAccessor,
			podcastAccessor: input.PodcastAccessor,
			contentAccessor: input.ContentAccessor,
		})
		switch {
		case err != nil:
			return nil, err
		case podcastSection == nil:
			// no-op
		default:
			out = append(out, *podcastSection)
		}
	default:
		return nil, fmt.Errorf("Unrecognized subscription level: %s", *userSubscriptionLevel)
	}
	for _, section := range documentSections {
		out = append(out, section)
	}
	var accountLinks []SectionLink
	var reinforcementLink, setTopicsLink, preferencesLink *string
	if input.UserAccessor.getDoesUserHaveAccount() {
		reinforcementLink = ptr.String(routes.MakeLoginLinkWithReinforcementRedirect())
		setTopicsLink = ptr.String(routes.MakeLoginLinkWithContentTopicsRedirect())
		preferencesLink = ptr.String(routes.MakeLoginLinkWithNewsletterPreferencesRedirect())
	} else {
		reinforcementLink, err = routes.MakeWordReinforcementLink(input.UserAccessor.getUserID())
		if err != nil {
			return nil, err
		}
		setTopicsLink, err = routes.MakeSetTopicsLink(input.UserAccessor.getUserID())
		if err != nil {
			return nil, err
		}
		preferencesLink, err = routes.MakeNewsletterPreferencesLink(input.UserAccessor.getUserID())
		if err != nil {
			return nil, err
		}
	}
	// TODO: make dynamic
	accountLinks = append(accountLinks, SectionLink{
		Title: "¿Has aprendido una palabra nueva? Haz clic aquí para añadirla a tu lista de vocabulario",
		URL:   *reinforcementLink,
	})
	if len(input.UserAccessor.getUserTopics()) == 0 {
		accountLinks = append(accountLinks, SectionLink{
			Title: "Puedes escoger temas interesantes para personalizar el próximo boletín",
			URL:   *setTopicsLink,
		})
	}
	accountLinks = append(accountLinks, SectionLink{
		Title: "¿Estás recibiendo demasiados emails? ¿Los emails tienen demasiadas historias? Puedes cambiar eso aquí.",
		URL:   *preferencesLink,
	})
	out = append(out, Section{
		Title:      "Enlaces para gestionar tu suscripción",
		OtherLinks: accountLinks,
	})
	return &NewsletterVersion2{
		UserID:        input.UserAccessor.getUserID(),
		EmailRecordID: emailRecordID,
		LanguageCode:  input.UserAccessor.getLanguageCode(),
		Body: NewsletterVersion2Body{
			Sections:              out,
			AdvertisingDisclaimer: advertisingDisclaimer,
		},
	}, nil
}

// Documents

const (
	maximumNumberOfDocumentsInSection     int = 3
	minimumNumberOfDocumentsInMainSection int = 2

	maximumNumberOfSections int = 3

	maximumNumberOfPodcastsPerEmail int = 3
)

type getDocumentSectionsInput struct {
	emailRecordID   email.ID
	userAccessor    userPreferencesAccessor
	docsAccessor    documentAccessor
	contentAccessor contentAccessor
}

func getDocumentSections(c ctx.LogContext, numberOfDocumentsInNewsletter int, input getDocumentSectionsInput) ([]Section, []documents.DocumentID, error) {
	topics := getSectionTopicsForUser(input.userAccessor, input.contentAccessor)
	allowableSourceIDs := input.userAccessor.getAllowableSources()
	numberOfArticlesInMainSection := int2.MustMinInt(numberOfDocumentsInNewsletter/2, maximumNumberOfDocumentsInSection)
	mainSectionEligibleTopics := make(map[content.TopicID]bool)
	documentsByTopic := make(map[content.TopicID][]documents.DocumentWithScore)
	for _, t := range topics {
		documentsForTopic, err := input.docsAccessor.GetDocumentsForUser(c, getDocumentsForUserInput{
			getDocumentsBaseInput: getDocumentsBaseInput{
				LanguageCode:        input.userAccessor.getLanguageCode(),
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
			return nil, nil, err
		case len(documentsForTopic.RecentDocuments)+len(documentsForTopic.NonRecentDocuments) == 0:
			continue
		case len(documentsForTopic.RecentDocuments) >= numberOfArticlesInMainSection:
			var hasEligibleFocusDocument bool
			for _, document := range documentsForTopic.RecentDocuments {
				isDocumentFocusContentEligible := isDocumentFocusContentEligible(document.Document)
				hasEligibleFocusDocument = hasEligibleFocusDocument || isDocumentFocusContentEligible
			}
			if hasEligibleFocusDocument {
				mainSectionEligibleTopics[t] = hasEligibleFocusDocument
			}
		}
		documentsByTopic[t] = append(documentsForTopic.RecentDocuments, documentsForTopic.NonRecentDocuments...)
		if len(mainSectionEligibleTopics) >= 1 && len(documentsByTopic) >= maximumNumberOfSections {
			break
		}
	}
	var topicIDs []content.TopicID
	for topicID := range documentsByTopic {
		topicIDs = append(topicIDs, topicID)
	}
	sort.SliceStable(topicIDs, func(i, j int) bool {
		isLeftTopicMainDocumentEligible := mainSectionEligibleTopics[topicIDs[i]]
		isRightTopicMainDocumentEligible := mainSectionEligibleTopics[topicIDs[j]]
		switch {
		case isLeftTopicMainDocumentEligible && !isRightTopicMainDocumentEligible:
			return true
		case !isLeftTopicMainDocumentEligible && isRightTopicMainDocumentEligible:
			return false
		default:
			leftDocumentMaxScore := documentsByTopic[topicIDs[i]][0].Score
			rightDocumentMaxScore := documentsByTopic[topicIDs[j]][0].Score
			return leftDocumentMaxScore.GreaterThan(rightDocumentMaxScore)
		}
	})
	var documentIDs []documents.DocumentID
	documentIDsHashSet := make(map[documents.DocumentID]bool)
	numberOfDocumentsRemainingInNewsletter := numberOfDocumentsInNewsletter
	var out []Section
	for i := 0; i < maximumNumberOfSections; i++ {
		documentsWithScore := documentsByTopic[topicIDs[i]]
		numberOfDocumentsInSection := int2.MustMinInt(numberOfDocumentsRemainingInNewsletter, maximumNumberOfDocumentsInSection)
		if numberOfDocumentsInSection == 0 {
			break
		}
		if len(out) == 0 {
			numberOfDocumentsInSection = numberOfArticlesInMainSection
		}
		var focusContent *SectionFocusContent
		var otherLinks []SectionLink
		for _, document := range documentsWithScore {
			if _, ok := documentIDsHashSet[document.Document.ID]; ok {
				continue
			}
			documentIDsHashSet[document.Document.ID] = true
			documentIDs = append(documentIDs, document.Document.ID)
			isDocumentFocusContentEligible := isDocumentFocusContentEligible(document.Document)
			link, err := makeLinkFromDocument(c, makeLinkFromDocumentInput{
				emailRecordID:   input.emailRecordID,
				userAccessor:    input.userAccessor,
				contentAccessor: input.contentAccessor,
				document:        document.Document,
			})
			switch {
			case err != nil:
				return nil, nil, err
			case link == nil:
				continue
			}
			switch {
			case focusContent == nil && isDocumentFocusContentEligible:
				focusContent = &SectionFocusContent{
					Title:       *link.Title,
					ImageURL:    *link.ImageURL,
					Description: *link.Description,
					URL:         link.URL,
				}
				if len(otherLinks) == numberOfDocumentsInSection {
					otherLinks = append([]SectionLink{}, otherLinks[:numberOfDocumentsInSection-1]...)
				} else {
					numberOfDocumentsRemainingInNewsletter--
				}
			case len(otherLinks) <= numberOfDocumentsInSection:
				var description *string
				if link.Domain != nil {
					description = ptr.String(fmt.Sprintf("por %s", link.Domain.Name))
				}
				otherLinks = append(otherLinks, SectionLink{
					Title:    deref.String(link.Title, document.Document.URL),
					BodyText: description,
					URL:      link.URL,
				})
				numberOfDocumentsRemainingInNewsletter--
			}
			if focusContent != nil && len(otherLinks)+1 == numberOfDocumentsInSection {
				break
			}
		}
		// TODO: make this dynamic
		sectionTitle := "En las noticias"
		displayName, err := input.contentAccessor.GetDisplayNameByTopicID(topicIDs[i])
		if err != nil {
			c.Errorf("Error generating display name: %s", err.Error())
		} else {
			sectionTitle = text.ToTitleCaseForLanguage(*displayName, input.userAccessor.getLanguageCode())
		}
		var otherLinksTitle *string
		if len(otherLinks) > 0 {
			// TODO: make this dynamic
			otherLinksTitle = ptr.String("Otros enlaces")
		}
		out = append(out, Section{
			Title:           sectionTitle,
			FocusContent:    focusContent,
			OtherLinks:      otherLinks,
			OtherLinksTitle: otherLinksTitle,
		})
	}
	return out, documentIDs, nil
}

type getPodcastSectionForUserInput struct {
	emailRecordID   email.ID
	userAccessor    userPreferencesAccessor
	podcastAccessor podcastAccessor
	contentAccessor contentAccessor
}

func getPodcastSectionForUser(c ctx.LogContext, input getPodcastSectionForUserInput) (*Section, error) {
	allUserTopics := input.userAccessor.getUserTopics()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allUserTopics), func(i, j int) { allUserTopics[i], allUserTopics[j] = allUserTopics[j], allUserTopics[i] })
	podcastsByTopic, err := input.podcastAccessor.LookupPodcastEpisodesForTopics(allUserTopics)
	if err != nil {
		return nil, err
	}
	var numberOfPodcasts int
	var otherLinks []SectionLink
	var focusContent *SectionFocusContent
	for numberOfPodcasts < maximumNumberOfPodcastsPerEmail && len(podcastsByTopic) > 0 {
		for topicID, episodes := range podcastsByTopic {
			episode := episodes[0]
			podcastMetadata, err := input.podcastAccessor.GetPodcastMetadataForSourceID(episode.SourceID)
			switch {
			case err != nil:
				c.Errorf("Error getting podcast metadata for podcast %s: %s", episode.SourceID, err.Error())
			case focusContent == nil && podcastMetadata.ImageURL != nil:
				podcastLink, err := makeLinkFromPodcast(c, input.podcastAccessor, input.contentAccessor, episode, input.emailRecordID)
				if err != nil {
					return nil, err
				}
				focusContent = &SectionFocusContent{
					Title:       podcastLink.EpisodeTitle,
					ImageURL:    *podcastLink.PodcastImageURL,
					Description: podcastLink.EpisodeDescription,
					URL:         podcastLink.ListenURL,
				}
				numberOfPodcasts++
			case len(otherLinks) < maximumNumberOfPodcastsPerEmail-1:
				podcastLink, err := makeLinkFromPodcast(c, input.podcastAccessor, input.contentAccessor, episode, input.emailRecordID)
				if err != nil {
					return nil, err
				}
				otherLinks = append(otherLinks, SectionLink{
					Title: podcastLink.EpisodeTitle,
					URL:   podcastLink.ListenURL,
				})
				numberOfPodcasts++
			default:
				// no-op
			}
			nextEpisodes := append([]podcasts.Episode{}, episodes[1:]...)
			podcastsByTopic[topicID] = nextEpisodes
			if len(nextEpisodes) == 0 {
				delete(podcastsByTopic, topicID)
			}
		}
	}
	if len(otherLinks) == 0 && focusContent == nil {
		return nil, nil
	}
	// TODO: make this dynamic
	return &Section{
		Title:           "Podcasts para ti",
		FocusContent:    focusContent,
		OtherLinks:      otherLinks,
		OtherLinksTitle: ptr.String("Otros episodios"),
	}, nil
}

func isDocumentFocusContentEligible(doc documents.Document) bool {
	documentMetadata := doc.Metadata
	return documentMetadata.Image != nil && documentMetadata.Title != nil && documentMetadata.Description != nil
}

func getSectionTopicsForUser(accessor userPreferencesAccessor, contentAccessor contentAccessor) []content.TopicID {
	allUserTopics := accessor.getUserTopics()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allUserTopics), func(i, j int) { allUserTopics[i], allUserTopics[j] = allUserTopics[j], allUserTopics[i] })
	userNonTopics := contentAccessor.GetAllTopicsNotInList(allUserTopics)
	rand.Shuffle(len(userNonTopics), func(i, j int) { userNonTopics[i], userNonTopics[j] = userNonTopics[j], userNonTopics[i] })
	return append(allUserTopics, userNonTopics...)
}
