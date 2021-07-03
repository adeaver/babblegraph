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
	// This is technically not right. We send on US/Eastern.
	// However, we are only calculating the day index - and unit tests might be flaky if we use US/Eastern which takes into account Daylight Savings time
	// When we need the exact hour, the job should transition to Etc/GMT+5
	emailSendJobTimezone, err = time.LoadLocation("Etc/GMT+5")
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
	// Since we want to send as close to 10am on in the target time zone as possible
	// We need to figure out when that is.
	requestTimeInOriginalTimezone = time.Date(requestTimeInOriginalTimezone.Year(), requestTimeInOriginalTimezone.Month(), requestTimeInOriginalTimezone.Day(), hourOfTargetSend, 0, 0, 0, requestTimezone)
	// Now we want to convert that time to the timezone that the send job uses (in our case Eastern Standard)
	requestDayIndexAsEmailSendJobTimezone := requestTimeInOriginalTimezone.In(emailSendJobTimezone).Day()
	// We want to figure out if a day before or a day after would be closer to 10am in the client's timezone.
	// To accomplish this, we construct a date object with 10am in our send job timezone
	requestDayInSendJobTimezone := time.Date(requestTimeInOriginalTimezone.Year(), requestTimeInOriginalTimezone.Month(), requestDayIndexAsEmailSendJobTimezone, hourOfTargetSend, 0, 0, 0, emailSendJobTimezone)
	// Now we calculate the difference between that and the client's timezone.
	// If it's greater than 12 hours, then we'll want to send on another day
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

func ConvertIndexedTimeUTCToUserTimezone(dayIndexUTC, hourIndexUTC, quarterHourIndexUTC int, toIANATimezone string) (_dayIndexClientTimezone, _hourIndexClientTimezone, _quarterHourIndexClientTimezone *int, _err error) {
	requestTimezone, err := time.LoadLocation(toIANATimezone)
	if err != nil {
		return nil, nil, nil, err
	}
	indexedTimeUTC := time.Now().UTC()
	for int(indexedTimeUTC.Weekday()) != dayIndexUTC {
		indexedTimeUTC = indexedTimeUTC.Add(24 * time.Hour)
	}
	indexedTimeUTC = time.Date(indexedTimeUTC.Year(), indexedTimeUTC.Month(), indexedTimeUTC.Day(), hourIndexUTC, 15*quarterHourIndexUTC, 0, 0, indexedTimeUTC.Location())
	requestIndexTime := indexedTimeUTC.In(requestTimezone)
	dayIndex := int(requestIndexTime.Weekday())
	hourIndex := requestIndexTime.Hour()
	quarterHourIndex := requestIndexTime.Minute() / 15
	return ptr.Int(dayIndex), ptr.Int(hourIndex), ptr.Int(quarterHourIndex), nil
}
