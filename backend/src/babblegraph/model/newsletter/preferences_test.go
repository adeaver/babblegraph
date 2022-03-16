package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlemma"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"time"
)

type testUserAccessor struct {
	userID                    users.UserID
	languageCode              wordsmith.LanguageCode
	doesUserHaveAccount       bool
	userCreatedDate           time.Time
	userSubscriptionLevel     *useraccounts.SubscriptionLevel
	userNewsletterPreferences *usernewsletterpreferences.UserNewsletterPreferences
	userNewsletterSchedule    usernewsletterschedule.UserNewsletterSchedule
	readingLevel              *userReadingLevel
	sentDocumentIDs           []documents.DocumentID
	userTopics                []content.TopicID
	trackingLemmas            []wordsmith.LemmaID
	allowableSourceIDs        []content.SourceID
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

func (t *testUserAccessor) getUserNewsletterSchedule() usernewsletterschedule.UserNewsletterSchedule {
	return t.userNewsletterSchedule
}

func (t *testUserAccessor) getReadingLevel() *userReadingLevel {
	return t.readingLevel
}

func (t *testUserAccessor) getSentDocumentIDs() []documents.DocumentID {
	return t.sentDocumentIDs
}

func (t *testUserAccessor) getUserTopics() []content.TopicID {
	return t.userTopics
}

func (t *testUserAccessor) getTrackingLemmas() []wordsmith.LemmaID {
	return t.trackingLemmas
}

func (t *testUserAccessor) getAllowableSources() []content.SourceID {
	return t.allowableSourceIDs
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

func (t *testUserAccessor) getUserCreatedDate() time.Time {
	return t.userCreatedDate
}
