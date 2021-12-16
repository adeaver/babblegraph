package usernewsletterschedule

import (
	"babblegraph/util/ctx"
	"fmt"
	"time"
)

func resolveUTCMidnightWithNewsletterSchedule(c ctx.LogContext, utcMidnight time.Time, userNewsletterSchedule dbUserNewsletterSchedule) (*time.Time, error) {
	if isMidnight := isDateMidnight(utcMidnight); !isMidnight {
		return nil, fmt.Errorf("date is not midnight")
	}
	// This function is a bit convoluted.
	// The entrypoint for this package is a newsletter send request is generated for a specific UTC day
	// such as 12/25/2020. However, a user may select to have a newsletter sent to them at 9pm PST.
	// This means that the newsletter corresponding to UTC 12/25/2020, will actually be sent on 12/24/2020 PST.
	userTimezone, err := time.LoadLocation(userNewsletterSchedule.IANATimezone)
	if err != nil {
		c.Errorf("Error loading timezone %s for user schedule: %s", userNewsletterSchedule.IANATimezone, err.Error())
		return nil, err
	}
	beforeNextMidnight := utcMidnight.Add(24*time.Hour - time.Second)
	minute := userNewsletterSchedule.QuarterHourIndex * 15
	userSendTime := time.Date(utcMidnight.Year(), utcMidnight.Month(), utcMidnight.Day(), userNewsletterSchedule.HourIndex, minute, 0, 0, userTimezone)
	switch {
	case userSendTime.After(beforeNextMidnight):
		nextPossibleSendTime := userSendTime.Add(-24 * time.Hour)
		for nextPossibleSendTime.After(beforeNextMidnight) {
			userSendTime = nextPossibleSendTime
			nextPossibleSendTime = userSendTime.Add(-24 * time.Hour)
		}
		return &nextPossibleSendTime, nil
	case userSendTime.Before(beforeNextMidnight):
		nextPossibleSendTime := userSendTime.Add(24 * time.Hour)
		for nextPossibleSendTime.Before(beforeNextMidnight) {
			userSendTime = nextPossibleSendTime
			nextPossibleSendTime = userSendTime.Add(24 * time.Hour)
		}
	}
	return &userSendTime, nil
}

func isDateMidnight(t time.Time) bool {
	return (t.Hour() == 0 &&
		t.Minute() == 0 &&
		t.Second() == 0 &&
		t.Nanosecond() == 0)
}
