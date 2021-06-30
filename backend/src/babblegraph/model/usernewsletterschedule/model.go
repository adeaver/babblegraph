package usernewsletterschedule

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"strings"
)

const (
	contentTopicDelimiter   = ";"
	minimumNumberOfArticles = 4
	maximumNumberOfArticles = 12
)

type id string

type dbUserNewsletterDayMetadata struct {
	ID                  id                     `db:"_id"`
	UserID              users.UserID           `db:"user_id"`
	DayOfWeekIndexUTC   int                    `db:"day_of_week_index_utc"`
	HourOfDayIndexUTC   int                    `db:"hour_of_day_index_utc"`
	QuarterHourIndexUTC int                    `db:"quarter_hour_index_utc"`
	LanguageCode        wordsmith.LanguageCode `db:"language_code"`
	ContentTopics       *string                `db:"content_topics"`
	NumberOfArticles    int                    `db:"number_of_articles"`
	IsActive            bool                   `db:"is_active"`
}

type UserNewsletterScheduleDayMetadata struct {
	UserID              users.UserID
	DayOfWeekIndexUTC   int
	HourOfDayIndexUTC   int
	QuarterHourIndexUTC int
	LanguageCode        wordsmith.LanguageCode
	ContentTopics       []contenttopics.ContentTopic
	NumberOfArticles    int
	IsActive            bool
}

func (d dbUserNewsletterDayMetadata) ToNonDB() (*UserNewsletterScheduleDayMetadata, error) {
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
	return &UserNewsletterScheduleDayMetadata{
		UserID:              d.UserID,
		DayOfWeekIndexUTC:   d.DayOfWeekIndexUTC,
		HourOfDayIndexUTC:   d.HourOfDayIndexUTC,
		QuarterHourIndexUTC: d.QuarterHourIndexUTC,
		LanguageCode:        d.LanguageCode,
		ContentTopics:       topics,
		NumberOfArticles:    d.NumberOfArticles,
		IsActive:            d.IsActive,
	}, nil
}
