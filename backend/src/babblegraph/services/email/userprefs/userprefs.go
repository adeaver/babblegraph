package userprefs

import (
	"babblegraph/model/documents"
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
		userInfos = append(userInfos, UserEmailInfo{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
			// TODO: don't hardcode these
			Languages: []wordsmith.LanguageCode{wordsmith.LanguageCodeSpanish},
			ReadingLevel: UserReadingLevel{
				LowerBound: 60,
				UpperBound: 70,
			},
		})
	}
	return userInfos, nil
}
