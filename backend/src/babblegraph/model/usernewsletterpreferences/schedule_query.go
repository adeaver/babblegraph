package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getAllNewsletterScheduleMetadataForUserQuery = "SELECT * FROM user_newsletter_schedule_day_metadata WHERE user_id = $1 AND language_code=$2"
	upsertNewsletterScheduleMetadataQuery        = `INSERT INTO
        user_newsletter_schedule_day_metadata (
            user_id,
            day_of_week_index,
            language_code,
            number_of_articles,
            is_active
        ) VALUES (
            $1, $2, $3, $4, $5
        ) ON CONFLICT (
            user_id, language_code, day_of_week_index
        ) DO UPDATE
        SET
            number_of_articles=$4,
            is_active=$5
        `

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

func lookupNewsletterDayMetadataForUser(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) ([]dbUserNewsletterDayMetadata, error) {
	var matches []dbUserNewsletterDayMetadata
	if err := tx.Select(&matches, getAllNewsletterScheduleMetadataForUserQuery, userID, languageCode); err != nil {
		return nil, err
	}
	return matches, nil
}

type upsertNewsletterDayMetadataForUserInput struct {
	UserID           users.UserID
	DayOfWeekIndex   int
	LanguageCode     wordsmith.LanguageCode
	NumberOfArticles int
	IsActive         bool
}

func upsertNewsletterDayMetadataForUser(tx *sqlx.Tx, input upsertNewsletterDayMetadataForUserInput) error {
	switch {
	case input.DayOfWeekIndex < 0 || input.DayOfWeekIndex > 6:
		return fmt.Errorf("Day of week must be between 0 and 6")
	case input.NumberOfArticles < minimumNumberOfArticles || input.NumberOfArticles > maximumNumberOfArticles:
		return fmt.Errorf("Number of articles must be between %d and %d", minimumNumberOfArticles, maximumNumberOfArticles)
	}
	if _, err := tx.Query(upsertNewsletterScheduleMetadataQuery, input.UserID, input.DayOfWeekIndex, input.LanguageCode, input.NumberOfArticles, input.IsActive); err != nil {
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

type upsertUserNewsletterScheduleInput struct {
	UserID           users.UserID
	LanguageCode     wordsmith.LanguageCode
	IANATimezone     *time.Location
	HourIndex        int
	QuarterHourIndex int
}

func upsertUserNewsletterSchedule(tx *sqlx.Tx, input upsertUserNewsletterScheduleInput) error {
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
