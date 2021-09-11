package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlemma"
	"babblegraph/model/userlinks"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"time"

	"github.com/jmoiron/sqlx"
)

type userReadingLevel struct {
	LowerBound int64
	UpperBound int64
}

type userPreferencesAccessor interface {
	getUserID() users.UserID
	getLanguageCode() wordsmith.LanguageCode

	getDoesUserHaveAccount() bool

	getUserSubscriptionLevel() *useraccounts.SubscriptionLevel
	getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences
	getUserScheduleForDay() *usernewsletterschedule.UserNewsletterScheduleDayMetadata
	getReadingLevel() *userReadingLevel
	getSentDocumentIDs() []documents.DocumentID
	getUserTopics() []contenttopics.ContentTopic
	getTrackingLemmas() []wordsmith.LemmaID
	getUserDomainCounts() []userlinks.UserDomainCount
	getSpotlightRecordsOrderedBySentOn() []userlemma.UserLemmaReinforcementSpotlightRecord

	insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error)
	insertSpotlightReinforcementRecord(lemmaID wordsmith.LemmaID) error
}

type DefaultUserPreferencesAccessor struct {
	tx *sqlx.Tx

	userID       users.UserID
	languageCode wordsmith.LanguageCode

	doesUserHaveAccount       bool
	userSubscriptionLevel     *useraccounts.SubscriptionLevel
	userNewsletterPreferences *usernewsletterpreferences.UserNewsletterPreferences
	userScheduleForDay        *usernewsletterschedule.UserNewsletterScheduleDayMetadata
	userReadingLevel          *userReadingLevel
	sentDocumentIDs           []documents.DocumentID
	userTopics                []contenttopics.ContentTopic
	trackingLemmas            []wordsmith.LemmaID
	userDomainCounts          []userlinks.UserDomainCount
	userSpotlightRecords      []userlemma.UserLemmaReinforcementSpotlightRecord
}

func GetDefaultUserPreferencesAccessor(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, dateOfSend time.Time) (*DefaultUserPreferencesAccessor, error) {
	userSubscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	doesUserHaveAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, userID)
	if err != nil {
		return nil, err
	}
	userNewsletterPreferences, err := usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	currentUTCWeekdayIndex := int(dateOfSend.Weekday())
	userScheduleForDay, err := usernewsletterschedule.LookupNewsletterDayMetadataForUserAndDay(tx, userID, currentUTCWeekdayIndex)
	if err != nil {
		return nil, err
	}
	readingLevel, err := userreadability.GetReadabilityScoreRangeForUser(tx, userreadability.GetReadabilityScoreRangeForUserInput{
		UserID:       userID,
		LanguageCode: languageCode,
	})
	if err != nil {
		return nil, err
	}
	sentDocumentIDs, err := userdocuments.GetDocumentIDsSentToUser(tx, userID)
	if err != nil {
		return nil, err
	}
	userTopics, err := usercontenttopics.GetContentTopicsForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	lemmaMappings, err := userlemma.GetVisibleMappingsForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	var trackingLemmas []wordsmith.LemmaID
	for _, m := range lemmaMappings {
		if m.IsActive && m.LanguageCode == languageCode {
			trackingLemmas = append(trackingLemmas, m.LemmaID)
		}
	}
	userDomainCounts, err := userlinks.GetDomainCountsByCurrentAccessMonthForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	userSpotlightRecords, err := userlemma.GetLemmaReinforcementRecordsForUserOrderedBySentOn(tx, userID)
	if err != nil {
		return nil, err
	}
	return &DefaultUserPreferencesAccessor{
		tx:                        tx,
		userID:                    userID,
		languageCode:              languageCode,
		doesUserHaveAccount:       doesUserHaveAccount,
		userSubscriptionLevel:     userSubscriptionLevel,
		userNewsletterPreferences: userNewsletterPreferences,
		userScheduleForDay:        userScheduleForDay,
		userReadingLevel: &userReadingLevel{
			LowerBound: readingLevel.MinScore.ToInt64Rounded(),
			UpperBound: readingLevel.MaxScore.ToInt64Rounded(),
		},
		sentDocumentIDs:      sentDocumentIDs,
		userTopics:           userTopics,
		trackingLemmas:       trackingLemmas,
		userDomainCounts:     userDomainCounts,
		userSpotlightRecords: userSpotlightRecords,
	}, nil
}

func (d *DefaultUserPreferencesAccessor) getUserID() users.UserID {
	return d.userID
}

func (d *DefaultUserPreferencesAccessor) getLanguageCode() wordsmith.LanguageCode {
	return d.languageCode
}

func (d *DefaultUserPreferencesAccessor) getDoesUserHaveAccount() bool {
	return d.doesUserHaveAccount
}

func (d *DefaultUserPreferencesAccessor) getUserSubscriptionLevel() *useraccounts.SubscriptionLevel {
	return d.userSubscriptionLevel
}

func (d *DefaultUserPreferencesAccessor) getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences {
	return d.userNewsletterPreferences
}

func (d *DefaultUserPreferencesAccessor) getUserScheduleForDay() *usernewsletterschedule.UserNewsletterScheduleDayMetadata {
	return d.userScheduleForDay
}

func (d *DefaultUserPreferencesAccessor) getReadingLevel() *userReadingLevel {
	return d.userReadingLevel
}

func (d *DefaultUserPreferencesAccessor) getSentDocumentIDs() []documents.DocumentID {
	return d.sentDocumentIDs
}

func (d *DefaultUserPreferencesAccessor) getUserTopics() []contenttopics.ContentTopic {
	return d.userTopics
}

func (d *DefaultUserPreferencesAccessor) getTrackingLemmas() []wordsmith.LemmaID {
	return d.trackingLemmas
}

func (d *DefaultUserPreferencesAccessor) getUserDomainCounts() []userlinks.UserDomainCount {
	return d.userDomainCounts
}

func (d *DefaultUserPreferencesAccessor) getSpotlightRecordsOrderedBySentOn() []userlemma.UserLemmaReinforcementSpotlightRecord {
	return d.userSpotlightRecords
}

func (d *DefaultUserPreferencesAccessor) insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error) {
	return userdocuments.InsertDocumentForUserAndReturnID(d.tx, d.userID, emailRecordID, doc)
}

func (d *DefaultUserPreferencesAccessor) insertSpotlightReinforcementRecord(lemmaID wordsmith.LemmaID) error {
	return userlemma.UpsertLemmaReinforcementSpotlightRecord(d.tx, userlemma.UpsertLemmaReinforcementSpotlightRecordInput{
		UserID:       d.userID,
		LanguageCode: d.languageCode,
		LemmaID:      lemmaID,
	})
}
