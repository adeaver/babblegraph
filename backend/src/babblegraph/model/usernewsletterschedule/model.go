package usernewsletterschedule

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	contentTopicDelimiter   = ";"
	minimumNumberOfArticles = 4
	maximumNumberOfArticles = 12

	maximumNumberOfTopics = 6

	defaultUTCSendTimeHour = 11 * time.Hour
)

type dayID string

type dbUserNewsletterDayMetadata struct {
	ID               dayID                  `db:"_id"`
	UserID           users.UserID           `db:"user_id"`
	LanguageCode     wordsmith.LanguageCode `db:"language_code"`
	DayOfWeekIndex   int                    `db:"day_of_week_index"`
	ContentTopics    *string                `db:"content_topics"`
	NumberOfArticles int                    `db:"number_of_articles"`
	IsActive         bool                   `db:"is_active"`
}

type scheduleDayTopicMappingID string

type dbUserNewsletterScheduleDayTopicMapping struct {
	ID             scheduleDayTopicMappingID `db:"_id"`
	CreatedAt      time.Time                 `db:"created_at"`
	LastModifiedAt time.Time                 `db:"last_modified_at"`
	TopicID        content.TopicID           `db:"topic_id"`
	DayID          dayID                     `db:"day_id"`
	IsActive       bool                      `db:"is_active"`
}

type UserNewsletterScheduleDayMetadata struct {
	UserID           users.UserID
	LanguageCode     wordsmith.LanguageCode
	DayOfWeekIndex   int
	ContentTopics    []contenttopics.ContentTopic
	TopicIDs         []content.TopicID
	NumberOfArticles int
	IsActive         bool
}

func (d dbUserNewsletterDayMetadata) ToNonDB(tx *sqlx.Tx) (*UserNewsletterScheduleDayMetadata, error) {
	var topics []contenttopics.ContentTopic
	if d.ContentTopics != nil {
		topicStrings := strings.Split(*d.ContentTopics, contentTopicDelimiter)
		for _, s := range topicStrings {
			t, err := contenttopics.GetContentTopicForString(s)
			if err != nil {
				return nil, err
			}
			topics = append(topics, *t)
		}
	}
	topicIDs, err := lookupTopicMappingsForDay(tx, d.ID)
	if err != nil {
		return nil, err
	}
	return &UserNewsletterScheduleDayMetadata{
		UserID:           d.UserID,
		DayOfWeekIndex:   d.DayOfWeekIndex,
		LanguageCode:     d.LanguageCode,
		ContentTopics:    topics,
		TopicIDs:         topicIDs,
		NumberOfArticles: d.NumberOfArticles,
		IsActive:         d.IsActive,
	}, nil
}

type id string

type dbUserNewsletterSchedule struct {
	ID               id                     `db:"_id"`
	UserID           users.UserID           `db:"user_id"`
	LanguageCode     wordsmith.LanguageCode `db:"language_code"`
	IANATimezone     string                 `db:"iana_timezone"`
	HourIndex        int                    `db:"hour_of_day_index"`
	QuarterHourIndex int                    `db:"quarter_hour_index"`
}

type UserNewsletterSchedule interface {
	IsSendRequested() bool
	GetUTCSendTime() time.Time
	GetContentTopicsForDay() []content.TopicID
	GetNumberOfDocuments() int
	GetSendTimeInUserTimezone() time.Time
}

type ScheduleWithMetadata struct {
	sendTimeAtUserTimezone time.Time
	userScheduleDay        *UserNewsletterScheduleDayMetadata
}

type GetUserNewsletterScheduleForUTCMidnightInput struct {
	UserID           users.UserID
	LanguageCode     wordsmith.LanguageCode
	DayAtUTCMidnight time.Time
}

func GetUserNewsletterScheduleForUTCMidnight(c ctx.LogContext, tx *sqlx.Tx, input GetUserNewsletterScheduleForUTCMidnightInput) (*ScheduleWithMetadata, error) {
	var userScheduleTime *time.Time
	userSchedule, err := lookupUserNewsletterScheduleForUser(tx, input.UserID, input.LanguageCode)
	switch {
	case err != nil:
		return nil, err
	case userSchedule == nil:
		userScheduleTime = ptr.Time(input.DayAtUTCMidnight.Add(defaultUTCSendTimeHour))
	case userSchedule != nil:
		userScheduleTime, err = resolveUTCMidnightWithNewsletterSchedule(c, input.DayAtUTCMidnight, *userSchedule)
		if err != nil {
			return nil, err
		}
	default:
		panic("Unreachable")
	}
	var userScheduleDay *UserNewsletterScheduleDayMetadata
	dbUserScheduleDay, err := lookupNewsletterDayMetadataForUserAndDay(tx, input.UserID, int(userScheduleTime.Weekday()))
	switch {
	case err != nil:
		return nil, err
	case dbUserScheduleDay == nil:
		// no-op
	case dbUserScheduleDay != nil:
		userScheduleDay, err = dbUserScheduleDay.ToNonDB(tx)
		if err != nil {
			return nil, err
		}
	}
	c.Debugf("Got user send time of %+v for UTC midnight %+v with day of %+v", userScheduleTime, input.DayAtUTCMidnight, userScheduleDay)
	return &ScheduleWithMetadata{
		sendTimeAtUserTimezone: *userScheduleTime,
		userScheduleDay:        userScheduleDay,
	}, nil
}

func (u ScheduleWithMetadata) IsSendRequested() bool {
	return u.userScheduleDay == nil || u.userScheduleDay.IsActive
}

func (u ScheduleWithMetadata) GetUTCSendTime() time.Time {
	return u.sendTimeAtUserTimezone.UTC()
}

func (u ScheduleWithMetadata) GetContentTopicsForDay() []content.TopicID {
	if u.userScheduleDay == nil {
		return nil
	}
	return u.userScheduleDay.TopicIDs
}

func (u ScheduleWithMetadata) GetNumberOfDocuments() int {
	if u.userScheduleDay == nil {
		return maximumNumberOfArticles
	}
	return u.userScheduleDay.NumberOfArticles
}

func (u ScheduleWithMetadata) GetSendTimeInUserTimezone() time.Time {
	return u.sendTimeAtUserTimezone
}
