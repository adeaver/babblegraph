package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type userEmailInfo struct {
	UserID        users.UserID
	EmailAddress  string
	ReadingLevel  userReadingLevel
	Languages     []wordsmith.LanguageCode
	SentDocuments []documents.DocumentID
}

type userReadingLevel struct {
	LowerBound int64
	UpperBound int64
}

func getPreferencesForUser(tx *sqlx.Tx, user users.User) (*userEmailInfo, error) {
	readingLevel, err := userreadability.GetReadabilityScoreRangeForUser(tx, userreadability.GetReadabilityScoreRangeForUserInput{
		UserID:       user.ID,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	})
	if err != nil {
		return nil, err
	}
	sentDocumentIDs, err := userdocuments.GetDocumentIDsSentToUser(tx, user.ID)
	if err != nil {
		return nil, err
	}
	return &userEmailInfo{
		UserID:       user.ID,
		EmailAddress: user.EmailAddress,
		ReadingLevel: userReadingLevel{
			LowerBound: readingLevel.MinScore.ToInt64Rounded(),
			UpperBound: readingLevel.MaxScore.ToInt64Rounded(),
		},
		Languages:     []wordsmith.LanguageCode{wordsmith.LanguageCodeSpanish},
		SentDocuments: sentDocumentIDs,
	}, nil
}
