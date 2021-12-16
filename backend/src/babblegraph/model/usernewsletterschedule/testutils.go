package usernewsletterschedule

import (
	"babblegraph/model/contenttopics"
	"time"
)

type TestNewsletterSchedule struct {
	SendRequested     bool
	UTCSendTime       time.Time
	ContentTopics     []contenttopics.ContentTopic
	NumberOfDocuments int
}

func (t TestNewsletterSchedule) IsSendRequested() bool {
	return t.SendRequested
}

func (t TestNewsletterSchedule) GetUTCSendTime() time.Time {
	return t.UTCSendTime
}

func (t TestNewsletterSchedule) GetContentTopicsForDay() []contenttopics.ContentTopic {
	return t.ContentTopics
}

func (t TestNewsletterSchedule) GetNumberOfDocuments() int {
	return t.NumberOfDocuments
}
