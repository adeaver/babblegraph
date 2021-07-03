package usernewsletterschedule

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	getAllNewsletterScheduleMetadataForUserQuery = "SELECT * FROM user_newsletter_schedule_day_metadata WHERE user_id = $1"
	upsertNewsletterScheduleMetadataQuery        = `INSERT INTO
        user_newsletter_schedule_day_metadata (
            user_id,
            day_of_week_index_utc,
            hour_of_day_index_utc,
            quarter_hour_index_utc,
            language_code,
            content_topics,
            number_of_articles,
            is_active
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        ) ON CONFLICT (
            user_id, language_code, day_of_week_index_utc
        ) DO UPDATE
        SET
            hour_of_day_index_utc=$3,
            quarter_hour_index_utc=$4,
            content_topics=$6,
            number_of_articles=$7,
            is_active=$8`
)

func GetNewsletterDayMetadataForUser(tx *sqlx.Tx, userID users.UserID) ([]UserNewsletterScheduleDayMetadata, error) {
	var matches []dbUserNewsletterDayMetadata
	if err := tx.Select(&matches, getAllNewsletterScheduleMetadataForUserQuery, userID); err != nil {
		return nil, err
	}
	var out []UserNewsletterScheduleDayMetadata
	for _, m := range matches {
		metadata, err := m.ToNonDB()
		if err != nil {
			return nil, err
		}
		out = append(out, *metadata)
	}
	return out, nil
}

type UpsertNewsletterDayMetadataForUserInput struct {
	UserID              users.UserID
	DayOfWeekIndexUTC   int
	HourOfDayIndexUTC   int
	QuarterHourIndexUTC int
	LanguageCode        wordsmith.LanguageCode
	ContentTopics       []string
	NumberOfArticles    int
	IsActive            bool
}

func UpsertNewsletterDayMetadataForUser(tx *sqlx.Tx, input UpsertNewsletterDayMetadataForUserInput) error {
	// Validate that all the topics passed are valid
	for _, topicString := range input.ContentTopics {
		_, err := contenttopics.GetContentTopicForString(topicString)
		if err != nil {
			return err
		}
	}
	switch {
	case input.DayOfWeekIndexUTC < 0 || input.DayOfWeekIndexUTC > 6:
		return fmt.Errorf("Day of week must be between 0 and 6")
	case input.HourOfDayIndexUTC < 0 || input.HourOfDayIndexUTC > 23:
		return fmt.Errorf("Hour of day must be between 0 and 23")
	case input.QuarterHourIndexUTC < 0 || input.QuarterHourIndexUTC > 4:
		return fmt.Errorf("Quarter hour must be between 0 and 3")
	case input.NumberOfArticles < minimumNumberOfArticles || input.NumberOfArticles > maximumNumberOfArticles:
		return fmt.Errorf("Number of articles must be between %d and %d", minimumNumberOfArticles, maximumNumberOfArticles)
	}
	var contentTopicString *string
	if len(input.ContentTopics) > 0 {
		contentTopicString = ptr.String(strings.Join(input.ContentTopics, contentTopicDelimiter))
	}
	if _, err := tx.Exec(upsertNewsletterScheduleMetadataQuery, input.UserID, input.DayOfWeekIndexUTC, input.HourOfDayIndexUTC, input.QuarterHourIndexUTC, input.LanguageCode, contentTopicString, input.NumberOfArticles, input.IsActive); err != nil {
		return err
	}
	return nil
}
