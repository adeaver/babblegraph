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
	UserSubscriptionLevel     *useraccounts.SubscriptionLevel
	UserNewsletterPreferences *usernewsletterpreferences.UserNewsletterPreferences
	UserNewsletterSchedule    *usernewsletterschedule.UserNewsletterSchedule
	SentDocumentIDs           []documents.DocumentID
	UserTopics                []contenttopics.ContentTopic
	TrackingLemmas            []wordsmith.LemmaID
	UserDomainCounts          []userlinks.UserDomainCount
	SpotlightRecords          []userlemma.UserLemmaReinforcementSpotlightRecord
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
	if input.TrackingLemmas != nil {
		defaultUserPreferencesAccessor.trackingLemmas = input.TrackingLemmas
	}
	if input.UserDomainCounts != nil {
		defaultUserPreferencesAccessor.userDomainCounts = input.UserDomainCounts
	}
	if input.SpotlightRecords != nil {
		defaultUserPreferencesAccessor.userSpotlightRecords = input.SpotlightRecords
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

func (s *SampleNewsletterUserAccessor) getDoesUserHaveAccount() bool {
	return s.defaultUserPreferencesAccessor.getDoesUserHaveAccount()
}

func (s *SampleNewsletterUserAccessor) getUserSubscriptionLevel() *useraccounts.SubscriptionLevel {
	return s.defaultUserPreferencesAccessor.getUserSubscriptionLevel()
}

func (s *SampleNewsletterUserAccessor) getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences {
	return s.defaultUserPreferencesAccessor.getUserNewsletterPreferences()
}

func (s *SampleNewsletterUserAccessor) getUserNewsletterSchedule() usernewsletterschedule.UserNewsletterSchedule {
	return s.defaultUserPreferencesAccessor.getUserNewsletterSchedule()
}

func (s *SampleNewsletterUserAccessor) getReadingLevel() *userReadingLevel {
	return s.defaultUserPreferencesAccessor.getReadingLevel()
}

func (s *SampleNewsletterUserAccessor) getSentDocumentIDs() []documents.DocumentID {
	return s.defaultUserPreferencesAccessor.getSentDocumentIDs()
}

func (s *SampleNewsletterUserAccessor) getUserTopics() []contenttopics.ContentTopic {
	return s.defaultUserPreferencesAccessor.getUserTopics()
}

func (s *SampleNewsletterUserAccessor) getTrackingLemmas() []wordsmith.LemmaID {
	return s.defaultUserPreferencesAccessor.getTrackingLemmas()
}

func (s *SampleNewsletterUserAccessor) getUserDomainCounts() []userlinks.UserDomainCount {
	return s.defaultUserPreferencesAccessor.getUserDomainCounts()
}

func (s *SampleNewsletterUserAccessor) getSpotlightRecordsOrderedBySentOn() []userlemma.UserLemmaReinforcementSpotlightRecord {
	return s.defaultUserPreferencesAccessor.getSpotlightRecordsOrderedBySentOn()
}

func (s *SampleNewsletterUserAccessor) insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error) {
	return userdocuments.UserDocumentID(uuid.New().String()).Ptr(), nil
}

func (s *SampleNewsletterUserAccessor) insertSpotlightReinforcementRecord(lemmaID wordsmith.LemmaID) error {
	return nil
}
