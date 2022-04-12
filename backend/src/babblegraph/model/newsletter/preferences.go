package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlemma"
	"babblegraph/model/userlinks"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/model/uservocabulary"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
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
	getUserCreatedDate() time.Time

	getUserSubscriptionLevel() *useraccounts.SubscriptionLevel
	getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences
	getUserNewsletterSchedule() usernewsletterpreferences.Schedule
	getReadingLevel() *userReadingLevel
	getSentDocumentIDs() []documents.DocumentID
	getUserTopics() []content.TopicID
	getUserVocabularyEntries() []uservocabulary.UserVocabularyEntry
	getAllowableSources() []content.SourceID
	getSpotlightRecordsOrderedBySentOn() []userlemma.UserLemmaReinforcementSpotlightRecord
	insertDocumentForUserAndReturnID(emailRecordID email.ID, doc documents.Document) (*userdocuments.UserDocumentID, error)
	insertSpotlightReinforcementRecord(lemmaID wordsmith.LemmaID) error
}

type DefaultUserPreferencesAccessor struct {
	tx *sqlx.Tx

	userID       users.UserID
	languageCode wordsmith.LanguageCode

	doesUserHaveAccount       bool
	userCreatedDate           time.Time
	userSubscriptionLevel     *useraccounts.SubscriptionLevel
	userNewsletterPreferences *usernewsletterpreferences.UserNewsletterPreferences
	userNewsletterSchedule    usernewsletterpreferences.Schedule
	userReadingLevel          *userReadingLevel
	sentDocumentIDs           []documents.DocumentID
	userTopics                []content.TopicID
	userVocabularyEntries     []uservocabulary.UserVocabularyEntry
	allowableSourceIDs        []content.SourceID
	userSpotlightRecords      []userlemma.UserLemmaReinforcementSpotlightRecord
}

func GetDefaultUserPreferencesAccessor(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, dateOfSendUTCMidnight time.Time) (*DefaultUserPreferencesAccessor, error) {
	user, err := users.GetUser(tx, userID)
	if err != nil {
		return nil, err
	}
	userSubscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	doesUserHaveAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, userID)
	if err != nil {
		return nil, err
	}
	userNewsletterPreferences, err := usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(c, tx, userID, languageCode, ptr.Time(dateOfSendUTCMidnight))
	if err != nil {
		return nil, err
	}
	sentDocumentIDs, err := userdocuments.GetDocumentIDsSentToUser(tx, userID)
	if err != nil {
		return nil, err
	}
	userTopics, err := usercontenttopics.GetTopicIDsForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	vocabularyEntries, err := uservocabulary.GetUserVocabularyEntries(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	allowableSourceIDs, err := getAllowableSourceIDsForUser(tx, userID)
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
		userCreatedDate:           user.CreatedDate,
		userSubscriptionLevel:     userSubscriptionLevel,
		userNewsletterPreferences: userNewsletterPreferences,
		userNewsletterSchedule:    userNewsletterPreferences.Schedule,
		// The current scoring system is pretty broken. So we just set everyone to use the middle.
		// TODO: create a better scoring system
		userReadingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		sentDocumentIDs:       sentDocumentIDs,
		userTopics:            userTopics,
		userVocabularyEntries: vocabularyEntries,
		allowableSourceIDs:    allowableSourceIDs,
		userSpotlightRecords:  userSpotlightRecords,
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

func (d *DefaultUserPreferencesAccessor) getUserCreatedDate() time.Time {
	return d.userCreatedDate
}

func (d *DefaultUserPreferencesAccessor) getUserSubscriptionLevel() *useraccounts.SubscriptionLevel {
	return d.userSubscriptionLevel
}

func (d *DefaultUserPreferencesAccessor) getUserNewsletterPreferences() *usernewsletterpreferences.UserNewsletterPreferences {
	return d.userNewsletterPreferences
}

func (d *DefaultUserPreferencesAccessor) getUserNewsletterSchedule() usernewsletterpreferences.Schedule {
	return d.userNewsletterSchedule
}

func (d *DefaultUserPreferencesAccessor) getReadingLevel() *userReadingLevel {
	return d.userReadingLevel
}

func (d *DefaultUserPreferencesAccessor) getSentDocumentIDs() []documents.DocumentID {
	return d.sentDocumentIDs
}

func (d *DefaultUserPreferencesAccessor) getUserTopics() []content.TopicID {
	return d.userTopics
}

func (d *DefaultUserPreferencesAccessor) getUserVocabularyEntries() []uservocabulary.UserVocabularyEntry {
	return d.userVocabularyEntries
}

func (d *DefaultUserPreferencesAccessor) getAllowableSources() []content.SourceID {
	return d.allowableSourceIDs
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

func getAllowableSourceIDsForUser(tx *sqlx.Tx, userID users.UserID) ([]content.SourceID, error) {
	allowableSources, err := content.GetAllowableSources(tx)
	if err != nil {
		return nil, err
	}
	currentDomainCounts, err := userlinks.GetDomainCountsByCurrentAccessMonthForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	countsBySourceID := make(map[content.SourceID]int64)
	for _, c := range currentDomainCounts {
		countsBySourceID[c.SourceID] = c.Count
	}
	var out []content.SourceID
	for _, s := range allowableSources {
		count, ok := countsBySourceID[s.ID]
		switch {
		case !ok,
			s.MonthlyAccessLimit == nil,
			count < *s.MonthlyAccessLimit:
			out = append(out, s.ID)
		}
	}
	return out, nil
}
