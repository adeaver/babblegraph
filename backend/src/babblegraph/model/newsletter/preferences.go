package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlemma"
	"babblegraph/model/userlinks"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type userReadingLevel struct {
	LowerBound int64
	UpperBound int64
}

type userPreferencesAccessor interface {
	getReadingLevel() (*userReadingLevel, error)
	getSentDocumentIDs() ([]documents.DocumentID, error)
	getUserTopics() ([]contenttopics.ContentTopic, error)
	getTrackingLemmas() ([]userlemma.Mapping, error)
	getUserDomainCounts() ([]userlinks.UserDomainCount, error)
}

type DefaultUserPreferencesAccessor struct {
	tx           *sqlx.Tx
	userID       users.UserID
	languageCode wordsmith.LanguageCode
}

func GetDefaultUserPreferencesAccessor(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) *DefaultUserPreferencesAccessor {
	return &DefaultUserPreferencesAccessor{
		tx:           tx,
		userID:       userID,
		languageCode: languageCode,
	}
}

func (d *DefaultUserPreferencesAccessor) getReadingLevel() (*userReadingLevel, error) {
	readingLevel, err := userreadability.GetReadabilityScoreRangeForUser(d.tx, userreadability.GetReadabilityScoreRangeForUserInput{
		UserID:       d.userID,
		LanguageCode: d.languageCode,
	})
	if err != nil {
		return nil, err
	}
	return &userReadingLevel{
		LowerBound: readingLevel.MinScore.ToInt64Rounded(),
		UpperBound: readingLevel.MaxScore.ToInt64Rounded(),
	}, nil
}

func (d *DefaultUserPreferencesAccessor) getSentDocumentIDs() ([]documents.DocumentID, error) {
	return userdocuments.GetDocumentIDsSentToUser(d.tx, d.userID)
}

func (d *DefaultUserPreferencesAccessor) getUserTopics() ([]contenttopics.ContentTopic, error) {
	return usercontenttopics.GetContentTopicsForUser(d.tx, d.userID)
}

func (d *DefaultUserPreferencesAccessor) getTrackingLemmas() ([]userlemma.Mapping, error) {
	return userlemma.GetVisibleMappingsForUser(d.tx, d.userID)
}

func (d *DefaultUserPreferencesAccessor) getUserDomainCounts() ([]userlinks.UserDomainCount, error) {
	return userlinks.GetDomainCountsByCurrentAccessMonthForUser(d.tx, d.userID)
}
