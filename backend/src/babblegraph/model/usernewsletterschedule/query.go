package usernewsletterschedule

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getAllNewsletterScheduleMetadataForUserQuery    = "SELECT * FROM user_newsletter_schedule_day_metadata WHERE user_id = $1"
	getNewsletterScheduleMetadataForUserForDayQuery = "SELECT * FROM user_newsletter_schedule_day_metadata WHERE user_id = $1 AND day_of_week_index=$2"
	upsertNewsletterScheduleMetadataQuery           = `INSERT INTO
        user_newsletter_schedule_day_metadata (
            user_id,
            day_of_week_index,
            language_code,
            content_topics,
            number_of_articles,
            is_active
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        ) ON CONFLICT (
            user_id, language_code, day_of_week_index
        ) DO UPDATE
        SET
            content_topics=$4,
            number_of_articles=$5,
            is_active=$6`

	getNewsletterScheduleForUserQuery    = "SELECT * FROM user_newsletter_schedule WHERE user_id = $1 AND language_code = $2"
	upsertNewsletterScheduleForUserQuery = `INSERT INTO
        user_newsletter_schedule (
            user_id,
            language_code,
            iana_timezone,
            hour_of_day_index,
            quarter_hour_index
        ) VALUES (
            $1, $2, $3, $4, $5
        ) ON CONFLICT (
            user_id, language_code
        ) DO UPDATE
        SET
            iana_timezone=$3,
            hour_of_day_index=$4,
            quarter_hour_index=$5`
)

func lookupNewsletterDayMetadataForUserAndDay(tx *sqlx.Tx, userID users.UserID, dayOfWeekIndex int) (*dbUserNewsletterDayMetadata, error) {
	var matches []dbUserNewsletterDayMetadata
	if err := tx.Select(&matches, getNewsletterScheduleMetadataForUserForDayQuery, userID, dayOfWeekIndex); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		return &matches[0], nil
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
	UserID           users.UserID
	DayOfWeekIndex   int
	LanguageCode     wordsmith.LanguageCode
	ContentTopics    []string
	NumberOfArticles int
	IsActive         bool
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
	case input.DayOfWeekIndex < 0 || input.DayOfWeekIndex > 6:
		return fmt.Errorf("Day of week must be between 0 and 6")
	case input.NumberOfArticles < minimumNumberOfArticles || input.NumberOfArticles > maximumNumberOfArticles:
		return fmt.Errorf("Number of articles must be between %d and %d", minimumNumberOfArticles, maximumNumberOfArticles)
	case len(input.ContentTopics) > maximumNumberOfTopics:
		return fmt.Errorf("Number of topics must be less than %d", maximumNumberOfTopics)
	}
	var contentTopicString *string
	if len(topicNamesToInsert) > 0 {
		contentTopicString = ptr.String(strings.Join(topicNamesToInsert, contentTopicDelimiter))
	}
	if _, err := tx.Exec(upsertNewsletterScheduleMetadataQuery, input.UserID, input.DayOfWeekIndex, input.LanguageCode, contentTopicString, input.NumberOfArticles, input.IsActive); err != nil {
		return err
	}
	return nil
}

func lookupUserNewsletterScheduleForUser(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*dbUserNewsletterSchedule, error) {
	var matches []dbUserNewsletterSchedule
	if err := tx.Select(&matches, getNewsletterScheduleForUserQuery, userID, languageCode); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected to find at most one schedule, but got %d", len(matches))
	case len(matches) == 1:
		return &matches[0], nil
	default:
		panic("Unreachable")
	}
}

type UpsertUserNewsletterScheduleInput struct {
	UserID           users.UserID
	LanguageCode     wordsmith.LanguageCode
	IANATimezone     *time.Location
	HourIndex        int
	QuarterHourIndex int
}

func UpsertUserNewsletterSchedule(tx *sqlx.Tx, input UpsertUserNewsletterScheduleInput) error {
	switch {
	case input.HourIndex < 0 || input.HourIndex > 23:
		return fmt.Errorf("Error should be between 0 and 23, but got %d", input.HourIndex)
	case input.QuarterHourIndex < 0 || input.QuarterHourIndex > 3:
		return fmt.Errorf("Error should be between 0 and 3, but got %d", input.QuarterHourIndex)
	}
	if _, err := tx.Exec(upsertNewsletterScheduleForUserQuery, input.UserID, input.LanguageCode, input.IANATimezone.String(), input.HourIndex, input.QuarterHourIndex); err != nil {
		return err
	}
	return nil

}
