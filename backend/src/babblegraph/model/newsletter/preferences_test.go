package newsletter

import (
	"babblegraph/model/billing"
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userdocuments"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/model/uservocabulary"
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
	userNewsletterSchedule    usernewsletterpreferences.Schedule
	readingLevel              *userReadingLevel
	sentDocumentIDs           []documents.DocumentID
	userTopics                []content.TopicID
	vocabularyEntries         []uservocabulary.UserVocabularyEntry
	allowableSourceIDs        []content.SourceID
	spotlightRecords          []uservocabulary.UserVocabularySpotlightRecord
	paymentState              *billing.PaymentState

	insertedDocuments        []documents.Document
	insertedSpotlightRecords []uservocabulary.UserVocabularyEntryID
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

func (t *testUserAccessor) getUserNewsletterSchedule() usernewsletterpreferences.Schedule {
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

func (t *testUserAccessor) getUserVocabularyEntries() []uservocabulary.UserVocabularyEntry {
	return t.vocabularyEntries
}

func (t *testUserAccessor) getAllowableSources() []content.SourceID {
	return t.allowableSourceIDs
}

func (t *testUserAccessor) getSpotlightRecordsOrderedBySentOn() []uservocabulary.UserVocabularySpotlightRecord {
	return t.spotlightRecords
}

func (t *testUserAccessor) insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error) {
	t.insertedDocuments = append(t.insertedDocuments, doc)
	docID := userdocuments.UserDocumentID(string(doc.ID))
	return &docID, nil
}

func (t *testUserAccessor) insertSpotlightReinforcementRecord(userVocabularyEntryID uservocabulary.UserVocabularyEntryID) error {
	t.insertedSpotlightRecords = append(t.insertedSpotlightRecords, userVocabularyEntryID)
	return nil
}

func (t *testUserAccessor) getSubscriptionPaymentState() *billing.PaymentState {
	return t.paymentState
}

func (t *testUserAccessor) getUserCreatedDate() time.Time {
	return t.userCreatedDate
}
