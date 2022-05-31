package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"time"

	"github.com/jmoiron/sqlx"
)

type getUserNewsletterScheduleRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
	LanguageCode                string `json:"language_code"`
}

type userSchedule struct {
	IANATimezone     string `json:"iana_timezone"`
	HourIndex        int    `json:"hour_index"`
	QuarterHourIndex int    `json:"quarter_hour_index"`
	IsActiveForDays  []bool `json:"is_active_for_days"`
}

type getUserNewsletterScheduleResponse struct {
	Schedule                 *userSchedule      `json:"schedule,omitempty"`
	NumberOfArticlesPerEmail int                `json:"number_of_articles_per_email"`
	Error                    *clienterror.Error `json:"error,omitempty"`
}

func getUserNewsletterSchedule(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getUserNewsletterScheduleRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, clientErr, err := routermiddleware.ValidateUserAuthWithToken(userAuth, routermiddleware.ValidateUserAuthWithTokenInput{
		Token:   req.SubscriptionManagementToken,
		KeyType: routes.SubscriptionManagementRouteEncryptionKey,
	})
	switch {
	case err != nil:
		return nil, err
	case clientErr != nil:
		return getUserNewsletterScheduleResponse{
			Error: clientErr,
		}, nil
	default:
		// no-op
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return getUserNewsletterScheduleResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	var schedule *usernewsletterpreferences.ScheduleWithMetadata
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		schedule, err = usernewsletterpreferences.GetUserNewsletterSchedule(r, tx, *userID, *languageCode, nil)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserNewsletterScheduleResponse{
		NumberOfArticlesPerEmail: schedule.NumberOfArticlesPerEmail,
		Schedule: &userSchedule{
			IANATimezone:     schedule.IANATimezone,
			HourIndex:        schedule.HourIndex,
			QuarterHourIndex: schedule.QuarterHourIndex,
			IsActiveForDays:  schedule.IsActiveForDay,
		},
	}, nil
}

type updateUserNewsletterScheduleRequest struct {
	SubscriptionManagementToken string       `json:"subscription_management_token"`
	LanguageCode                string       `json:"language_code"`
	EmailAddress                *string      `json:"email_address"`
	Schedule                    userSchedule `json:"schedule,omitempty"`
	NumberOfArticlesPerEmail    int          `json:"number_of_articles_per_email"`
}

const (
	errorInvalidTimezone clienterror.Error = "invalid-timezone"
	errorNoActiveDay     clienterror.Error = "no-active-day"
)

type updateUserNewsletterScheduleResponse struct {
	Success bool               `json:"success"`
	Error   *clienterror.Error `json:"error,omitempty"`
}

func updateUserNewsletterSchedule(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req updateUserNewsletterScheduleRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, clientErr, err := routermiddleware.ValidateUserAuthWithToken(userAuth, routermiddleware.ValidateUserAuthWithTokenInput{
		EmailAddress:        req.EmailAddress,
		RequireEmailAddress: true,
		Token:               req.SubscriptionManagementToken,
		KeyType:             routes.SubscriptionManagementRouteEncryptionKey,
	})
	switch {
	case err != nil:
		return nil, err
	case clientErr != nil:
		return updateUserNewsletterScheduleResponse{
			Error: clientErr,
		}, nil
	default:
		// no-op
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return updateUserNewsletterScheduleResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	location, err := time.LoadLocation(req.Schedule.IANATimezone)
	if err != nil {
		return updateUserNewsletterScheduleResponse{
			Error: errorInvalidTimezone.Ptr(),
		}, nil
	}
	var hasActiveDay bool
	for _, d := range req.Schedule.IsActiveForDays {
		hasActiveDay = hasActiveDay || d
	}
	if !hasActiveDay {
		return updateUserNewsletterScheduleResponse{
			Error: errorNoActiveDay.Ptr(),
		}, nil
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return usernewsletterpreferences.UpsertUserNewsletterSchedule(r, tx, usernewsletterpreferences.UpsertUserNewsletterScheduleInput{
			UserID:                   *userID,
			LanguageCode:             *languageCode,
			IANATimezone:             location,
			HourIndex:                req.Schedule.HourIndex,
			QuarterHourIndex:         req.Schedule.QuarterHourIndex,
			NumberOfArticlesPerEmail: req.NumberOfArticlesPerEmail,
			IsActiveForDays:          req.Schedule.IsActiveForDays,
		})
	}); err != nil {
		return nil, err
	}
	return updateUserNewsletterScheduleResponse{
		Success: true,
	}, nil
}
