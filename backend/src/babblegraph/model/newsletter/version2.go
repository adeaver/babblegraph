package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/math/int2"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
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
	sections, err := getDocumentSections(c, numberOfDocumentsInNewsletter, getDocumentSectionsInput{
		emailRecordID:   emailRecordID,
		userAccessor:    input.UserAccessor,
		docsAccessor:    input.DocsAccessor,
		contentAccessor: input.ContentAccessor,
	})
	if err != nil {
		return nil, err
	}
	return &NewsletterVersion2{
		UserID:        input.UserAccessor.getUserID(),
		EmailRecordID: emailRecordID,
		LanguageCode:  input.UserAccessor.getLanguageCode(),
		Body: NewsletterVersion2Body{
			Sections: sections,
		},
	}, nil
}

// Documents

const (
	maximumNumberOfDocumentsInSection     int = 3
	minimumNumberOfDocumentsInMainSection int = 2

	maximumNumberOfSections int = 3
)

type getDocumentSectionsInput struct {
	emailRecordID   email.ID
	userAccessor    userPreferencesAccessor
	docsAccessor    documentAccessor
	contentAccessor contentAccessor
}

func getDocumentSections(c ctx.LogContext, numberOfDocumentsInNewsletter int, input getDocumentSectionsInput) ([]Section, error) {
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
			return nil, err
		case len(documentsForTopic.RecentDocuments)+len(documentsForTopic.NonRecentDocuments) == 0:
			continue
		case len(documentsForTopic.RecentDocuments) >= numberOfArticlesInMainSection:
			var hasEligibleFocusDocument bool
			for _, document := range documentsForTopic.RecentDocuments {
				isDocumentFocusContentEligible := isDocumentFocusContentEligible(document.Document)
				hasEligibleFocusDocument = hasEligibleFocusDocument || isDocumentFocusContentEligible
			}
			mainSectionEligibleTopics[t] = hasEligibleFocusDocument
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
	numberOfDocumentsRemainingInNewsletter := numberOfDocumentsInNewsletter
	var out []Section
	for i := 0; i <= maximumNumberOfSections; i++ {
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
			isDocumentFocusContentEligible := isDocumentFocusContentEligible(document.Document)
			link, err := makeLinkFromDocument(c, makeLinkFromDocumentInput{
				emailRecordID:   input.emailRecordID,
				userAccessor:    input.userAccessor,
				contentAccessor: input.contentAccessor,
				document:        document.Document,
			})
			switch {
			case err != nil:
				return nil, err
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
				otherLinks = append(otherLinks, SectionLink{
					Title: deref.String(link.Title, document.Document.URL),
					URL:   link.URL,
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
	return out, nil
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
