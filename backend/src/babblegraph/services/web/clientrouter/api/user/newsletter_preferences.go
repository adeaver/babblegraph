package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/deref"
	"babblegraph/util/email"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type userSchedule struct {
	IANATimezone     string `json:"iana_timezone"`
	HourIndex        int    `json:"hour_index"`
	QuarterHourIndex int    `json:"quarter_hour_index"`
	IsActiveForDays  []bool `json:"is_active_for_days"`
}

type userNewsletterPreferences struct {
	LanguageCode                        wordsmith.LanguageCode `json:"language_code"`
	IsLemmaReinforcementSpotlightActive bool                   `json:"is_lemma_reinforcement_spotlight_active"`
	ArePodcastsEnabled                  *bool                  `json:"are_podcasts_enabled,omitempty"`
	IncludeExplicitPodcasts             *bool                  `json:"include_explicit_podcasts,omitempty"`
	MinimumPodcastDurationSeconds       *int64                 `json:"minimum_podcast_duration_seconds,omitempty"`
	MaximumPodcastDurationSeconds       *int64                 `json:"maximum_podcast_duration_seconds,omitempty"`
	NumberOfArticlesPerEmail            int                    `json:"number_of_articles_per_email"`
	Schedule                            userSchedule           `json:"schedule"`
}

type getUserNewsletterPreferencesRequest struct {
	LanguageCode                string `json:"language_code"`
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type getUserNewsletterPreferencesResponse struct {
	Preferences *userNewsletterPreferences `json:"preferences,omitempty"`
	Error       *clienterror.Error         `json:"error,omitempty"`
}

func getUserNewsletterPreferences(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getUserNewsletterPreferencesRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return getUserNewsletterPreferencesResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return getUserNewsletterPreferencesResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	var doesUserHaveAccount bool
	var isLegacySubscriber bool
	var prefs *usernewsletterpreferences.UserNewsletterPreferences
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		prefs, err = usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(r, tx, *userID, *languageCode, nil)
		if err != nil {
			return err
		}
		doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
		if err != nil {
			return err
		}
		userSubscription, err := useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
		if err != nil {
			return err
		}
		isLegacySubscriber = userSubscription != nil && *userSubscription == useraccounts.SubscriptionLevelLegacy
		return err
	}); err != nil {
		return nil, err
	}
	schedule, ok := prefs.Schedule.(*usernewsletterpreferences.ScheduleWithMetadata)
	if !ok {
		return nil, fmt.Errorf("Schedule did not correctly convert to schedule with metadata")
	}
	userPreferences := &userNewsletterPreferences{
		LanguageCode:                        *languageCode,
		IsLemmaReinforcementSpotlightActive: prefs.ShouldIncludeLemmaReinforcementSpotlight,
		NumberOfArticlesPerEmail:            schedule.NumberOfArticlesPerEmail,
		Schedule: userSchedule{
			IANATimezone:     schedule.IANATimezone,
			HourIndex:        schedule.HourIndex,
			QuarterHourIndex: schedule.QuarterHourIndex,
			IsActiveForDays:  schedule.IsActiveForDay,
		},
	}
	switch {
	case userAuth != nil:
		if userAuth.UserID != *userID {
			return getUserNewsletterPreferencesResponse{
				Error: clienterror.ErrorIncorrectKey.Ptr(),
			}, nil
		}
		if userAuth.SubscriptionLevel != nil {
			userPreferences.ArePodcastsEnabled = ptr.Bool(prefs.PodcastPreferences.ArePodcastsEnabled)
			userPreferences.IncludeExplicitPodcasts = ptr.Bool(prefs.PodcastPreferences.IncludeExplicitPodcasts)
			if prefs.PodcastPreferences.MinimumDurationNanoseconds != nil {
				userPreferences.MinimumPodcastDurationSeconds = ptr.Int64(int64(*prefs.PodcastPreferences.MinimumDurationNanoseconds / time.Second))
			}
			if prefs.PodcastPreferences.MaximumDurationNanoseconds != nil {
				userPreferences.MaximumPodcastDurationSeconds = ptr.Int64(int64(*prefs.PodcastPreferences.MaximumDurationNanoseconds / time.Second))
			}
			r.Infof("Inserting podcast preferences")
		}
	case doesUserHaveAccount:
		r.RespondWithStatus(http.StatusForbidden)
		return getUserNewsletterPreferencesResponse{
			Error: clienterror.ErrorNoAuth.Ptr(),
		}, nil
	case !isLegacySubscriber:
		userPreferences.ArePodcastsEnabled = ptr.Bool(prefs.PodcastPreferences.ArePodcastsEnabled)
		userPreferences.IncludeExplicitPodcasts = ptr.Bool(prefs.PodcastPreferences.IncludeExplicitPodcasts)
		if prefs.PodcastPreferences.MinimumDurationNanoseconds != nil {
			userPreferences.MinimumPodcastDurationSeconds = ptr.Int64(int64(*prefs.PodcastPreferences.MinimumDurationNanoseconds / time.Second))
		}
		if prefs.PodcastPreferences.MaximumDurationNanoseconds != nil {
			userPreferences.MaximumPodcastDurationSeconds = ptr.Int64(int64(*prefs.PodcastPreferences.MaximumDurationNanoseconds / time.Second))
		}
	default:
		// no-op
	}
	return getUserNewsletterPreferencesResponse{
		Preferences: userPreferences,
	}, nil
}

