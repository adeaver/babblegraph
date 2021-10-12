package newslettersendrequests

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"strconv"
	"time"
)

func getDateOfSendForTime(t time.Time) string {
	utcTime := t.UTC()
	return fmt.Sprintf("%02d%02d%04d", utcTime.Month(), utcTime.Day(), utcTime.Year())
}

func makeSendRequestID(userID users.UserID, languageCode wordsmith.LanguageCode, dateOfSendString string) ID {
	return ID(fmt.Sprintf("%s-%s-%s", dateOfSendString, userID, languageCode))
}

func getUTCMidnightDateOfSend(dateOfSendString string) (*time.Time, error) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	year, err := strconv.Atoi(dateOfSendString[4:])
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(dateOfSendString[2:4])
	if err != nil {
		return nil, err
	}
	month, err := strconv.Atoi(dateOfSendString[:2])
	if err != nil {
		return nil, err
	}
	utcDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, utc)
	return &utcDate, nil
}
