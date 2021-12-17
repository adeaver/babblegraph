package user

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/newsletter"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/timeutils"
	"babblegraph/wordsmith"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type getUserScheduleRequest struct {
	EmailAddress string                 `json:"email_address"`
	Token        string                 `json:"token"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
}

type getUserScheduleResponse struct {
	UserIANATimezone string           `json:"user_iana_timezone"`
	HourIndex        int              `json:"hour_index"`
	QuarterHourIndex int              `json:"quarter_hour_index"`
	PreferencesByDay []dayPreferences `json:"preferences_by_day,omitempty"`
}

type dayPreferences struct {
	IsActive         bool                         `json:"is_active"`
	NumberOfArticles int                          `json:"number_of_articles"`
	ContentTopics    []contenttopics.ContentTopic `json:"content_topics"`
	DayIndex         int                          `json:"day_index"`
}

func handleGetUserSchedule(body []byte) (interface{}, error) {
	var req getUserScheduleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndEmailAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey, req.EmailAddress)
	if err != nil {
		return nil, err
	}
	// TODO(web-context): don't use default log context here
	c := ctx.GetDefaultLogContext()
	var userNewsletterSchedule usernewsletterschedule.UserNewsletterSchedule
	var userDayPreferences []usernewsletterschedule.UserNewsletterScheduleDayMetadata
	var doesUserHaveSubscription bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userNewsletterSchedule, err = usernewsletterschedule.GetUserNewsletterScheduleForUTCMidnight(c, tx, usernewsletterschedule.GetUserNewsletterScheduleForUTCMidnightInput{
			UserID:           *userID,
			DayAtUTCMidnight: timeutils.ConvertToMidnight(time.Now().UTC()),
			LanguageCode:     req.LanguageCode,
		})
		if err != nil {
			return err
		}
		_, doesUserHaveSubscription, err = useraccounts.DoesUserHaveSubscription(tx, *userID)
		switch {
		case err != nil:
			return err
		case !doesUserHaveSubscription:
			return nil
		}
		// TODO(multiple-languages): Modify this to take in a language code
		userDayPreferences, err = usernewsletterschedule.GetNewsletterDayMetadataForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	userSendTime := userNewsletterSchedule.GetSendTimeInUserTimezone()
	var apiDayPreferences []dayPreferences = nil
	if doesUserHaveSubscription {
		for i := 0; i < 7; i++ {
			apiDayPreferences = append(apiDayPreferences, dayPreferences{
				IsActive:         true,
				NumberOfArticles: newsletter.DefaultNumberOfArticlesPerEmail,
				DayIndex:         i,
			})
		}
		for _, userDayPreference := range userDayPreferences {
			apiDayPreferences[userDayPreference.DayOfWeekIndex] = dayPreferences{
				IsActive:         userDayPreference.IsActive,
				NumberOfArticles: userDayPreference.NumberOfArticles,
				ContentTopics:    userDayPreference.ContentTopics,
			}
		}
	}
	return getUserScheduleResponse{
		UserIANATimezone: userSendTime.Location().String(),
		HourIndex:        userSendTime.Hour(),
		QuarterHourIndex: userSendTime.Minute() / 15,
		PreferencesByDay: apiDayPreferences,
	}, nil
}

type updateUserScheduleRequest struct {
	EmailAddress     string `json:"email_address"`
	Token            string `json:"token"`
	LanguageCode     string `json:"language_code"`
	HourIndex        int    `json:"hour_index"`
	QuarterHourIndex int    `json:"quarter_hour_index"`
	IANATimezone     string `json:"iana_timezone"`
}

type updateUserScheduleResponse struct {
	Success bool `json:"success"`
}

func handleUpdateUserSchedule(body []byte) (interface{}, error) {
	var req updateUserScheduleRequest
	if err := json.Unmarshal(body, req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndEmailAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey, req.EmailAddress)
	if err != nil {
		return updateUserScheduleResponse{
			Success: false,
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return updateUserScheduleResponse{
			Success: false,
		}, nil
	}
	loc, err := time.LoadLocation(req.IANATimezone)
	if err != nil {
		return updateUserScheduleResponse{
			Success: false,
		}, nil
	}
	switch {
	case req.HourIndex < 0 || req.HourIndex > 23,
		req.QuarterHourIndex < 0 || req.QuarterHourIndex > 3:
		return updateUserScheduleResponse{
			Success: false,
		}, nil
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return usernewsletterschedule.UpsertUserNewsletterSchedule(tx, usernewsletterschedule.UpsertUserNewsletterScheduleInput{
			UserID:           *userID,
			LanguageCode:     *languageCode,
			IANATimezone:     loc,
			QuarterHourIndex: req.QuarterHourIndex,
			HourIndex:        req.HourIndex,
		})
	}); err != nil {
		return nil, err
	}
	return updateUserScheduleResponse{
		Success: true,
	}, nil
}

type updateUserScheduleWithDayPreferencesRequest struct {
	Token            string           `json:"token"`
	LanguageCode     string           `json:"language_code"`
	HourIndex        int              `json:"hour_index"`
	QuarterHourIndex int              `json:"quarter_hour_index"`
	IANATimezone     string           `json:"iana_timezone"`
	DayPreferences   []dayPreferences `json:"day_preferences"`
}

type updateUserScheduleWithDayPreferencesResponse struct {
	Success bool `json:"success"`
}

func handleUpdateUserScheduleWithDayPreferences(userID users.UserID, body []byte) (interface{}, error) {
	var req updateUserScheduleWithDayPreferencesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	tokenUserID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil || userID != *tokenUserID {
		return updateUserScheduleWithDayPreferencesResponse{
			Success: false,
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return updateUserScheduleWithDayPreferencesResponse{
			Success: false,
		}, nil
	}
	loc, err := time.LoadLocation(req.IANATimezone)
	if err != nil {
		return updateUserScheduleWithDayPreferencesResponse{
			Success: false,
		}, nil
	}
	switch {
	case req.HourIndex < 0 || req.HourIndex > 23,
		req.QuarterHourIndex < 0 || req.QuarterHourIndex > 3:
		return updateUserScheduleWithDayPreferencesResponse{
			Success: false,
		}, nil
	}
	for _, pref := range req.DayPreferences {
		if pref.DayIndex < 0 || pref.DayIndex > 6 {
			return updateUserScheduleWithDayPreferencesResponse{
				Success: false,
			}, nil
		}
		if pref.NumberOfArticles < 0 || pref.NumberOfArticles > newsletter.DefaultNumberOfArticlesPerEmail {
			return updateUserScheduleWithDayPreferencesResponse{
				Success: false,
			}, nil
		}
		// Validate Content Topics
		for _, t := range pref.ContentTopics {
			if _, err := contenttopics.GetContentTopicForString(t.Str()); err != nil {
				return updateUserScheduleWithDayPreferencesResponse{
					Success: false,
				}, nil
			}
		}
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, pref := range req.DayPreferences {
			var topics []string
			for _, t := range pref.ContentTopics {
				topics = append(topics, t.Str())
			}
			if err := usernewsletterschedule.UpsertNewsletterDayMetadataForUser(tx, usernewsletterschedule.UpsertNewsletterDayMetadataForUserInput{
				UserID:           userID,
				DayOfWeekIndex:   pref.DayIndex,
				LanguageCode:     *languageCode,
				ContentTopics:    topics,
				NumberOfArticles: pref.NumberOfArticles,
				IsActive:         pref.IsActive,
			}); err != nil {
				return err
			}
		}
		return usernewsletterschedule.UpsertUserNewsletterSchedule(tx, usernewsletterschedule.UpsertUserNewsletterScheduleInput{
			UserID:           userID,
			LanguageCode:     *languageCode,
			IANATimezone:     loc,
			QuarterHourIndex: req.QuarterHourIndex,
			HourIndex:        req.HourIndex,
		})
	}); err != nil {
		return nil, err
	}
	return updateUserScheduleWithDayPreferencesResponse{
		Success: true,
	}, nil
}
