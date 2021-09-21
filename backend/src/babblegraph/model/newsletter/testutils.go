package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlemma"
	"babblegraph/model/userlinks"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"strings"
)

type testDocsAccessor struct {
	documents []documents.DocumentWithScore
}

func isIDExcluded(id documents.DocumentID, excludedIDs []documents.DocumentID) bool {
	for _, i := range excludedIDs {
		if i == id {
			return true
		}
	}
	return false
}

func isDomainValid(domain string, validDomains []string) bool {
	for _, d := range validDomains {
		if d == domain {
			return true
		}
	}
	return false
}

func containsTopic(topic contenttopics.ContentTopic, topics []contenttopics.ContentTopic) bool {
	for _, t := range topics {
		if t == topic {
			return true
		}
	}
	return false
}

func containsLemma(lemma wordsmith.LemmaID, description string) bool {
	tokens := text.TokenizeUnique(description)
	for _, t := range tokens {
		if t == lemma.Str() {
			return true
		}
	}
	return false
}

func (t *testDocsAccessor) GetDocumentsForUser(input getDocumentsForUserInput) ([]documents.DocumentWithScore, error) {
	var docs []documents.DocumentWithScore
	for _, docWithScore := range t.documents {
		doc := docWithScore.Document
		switch {
		case doc.LanguageCode != input.LanguageCode:
			log.Println(fmt.Sprintf("Language code is not valid"))
		case isIDExcluded(doc.ID, input.ExcludedDocumentIDs):
			log.Println(fmt.Sprintf("ID %s is excluded", doc.ID))
		case !isDomainValid(doc.Domain, input.ValidDomains):
			log.Println(fmt.Sprintf("Domain %s is invalid", doc.Domain))
		case input.MinimumReadingLevel != nil && *input.MinimumReadingLevel > doc.ReadabilityScore:
			log.Println(fmt.Sprintf("Reading level %d is too low", doc.ReadabilityScore))
		case input.MaximumReadingLevel != nil && *input.MaximumReadingLevel < doc.ReadabilityScore:
			log.Println(fmt.Sprintf("Reading level %d is too high", doc.ReadabilityScore))
		case input.Topic != nil && !containsTopic(*input.Topic, doc.Topics):
			log.Println(fmt.Sprintf("Document does not contain topic %s", input.Topic.Str()))
		default:
			docs = append(docs, docWithScore)
		}
	}
	return docs, nil
}

func (t *testDocsAccessor) GetDocumentsForUserForLemma(input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error) {
	var docs []documents.DocumentWithScore
	for _, docWithScore := range t.documents {
		doc := docWithScore.Document
		switch {
		case doc.LanguageCode != input.LanguageCode,
			isIDExcluded(doc.ID, input.ExcludedDocumentIDs),
			isDomainValid(doc.Domain, input.ValidDomains),
			input.MinimumReadingLevel != nil && *input.MinimumReadingLevel < doc.ReadabilityScore,
			input.MaximumReadingLevel != nil && *input.MaximumReadingLevel > doc.ReadabilityScore,
			doc.LemmatizedDescription == nil || !containsLemma(input.Lemma, *doc.LemmatizedDescription):
			continue
		}
		docs = append(docs, docWithScore)
	}
	return docs, nil
}

type testUserAccessor struct {
	userID                    users.UserID
	languageCode              wordsmith.LanguageCode
	doesUserHaveAccount       bool
	userSubscriptionLevel     *useraccounts.SubscriptionLevel
	userNewsletterPreferences *usernewsletterpreferences.UserNewsletterPreferences
	userScheduleForDay        *usernewsletterschedule.UserNewsletterScheduleDayMetadata
	readingLevel              *userReadingLevel
	sentDocumentIDs           []documents.DocumentID
	userTopics                []contenttopics.ContentTopic
	trackingLemmas            []wordsmith.LemmaID
	userDomainCount           []userlinks.UserDomainCount
	spotlightRecords          []userlemma.UserLemmaReinforcementSpotlightRecord

	insertedDocuments        []documents.Document
	insertedSpotlightRecords []wordsmith.LemmaID
}

func (t *testUserAccessor) getUserID() users.UserID {
	return t.userID
}

func (t *testUserAccessor) getLanguageCode() wordsmith.LanguageCode {
	return t.languageCode
}

func (t *testUserAccessor) getDoesUserHaveAccount() bool {
	return t.doesUserHaveAccount
}

func (t *testUserAccessor) getUserSubscriptionLevel() *useraccounts.SubscriptionLevel {
	return t.userSubscriptionLevel
}

func (t *testUserAccessor) getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences {
	return t.userNewsletterPreferences
}

