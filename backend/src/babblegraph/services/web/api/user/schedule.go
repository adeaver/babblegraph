package user

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

type handleAddUserNewsletterScheduleRequest struct {
	UserScheduleDayRequests []userScheduleDayRequest `json:"user_schedule_day_requests"`
	LanguageCode            wordsmith.LanguageCode   `json:"language_code"`
	IANATimezone            string                   `json:"iana_timezone"`
}

type userScheduleDayRequest struct {
	DayOfWeekIndex   int      `json:"day_of_week_index"`
	ContentTopics    []string `json:"content_topics"`
	NumberOfArticles int      `json:"number_of_articles"`
	IsActive         bool     `json:"is_active"`
}

type handleAddUserNewsletterScheduleResponse struct {
	Success bool `json:"success"`
}

func handleAddUserNewsletterSchedule(userID users.UserID, body []byte) (interface{}, error) {
	var req handleAddUserNewsletterScheduleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, day := range req.UserScheduleDayRequests {
			dayIndexUTC, hourIndexUTC, quarterHourIndexUTC, err := usernewsletterschedule.GetClosetSendTimeInUTC(day.DayOfWeekIndex, req.IANATimezone)
			if err != nil {
				return err
			}
			if err := usernewsletterschedule.UpsertNewsletterDayMetadataForUser(tx, usernewsletterschedule.UpsertNewsletterDayMetadataForUserInput{
				UserID:              userID,
				DayOfWeekIndexUTC:   *dayIndexUTC,
				HourOfDayIndexUTC:   *hourIndexUTC,
				QuarterHourIndexUTC: *quarterHourIndexUTC,
				LanguageCode:        req.LanguageCode,
				NumberOfArticles:    day.NumberOfArticles,
				ContentTopics:       day.ContentTopics,
				IsActive:            day.IsActive,
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Println(fmt.Sprintf("Got error upserting schedule for user %s: %s", userID, err.Error()))
		return handleAddUserNewsletterScheduleResponse{
			Success: false,
		}, nil
	}
	return handleAddUserNewsletterScheduleResponse{
		Success: true,
	}, nil
}

type handleGetUserNewsletterScheduleRequest struct {
	IANATimezone string `json:"iana_timezone"`
}

type handleGetUserNewsletterScheduleResponse struct {
	ScheduleByLanguageCode []scheduleByLanguageCode `json:"schedule_by_language_code"`
}

type scheduleByLanguageCode struct {
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	ScheduleDays []scheduleDay          `json:"schedule_days"`
}

type scheduleDay struct {
	DayOfWeekIndex   int                          `json:"day_of_week_index"`
	HourOfDayIndex   int                          `json:"hour_of_day_index"`
	QuarterHourIndex int                          `json:"quarter_hour_index"`
	ContentTopics    []contenttopics.ContentTopic `json:"content_topics"`
	NumberOfArticles int                          `json:"number_of_articles"`
	IsActive         bool                         `json:"is_active"`
}

func handleGetUserNewsletterSchedule(userID users.UserID, body []byte) (interface{}, error) {
	var req handleGetUserNewsletterScheduleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var daySchedules []usernewsletterschedule.UserNewsletterScheduleDayMetadata
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		daySchedules, err = usernewsletterschedule.GetNewsletterDayMetadataForUser(tx, userID)
		return err
	}); err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	daySchedulesByLanguageCode := make(map[wordsmith.LanguageCode][]scheduleDay)
	for _, d := range daySchedules {
		requestDayIndex, requestHourIndex, requestQuarterHourIndex, err := usernewsletterschedule.ConvertIndexedTimeUTCToUserTimezone(d.DayOfWeekIndexUTC, d.HourOfDayIndexUTC, d.QuarterHourIndexUTC, req.IANATimezone)
		if err != nil {
			sentry.CaptureException(err)
			return nil, err
		}
		daysForLanguageCode, _ := daySchedulesByLanguageCode[d.LanguageCode]
		daySchedulesByLanguageCode[d.LanguageCode] = append(daysForLanguageCode, scheduleDay{
			DayOfWeekIndex:   *requestDayIndex,
			HourOfDayIndex:   *requestHourIndex,
			QuarterHourIndex: *requestQuarterHourIndex,
			ContentTopics:    d.ContentTopics,
			NumberOfArticles: d.NumberOfArticles,
			IsActive:         d.IsActive,
		})
	}
	var outputSchedules []scheduleByLanguageCode
	for languageCode, scheduleDays := range daySchedulesByLanguageCode {
		outputSchedules = append(outputSchedules, scheduleByLanguageCode{
			LanguageCode: languageCode,
			ScheduleDays: scheduleDays,
		})
	}
	return handleGetUserNewsletterScheduleResponse{
		ScheduleByLanguageCode: outputSchedules,
	}, nil
}
