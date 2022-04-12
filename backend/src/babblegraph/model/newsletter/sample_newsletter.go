package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userdocuments"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/model/uservocabulary"
	"babblegraph/util/ctx"
	"babblegraph/util/timeutils"
	"babblegraph/wordsmith"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// This is used for generating sample emails locally
type GetSampleNewsletterUserAccessorInput struct {
	UserID       users.UserID
	LanguageCode wordsmith.LanguageCode

	// Overrideable Features
	DoesUserHaveAccount       *bool
	CreatedDate               *time.Time
	UserSubscriptionLevel     *useraccounts.SubscriptionLevel
	UserNewsletterPreferences *usernewsletterpreferences.UserNewsletterPreferences
	UserNewsletterSchedule    *usernewsletterpreferences.Schedule
	SentDocumentIDs           []documents.DocumentID
	UserTopics                []content.TopicID
	UserVocabularyEntries     []uservocabulary.UserVocabularyEntry
	AllowableSourceIDs        []content.SourceID
	SpotlightRecords          []uservocabulary.UserVocabularySpotlightRecord
}

func GetSampleNewsletterUserAccessor(c ctx.LogContext, tx *sqlx.Tx, input GetSampleNewsletterUserAccessorInput) (*SampleNewsletterUserAccessor, error) {
	utcMidnight := timeutils.ConvertToMidnight(time.Now().UTC())
	defaultUserPreferencesAccessor, err := GetDefaultUserPreferencesAccessor(c, tx, input.UserID, input.LanguageCode, utcMidnight)
	if err != nil {
		return nil, err
	}
	if input.DoesUserHaveAccount != nil {
		defaultUserPreferencesAccessor.doesUserHaveAccount = *input.DoesUserHaveAccount
	}
	if input.UserSubscriptionLevel != nil {
		defaultUserPreferencesAccessor.userSubscriptionLevel = input.UserSubscriptionLevel
	}
	if input.UserNewsletterPreferences != nil {
		defaultUserPreferencesAccessor.userNewsletterPreferences = input.UserNewsletterPreferences
	}
	if input.UserNewsletterSchedule != nil {
		defaultUserPreferencesAccessor.userNewsletterSchedule = *input.UserNewsletterSchedule
	}
	if input.SentDocumentIDs != nil {
		defaultUserPreferencesAccessor.sentDocumentIDs = input.SentDocumentIDs
	}
	if input.UserTopics != nil {
		defaultUserPreferencesAccessor.userTopics = input.UserTopics
	}
	if input.UserVocabularyEntries != nil {
		defaultUserPreferencesAccessor.userVocabularyEntries = input.UserVocabularyEntries
	}
	if input.AllowableSourceIDs != nil {
		defaultUserPreferencesAccessor.allowableSourceIDs = input.AllowableSourceIDs
	}
	if input.SpotlightRecords != nil {
		defaultUserPreferencesAccessor.userSpotlightRecords = input.SpotlightRecords
	}
	if input.CreatedDate != nil {
		defaultUserPreferencesAccessor.userCreatedDate = *input.CreatedDate
	}
	return &SampleNewsletterUserAccessor{
		defaultUserPreferencesAccessor: defaultUserPreferencesAccessor,
	}, nil
}

type SampleNewsletterUserAccessor struct {
	defaultUserPreferencesAccessor *DefaultUserPreferencesAccessor
}

func (s *SampleNewsletterUserAccessor) getUserID() users.UserID {
	return s.defaultUserPreferencesAccessor.getUserID()
}

func (s *SampleNewsletterUserAccessor) getLanguageCode() wordsmith.LanguageCode {
	return s.defaultUserPreferencesAccessor.getLanguageCode()
}

func (s *SampleNewsletterUserAccessor) getUserCreatedDate() time.Time {
	return s.defaultUserPreferencesAccessor.getUserCreatedDate()
}

func (s *SampleNewsletterUserAccessor) getDoesUserHaveAccount() bool {
	return s.defaultUserPreferencesAccessor.getDoesUserHaveAccount()
}

func (s *SampleNewsletterUserAccessor) getUserSubscriptionLevel() *useraccounts.SubscriptionLevel {
	return s.defaultUserPreferencesAccessor.getUserSubscriptionLevel()
}

func (s *SampleNewsletterUserAccessor) getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences {
	return s.defaultUserPreferencesAccessor.getUserNewsletterPreferences()
}

func (s *SampleNewsletterUserAccessor) getUserNewsletterSchedule() usernewsletterpreferences.Schedule {
	return s.defaultUserPreferencesAccessor.getUserNewsletterSchedule()
}

func (s *SampleNewsletterUserAccessor) getReadingLevel() *userReadingLevel {
	return s.defaultUserPreferencesAccessor.getReadingLevel()
}

func (s *SampleNewsletterUserAccessor) getSentDocumentIDs() []documents.DocumentID {
	return s.defaultUserPreferencesAccessor.getSentDocumentIDs()
}

func (s *SampleNewsletterUserAccessor) getUserTopics() []content.TopicID {
	return s.defaultUserPreferencesAccessor.getUserTopics()
}

func (s *SampleNewsletterUserAccessor) getUserVocabularyEntries() []uservocabulary.UserVocabularyEntry {
	return s.defaultUserPreferencesAccessor.getUserVocabularyEntries()
}

func (s *SampleNewsletterUserAccessor) getAllowableSources() []content.SourceID {
	return s.defaultUserPreferencesAccessor.getAllowableSources()
}

func (s *SampleNewsletterUserAccessor) getSpotlightRecordsOrderedBySentOn() []uservocabulary.UserVocabularySpotlightRecord {
	return s.defaultUserPreferencesAccessor.getSpotlightRecordsOrderedBySentOn()
}

func (s *SampleNewsletterUserAccessor) insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error) {
	return userdocuments.UserDocumentID(uuid.New().String()).Ptr(), nil
}

func (s *SampleNewsletterUserAccessor) insertSpotlightReinforcementRecord(lemmaID wordsmith.LemmaID) error {
	return nil
}
