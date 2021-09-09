package newslettersendrequests

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"time"
)

func getDateOfSendForTime(t time.Time) string {
	utcTime := t.UTC()
	return fmt.Sprintf("%02d%02d%4d", utcTime.Month(), utcTime.Day(), utcTime.Year())
}

func makeSendRequestID(userID users.UserID, languageCode wordsmith.LanguageCode, dateOfSendString string) ID {
	return ID(fmt.Sprintf("%s-%s-%s", dateOfSendString, userID, languageCode))
}
