package usernewsletterschedule

import (
	"babblegraph/util/ptr"
	"log"
	"time"
)

var emailSendJobTimezone *time.Location

const (
	// This means we try to send the email as close to 10am as possible
	// The reason that this is 10am and not 5:30am is that it takes a few hours
	// to complete the email send job
	hourOfTargetSend int = 10
)

func init() {
	var err error
	emailSendJobTimezone, err = time.LoadLocation("US/Eastern")
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func GetClosetSendTimeInUTC(dayIndex int, ianaTimezoneString string) (_dayIndexUTC, _hourIndexUTC, _quarterHourIndexUTC *int, _err error) {
	requestTimezone, err := time.LoadLocation(ianaTimezoneString)
	if err != nil {
		return nil, nil, nil, err
	}
	requestTimeInOriginalTimezone := time.Now().In(requestTimezone)
	for int(requestTimeInOriginalTimezone.Weekday()) != dayIndex {
		requestTimeInOriginalTimezone = requestTimeInOriginalTimezone.Add(24 * time.Hour)
	}
	// requestTimeInOriginalTimezone now has the correct day of the week, but we need to convert it
	// to the time we would like to try to send the email
	requestTimeInOriginalTimezone = time.Date(requestTimeInOriginalTimezone.Year(), requestTimeInOriginalTimezone.Month(), requestTimeInOriginalTimezone.Day(), hourOfTargetSend, 0, 0, 0, requestTimezone)
	requestDayIndexAsEmailSendJobTimezone := requestTimeInOriginalTimezone.In(emailSendJobTimezone).Day()
	requestDayInSendJobTimezone := time.Date(requestTimeInOriginalTimezone.Year(), requestTimeInOriginalTimezone.Month(), requestDayIndexAsEmailSendJobTimezone, 10, 0, 0, 0, emailSendJobTimezone)
	differenceBetweenSendJobTimeAndRequestTime := getAbsoluteTimeDifferenceInMilliseconds(requestDayInSendJobTimezone, requestTimeInOriginalTimezone)
	currentLowestDifferenceTime := requestDayInSendJobTimezone
	for i := -1; i < 2; i++ {
		adjustedEmailSendJobTime := requestDayInSendJobTimezone.Add(time.Duration(i) * 24 * time.Hour)
		timeDifferenceToOriginalRequestTime := getAbsoluteTimeDifferenceInMilliseconds(adjustedEmailSendJobTime, requestTimeInOriginalTimezone)
		if timeDifferenceToOriginalRequestTime < differenceBetweenSendJobTimeAndRequestTime {
			differenceBetweenSendJobTimeAndRequestTime = timeDifferenceToOriginalRequestTime
			currentLowestDifferenceTime = adjustedEmailSendJobTime
		}
	}
	lowestDifferenceTimeInUTC := currentLowestDifferenceTime.UTC()
	utcDayIndex := int(lowestDifferenceTimeInUTC.Weekday())
	utcHourIndex := lowestDifferenceTimeInUTC.Hour()
	utcQuarterHourIndex := lowestDifferenceTimeInUTC.Minute() % 15
	return ptr.Int(utcDayIndex), ptr.Int(utcHourIndex), ptr.Int(utcQuarterHourIndex), nil
}

func getAbsoluteTimeDifferenceInMilliseconds(t1, t2 time.Time) int64 {
	diff := t1.Sub(t2).Milliseconds()
	if diff < 0 {
		diff = diff * -1
	}
	return diff
}
