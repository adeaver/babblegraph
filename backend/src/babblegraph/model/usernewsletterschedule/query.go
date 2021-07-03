package usernewsletterschedule

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	getAllNewsletterScheduleMetadataForUserQuery    = "SELECT * FROM user_newsletter_schedule_day_metadata WHERE user_id = $1"
	getNewsletterScheduleMetadataForUserForDayQuery = "SELECT * FROM user_newsletter_schedule_day_metadata WHERE user_id = $1 AND day_of_week_index_utc=$2"
	upsertNewsletterScheduleMetadataQuery           = `INSERT INTO
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

func LookupNewsletterDayMetadataForUserAndDay(tx *sqlx.Tx, userID users.UserID, dayOfWeekIndexUTC int) (*UserNewsletterScheduleDayMetadata, error) {
	var matches []dbUserNewsletterDayMetadata
	if err := tx.Select(&matches, getNewsletterScheduleMetadataForUserForDayQuery, userID, dayOfWeekIndexUTC); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		m, err := matches[0].ToNonDB()
		if err != nil {
			return nil, err
		}
		return m, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected 0 or 1 record, but got %d", len(matches))
	}
	return nil, nil
}

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
	// Validate that all the topics passed are valid and only occur once
	seenTopicNames := make(map[string]bool)
	var topicNamesToInsert []string
	for _, topicString := range input.ContentTopics {
		if _, ok := seenTopicNames[topicString]; ok {
			continue
		}
		_, err := contenttopics.GetContentTopicForString(topicString)
		if err != nil {
			log.Println(fmt.Sprintf("Got invalid content topic %s", err.Error()))
			continue
		}
		seenTopicNames[topicString] = true
		topicNamesToInsert = append(topicNamesToInsert, topicString)
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
	case len(input.ContentTopics) > maximumNumberOfTopics:
		return fmt.Errorf("Number of topics must be less than %d", maximumNumberOfTopics)
	}
	var contentTopicString *string
	if len(topicNamesToInsert) > 0 {
		contentTopicString = ptr.String(strings.Join(topicNamesToInsert, contentTopicDelimiter))
	}
	if _, err := tx.Exec(upsertNewsletterScheduleMetadataQuery, input.UserID, input.DayOfWeekIndexUTC, input.HourOfDayIndexUTC, input.QuarterHourIndexUTC, input.LanguageCode, contentTopicString, input.NumberOfArticles, input.IsActive); err != nil {
		return err
	}
	return nil
}
