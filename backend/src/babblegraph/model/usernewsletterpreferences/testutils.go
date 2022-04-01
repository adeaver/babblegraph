package usernewsletterpreferences

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"time"
)

type TestNewsletterSchedule struct {
	SendRequested     bool
	UserSendTime      time.Time
	ContentTopics     []contenttopics.ContentTopic
	TopicIDs          []content.TopicID
	NumberOfDocuments int
}

func (t TestNewsletterSchedule) IsSendRequested() bool {
	return t.SendRequested
}

func (t TestNewsletterSchedule) GetUTCSendTime() time.Time {
	return t.UserSendTime.UTC()
}

func (t TestNewsletterSchedule) GetContentTopicsForDay() []content.TopicID {
	return t.TopicIDs
}

func (t TestNewsletterSchedule) GetNumberOfDocuments() int {
	return t.NumberOfDocuments
}

func (t TestNewsletterSchedule) GetSendTimeInUserTimezone() time.Time {
	return t.UserSendTime
}