func (t *testUserAccessor) getUserScheduleForDay() *usernewsletterschedule.UserNewsletterScheduleDayMetadata {
	return t.userScheduleForDay
}

func (t *testUserAccessor) getReadingLevel() *userReadingLevel {
	return t.readingLevel
}

func (t *testUserAccessor) getSentDocumentIDs() []documents.DocumentID {
	return t.sentDocumentIDs
}

func (t *testUserAccessor) getUserTopics() []contenttopics.ContentTopic {
	return t.userTopics
}

func (t *testUserAccessor) getTrackingLemmas() []wordsmith.LemmaID {
	return t.trackingLemmas
}

func (t *testUserAccessor) getUserDomainCounts() []userlinks.UserDomainCount {
	return t.userDomainCount
}

func (t *testUserAccessor) getSpotlightRecordsOrderedBySentOn() []userlemma.UserLemmaReinforcementSpotlightRecord {
	return t.spotlightRecords
}

func (t *testUserAccessor) insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error) {
	t.insertedDocuments = append(t.insertedDocuments, doc)
	docID := userdocuments.UserDocumentID(string(doc.ID))
	return &docID, nil
}

func (t *testUserAccessor) insertSpotlightReinforcementRecord(lemmaID wordsmith.LemmaID) error {
	t.insertedSpotlightRecords = append(t.insertedSpotlightRecords, lemmaID)
	return nil
}

type getDefaultDocumentInput struct {
	Topics []contenttopics.ContentTopic
	Lemmas []wordsmith.LemmaID
}

func getDefaultDocumentWithLink(idx int, emailRecordID email.ID, userAccessor userPreferencesAccessor, input getDefaultDocumentInput) (*documents.DocumentWithScore, *Link, error) {
	doc := documents.Document{
		ID:               documents.DocumentID(fmt.Sprintf("web_doc-%d", idx)),
		Version:          documents.Version4,
		URL:              fmt.Sprintf("https://www.elmundo.es/%d", idx),
		ReadabilityScore: 50,
		LanguageCode:     wordsmith.LanguageCodeSpanish,
		DocumentType:     documents.TypeArticle,
		Metadata: documents.Metadata{
			Title:       ptr.String(fmt.Sprintf("Document %d", idx)),
			Image:       ptr.String(fmt.Sprintf("https://www.elmundo.es/%d.jpg", idx)),
			URL:         ptr.String(fmt.Sprintf("https://www.elmundo.es/%d", idx)),
			Description: ptr.String(fmt.Sprintf("This is document #%d", idx)),
		},
		Domain:     "elmundo.es",
		Topics:     input.Topics,
		HasPaywall: ptr.Bool(false),
	}
	link, err := makeLinkFromDocument(emailRecordID, userAccessor, doc)
	if err != nil {
		return nil, nil, err
	}
	return &documents.DocumentWithScore{
		Score:    decimal.FromInt64(1),
		Document: doc,
	}, link, nil
}

func testCategory(expected, result Category) error {
	var errs []string
	matchedLinks := make(map[documents.DocumentID]bool)
	for _, expectedLink := range expected.Links {
		var didFindLink bool
		for _, resultLink := range result.Links {
			isSameLink, err := testLink(expectedLink, resultLink)
			if isSameLink {
				matchedLinks[resultLink.DocumentID] = true
				if err != nil {
					errs = append(errs, err.Error())
				}
				didFindLink = true
				break
			}
		}
		if !didFindLink {
			errs = append(errs, fmt.Sprintf("Expected link for document ID %s, but didn't get it", expectedLink.DocumentID))
		}
	}
	for _, resultLink := range result.Links {
		if _, ok := matchedLinks[resultLink.DocumentID]; !ok {
			errs = append(errs, fmt.Sprintf("Got link for document ID %s, but didn't expect it", resultLink.DocumentID))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

func testLink(expected, result Link) (bool, error) {
	var errs []string
	if expected.DocumentID != result.DocumentID {
		return false, nil
	}
	if err := testutils.CompareNullableString(expected.ImageURL, result.ImageURL); err != nil {
		errs = append(errs, fmt.Sprintf("Image URL for link %s: %s", expected.DocumentID, err.Error()))
	}
	if err := testutils.CompareNullableString(expected.Title, result.Title); err != nil {
		errs = append(errs, fmt.Sprintf("Title for link %s: %s", expected.DocumentID, err.Error()))
	}
	if err := testutils.CompareNullableString(expected.Description, result.Description); err != nil {
		errs = append(errs, fmt.Sprintf("Description for link %s: %s", expected.DocumentID, err.Error()))
	}
	if len(errs) > 0 {
		return true, fmt.Errorf(strings.Join(errs, "\n"))
	}
	return true, nil
}
