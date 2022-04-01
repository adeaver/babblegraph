package usernewsletterpreferences

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"time"
)

type TestNewsletterSchedule struct {
	SendRequested     bool
	UserSendTime      time.Time
	ContentTopics     []contenttopics.ContentTopic
	TopicIDs          []content.TopicID
	NumberOfDocuments int
}

func (t TestNewsletterSchedule) IsSendRequested(utcWeekday time.Weekday) bool {
	return t.SendRequested
}

func (t TestNewsletterSchedule) GetUTCSendTime(utcWeekday time.Weekday) time.Time {
	return t.UserSendTime.UTC()
}

func (t TestNewsletterSchedule) GetNumberOfDocuments(utcWeekday time.Weekday) int {
	return t.NumberOfDocuments
}

func (t TestNewsletterSchedule) ConvertUTCTimeToUserDate(c ctx.LogContext, utcTime time.Time) (*time.Time, error) {
	return ptr.Time(t.UserSendTime), nil
}
