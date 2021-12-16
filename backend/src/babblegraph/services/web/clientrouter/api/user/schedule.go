package user

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/newsletter"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterschedule"
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
	if err := json.Unmarshal(body, req); err != nil {
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
