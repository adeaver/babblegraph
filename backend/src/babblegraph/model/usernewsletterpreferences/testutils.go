package usernewsletterpreferences

import (
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"time"
)

type TestNewsletterSchedule struct {
	SendRequested     bool
	UserSendTime      time.Time
	NumberOfDocuments int
}

func (t TestNewsletterSchedule) IsSendRequested(utcWeekday time.Weekday) bool {
	return t.SendRequested
}

func (t TestNewsletterSchedule) GetUTCHourAndQuarterHourIndex(utcWeekday time.Weekday) (_hourIndex, _quarterHourIndex int) {
	sendTime := t.UserSendTime.UTC()
	return sendTime.Hour(), sendTime.Minute() / 15
}

func (t TestNewsletterSchedule) GetNumberOfDocuments() int {
	return t.NumberOfDocuments
}

func (t TestNewsletterSchedule) ConvertUTCTimeToUserDate(c ctx.LogContext, utcTime time.Time) (*time.Time, error) {
	return ptr.Time(t.UserSendTime), nil
}
