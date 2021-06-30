package user

import (
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"log"

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
