package usernewsletterschedule

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/wordsmith"
	"time"
)

type newsletterScheduleID string

type dbUserNewsletterSchedule struct {
	ID                   newsletterScheduleID   `db:"_id"`
	UserID               users.UserID           `db:"user_id"`
	LanguageCode         wordsmith.LanguageCode `db:"language_code"`
	IANATimezone         string                 `db:"iana_timezone"`
	SendHourIndex        int                    `db:"send_hour_index"`
	SendQuarterHourIndex int                    `db:"send_quarter_hour_index"`
}

type UserNewsletterSchedule struct{}

func GetUserNewsletterScheduleForUTCMidnight(c ctx.LogContext, dayAtUTCMidnight time.Time) (*UserNewsletterSchedule, error) {
	return nil, nil
}

func (u UserNewsletterSchedule) IsSendRequested() bool {
	return true
}

func (u UserNewsletterSchedule) GetSendUTCSendTime() *time.Time {
	return nil
}

func (u UserNewsletterSchedule) GetContentTopicsForDay() []contenttopics.ContentTopic {
	return nil
}

func (u UserNewsletterSchedule) GetNumberOfDocuments() int64 {
	return 0
}