type updateUserNewsletterPreferencesRequest struct {
	SubscriptionManagementToken string                    `json:"subscription_management_token"`
	EmailAddress                *string                   `json:"email_address,omitempty"`
	Preferences                 userNewsletterPreferences `json:"preferences"`
}

type updateUserNewsletterPreferencesResponse struct {
	Error   *clienterror.Error `json:"error,omitempty"`
	Success bool               `json:"success"`
}

const (
	errorEmptyEmailAddress clienterror.Error = "no-email-address"
	errorInvalidTimezone   clienterror.Error = "invalid-timezone"
	errorNoActiveDay       clienterror.Error = "no-active-day"
)

func updateUserNewsletterPreferences(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req updateUserNewsletterPreferencesRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return updateUserNewsletterPreferencesResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.Preferences.LanguageCode.Str())
	if err != nil {
		return getUserNewsletterPreferencesResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	userTimezone, err := time.LoadLocation(req.Preferences.Schedule.IANATimezone)
	if err != nil {
		return getUserNewsletterPreferencesResponse{
			Error: errorInvalidTimezone.Ptr(),
		}, nil
	}
	var hasAtLeastOneActiveDay bool
	for _, isActive := range req.Preferences.Schedule.IsActiveForDays {
		hasAtLeastOneActiveDay = hasAtLeastOneActiveDay || isActive
	}
	if !hasAtLeastOneActiveDay {
		return getUserNewsletterPreferencesResponse{
			Error: errorNoActiveDay.Ptr(),
		}, nil
	}
	if userAuth != nil {
		if userAuth.UserID != *userID {
			return getUserNewsletterPreferencesResponse{
				Error: clienterror.ErrorIncorrectKey.Ptr(),
			}, nil
		}
		var podcastPreferences *usernewsletterpreferences.PodcastPreferencesInput
		if userAuth.SubscriptionLevel != nil {
			podcastPreferences = &usernewsletterpreferences.PodcastPreferencesInput{
				ArePodcastsEnabled:      deref.Bool(req.Preferences.ArePodcastsEnabled, true),
				IncludeExplicitPodcasts: deref.Bool(req.Preferences.IncludeExplicitPodcasts, true),
			}
			if req.Preferences.MinimumPodcastDurationSeconds != nil {
				podcastPreferences.MinimumDurationNanoseconds = ptr.Duration(time.Duration(*req.Preferences.MinimumPodcastDurationSeconds) * time.Second)
			}
			if req.Preferences.MaximumPodcastDurationSeconds != nil {
				podcastPreferences.MaximumDurationNanoseconds = ptr.Duration(time.Duration(*req.Preferences.MaximumPodcastDurationSeconds) * time.Second)
			}
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			return usernewsletterpreferences.UpdateUserNewsletterPreferences(r, tx, usernewsletterpreferences.UpdateUserNewsletterPreferencesInput{
				UserID:                              *userID,
				LanguageCode:                        *languageCode,
				IsLemmaReinforcementSpotlightActive: req.Preferences.IsLemmaReinforcementSpotlightActive,
				PodcastPreferences:                  podcastPreferences,
				IANATimezone:                        userTimezone,
				HourIndex:                           req.Preferences.Schedule.HourIndex,
				QuarterHourIndex:                    req.Preferences.Schedule.QuarterHourIndex,
				IsActiveForDays:                     req.Preferences.Schedule.IsActiveForDays,
				NumberOfArticlesPerEmail:            req.Preferences.NumberOfArticlesPerEmail,
			})
		}); err != nil {
			return nil, err
		}
	} else {
		if req.EmailAddress == nil {
			return getUserNewsletterPreferencesResponse{
				Error: errorEmptyEmailAddress.Ptr(),
			}, nil
		}
		var cErr *clienterror.Error
		formattedEmailAddress := email.FormatEmailAddress(*req.EmailAddress)
		err := database.WithTx(func(tx *sqlx.Tx) error {
			user, err := users.GetUser(tx, *userID)
			switch {
			case err != nil:
				return err
			case user.EmailAddress != formattedEmailAddress:
				cErr = clienterror.ErrorInvalidEmailAddress.Ptr()
				return nil
			}
			doesUserHaveAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
			switch {
			case err != nil:
				return err
			case doesUserHaveAccount:
				cErr = clienterror.ErrorNoAuth.Ptr()
				return nil
			}
			input := usernewsletterpreferences.UpdateUserNewsletterPreferencesInput{
				UserID:                              *userID,
				LanguageCode:                        *languageCode,
				IsLemmaReinforcementSpotlightActive: req.Preferences.IsLemmaReinforcementSpotlightActive,
				IANATimezone:                        userTimezone,
				HourIndex:                           req.Preferences.Schedule.HourIndex,
				QuarterHourIndex:                    req.Preferences.Schedule.QuarterHourIndex,
				IsActiveForDays:                     req.Preferences.Schedule.IsActiveForDays,
				NumberOfArticlesPerEmail:            req.Preferences.NumberOfArticlesPerEmail,
			}
			userSubscription, err := useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
			switch {
			case err != nil:
				return err
			case userSubscription == nil,
				*userSubscription != useraccounts.SubscriptionLevelLegacy:
				input.PodcastPreferences = &usernewsletterpreferences.PodcastPreferencesInput{
					ArePodcastsEnabled:      deref.Bool(req.Preferences.ArePodcastsEnabled, true),
					IncludeExplicitPodcasts: deref.Bool(req.Preferences.IncludeExplicitPodcasts, true),
				}
				if req.Preferences.MinimumPodcastDurationSeconds != nil {
					input.PodcastPreferences.MinimumDurationNanoseconds = ptr.Duration(time.Duration(*req.Preferences.MinimumPodcastDurationSeconds) * time.Second)
				}
				if req.Preferences.MaximumPodcastDurationSeconds != nil {
					input.PodcastPreferences.MaximumDurationNanoseconds = ptr.Duration(time.Duration(*req.Preferences.MaximumPodcastDurationSeconds) * time.Second)
				}
			}
			return usernewsletterpreferences.UpdateUserNewsletterPreferences(r, tx, input)
		})
		switch {
		case err != nil:
			return nil, err
		case cErr != nil:
			return getUserNewsletterPreferencesResponse{
				Error: cErr,
			}, nil
		}
	}
	return updateUserNewsletterPreferencesResponse{
		Success: true,
	}, nil
}
