package userprefs

import (
	"babblegraph/model/documents"
	"babblegraph/model/userreadability"
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type UserEmailInfo struct {
	UserID        users.UserID
	EmailAddress  string
	ReadingLevel  UserReadingLevel
	Languages     []wordsmith.LanguageCode
	SentDocuments []documents.DocumentID
}

type UserReadingLevel struct {
	LowerBound int64
	UpperBound int64
}

func GetActiveUserEmailInfo(tx *sqlx.Tx) ([]UserEmailInfo, error) {
	activeUsers, err := users.GetAllActiveUsers(tx)
	if err != nil {
		return nil, err
	}
	var userInfos []UserEmailInfo
	for _, user := range activeUsers {
		readingLevel, err := userreadability.GetReadabilityScoreRangeForUser(tx, userreadability.GetReadabilityScoreRangeForUserInput{
			UserID:       user.ID,
			LanguageCode: wordsmith.LanguageCodeSpanish,
		})
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, UserEmailInfo{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
			ReadingLevel: UserReadingLevel{
				LowerBound: readingLevel.MinScore.ToInt64Rounded(),
				UpperBound: readingLevel.MaxScore.ToInt64Rounded(),
			},
			Languages: []wordsmith.LanguageCode{wordsmith.LanguageCodeSpanish},
		})
	}
	return userInfos, nil
}
