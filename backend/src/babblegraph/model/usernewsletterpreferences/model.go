package usernewsletterpreferences

import (
	"babblegraph/model/content"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/timeutils"
	"babblegraph/wordsmith"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	minimumNumberOfArticles = 4
	maximumNumberOfArticles = 12
	defaultNumberOfArticles = 12

	defaultUTCSendTimeHour = 11
)

type UserNewsletterPreferences struct {
	UserID                                   users.UserID
	LanguageCode                             wordsmith.LanguageCode
	ShouldIncludeLemmaReinforcementSpotlight bool
	PodcastPreferences                       PodcastPreferences
	Schedule                                 Schedule
}

type PodcastPreferences struct {
	ArePodcastsEnabled         bool
	IncludeExplicitPodcasts    bool
	MinimumDurationNanoseconds *time.Duration
	MaximumDurationNanoseconds *time.Duration
	ExcludedSourceIDs          []content.SourceID
}

type userLemmaReinforcementSpotlightPreferencesID string

type dbUserLemmaReinforcementSpotlightPreferences struct {
	ID                                       userLemmaReinforcementSpotlightPreferencesID `db:"_id"`
	LanguageCode                             wordsmith.LanguageCode                       `db:"language_code"`
	UserID                                   users.UserID                                 `db:"user_id"`
	ShouldIncludeLemmaReinforcementSpotlight bool                                         `db:"should_include_lemma_reinforcement_spotlight"`
}

type userPodcastPreferecesID string

type dbUserPodcastPreferences struct {
	CreatedAt                  time.Time               `db:"created_at"`
	LastModifiedAt             time.Time               `db:"last_modified_at"`
	ID                         userPodcastPreferecesID `db:"_id"`
	LanguageCode               wordsmith.LanguageCode  `db:"language_code"`
	UserID                     users.UserID            `db:"user_id"`
	ArePodcastsEnabled         bool                    `db:"podcasts_enabled"`
	IncludeExplicitPodcasts    bool                    `db:"include_explicit_podcasts"`
	MinimumDurationNanoseconds *time.Duration          `db:"minimum_duration_nanoseconds"`
	MaximumDurationNanoseconds *time.Duration          `db:"maximum_duration_nanoseconds"`
}

type dbUserPodcastSourceMappingID string

type dbUserPodcastSourcePreferences struct {
	CreatedAt      time.Time                    `db:"created_at"`
	LastModifiedAt time.Time                    `db:"last_modified_at"`
	ID             dbUserPodcastSourceMappingID `db:"_id"`
	LanguageCode   wordsmith.LanguageCode       `db:"language_code"`
	UserID         users.UserID                 `db:"user_id"`
	SourceID       content.SourceID             `db:"source_id"`
	IsActive       bool                         `db:"is_active"`
}

type scheduleID string

type dbUserNewsletterSchedule struct {
	ID                       scheduleID             `db:"_id"`
	UserID                   users.UserID           `db:"user_id"`
	LanguageCode             wordsmith.LanguageCode `db:"language_code"`
	IANATimezone             string                 `db:"iana_timezone"`
	HourIndex                int                    `db:"hour_of_day_index"`
	QuarterHourIndex         int                    `db:"quarter_hour_index"`
	NumberOfArticlesPerEmail int                    `db:"number_of_articles_per_email"`
}

type dayID string

type dbUserNewsletterDayMetadata struct {
	ID             dayID                  `db:"_id"`
	UserID         users.UserID           `db:"user_id"`
	LanguageCode   wordsmith.LanguageCode `db:"language_code"`
	DayOfWeekIndex int                    `db:"day_of_week_index"`
	IsActive       bool                   `db:"is_active"`
}

type Schedule interface {
	IsSendRequested(utcWeekday time.Weekday) bool
	GetUTCHourAndQuarterHourIndex() (_hourIndex, _quarterHourIndex int)
	GetNumberOfDocuments() int
	ConvertUTCTimeToUserDate(c ctx.LogContext, utcTime time.Time) (*time.Time, error)
}

type ScheduleWithMetadata struct {
	userScheduleDays    []dbUserNewsletterDayMetadata
	utcHourIndex        int
	utcQuarterHourIndex int

	NumberOfArticlesPerEmail int
	IANATimezone             string
	HourIndex                int
	QuarterHourIndex         int
	IsActiveForDay           []bool
}

func getUserNewsletterSchedule(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, utcMidnight *time.Time) (*ScheduleWithMetadata, error) {
	isActiveForDay := []bool{true, true, true, true, true, true, true}
	var userScheduleDays []dbUserNewsletterDayMetadata
	var ianaTimezone string
	var utcHourIndex, utcQuarterHourIndex, hourIndex, quarterHourIndex, numberOfArticlesPerEmail int
	userSchedule, err := lookupUserNewsletterScheduleForUser(tx, userID, languageCode)
	switch {
	case err != nil:
		return nil, err
	case userSchedule == nil:
		ianaTimezone = "UTC"
		utcHourIndex, utcQuarterHourIndex, hourIndex, quarterHourIndex = defaultUTCSendTimeHour, 0, defaultUTCSendTimeHour, 0
		numberOfArticlesPerEmail = defaultNumberOfArticles
	default:
		todayUTCMidnight := timeutils.ConvertToMidnight(deref.Time(utcMidnight, time.Now().UTC()))
		userSendTime, err := resolveUTCMidnightWithNewsletterSchedule(c, todayUTCMidnight, *userSchedule)
		if err != nil {
			return nil, err
		}
		numberOfArticlesPerEmail = userSchedule.NumberOfArticlesPerEmail
		ianaTimezone = userSchedule.IANATimezone
		hourIndex = userSendTime.Hour()
		quarterHourIndex = userSendTime.Minute() / 15
		userSendTimeUTC := userSendTime.UTC()
		utcHourIndex = userSendTimeUTC.Hour()
		utcQuarterHourIndex = userSendTimeUTC.Minute() / 15
		userScheduleDays, err = lookupNewsletterDayMetadataForUser(tx, userID, languageCode)
		switch {
		case err != nil:
			return nil, err
		case len(userScheduleDays) == 0:
			// no-op
		default:
			for _, d := range userScheduleDays {
				isActiveForDay[d.DayOfWeekIndex] = d.IsActive
			}
			offset := int(todayUTCMidnight.Weekday() - userSendTime.Weekday())
			sort.SliceStable(userScheduleDays, func(i, j int) bool {
				return userScheduleDays[i].DayOfWeekIndex+offset < userScheduleDays[j].DayOfWeekIndex+offset
			})
		}
	}
	return &ScheduleWithMetadata{
		NumberOfArticlesPerEmail: numberOfArticlesPerEmail,
		userScheduleDays:         userScheduleDays,
		utcHourIndex:             utcHourIndex,
		utcQuarterHourIndex:      utcQuarterHourIndex,
		IANATimezone:             ianaTimezone,
		HourIndex:                hourIndex,
		QuarterHourIndex:         quarterHourIndex,
		IsActiveForDay:           isActiveForDay,
	}, nil
}

func (s *ScheduleWithMetadata) IsSendRequested(utcWeekday time.Weekday) bool {
	return len(s.userScheduleDays) == 0 || s.userScheduleDays[int(utcWeekday)].IsActive
}

func (s *ScheduleWithMetadata) GetUTCHourAndQuarterHourIndex() (_hourIndex, _quarterHourIndex int) {
	return s.utcHourIndex, s.utcQuarterHourIndex
}

func (s *ScheduleWithMetadata) GetNumberOfDocuments() int {
	return s.NumberOfArticlesPerEmail
}

// NOTE: This function does not return an accurate time - just an accurate date
func (s *ScheduleWithMetadata) ConvertUTCTimeToUserDate(c ctx.LogContext, utcTime time.Time) (*time.Time, error) {
	return resolveUTCMidnightWithNewsletterSchedule(c, timeutils.ConvertToMidnight(utcTime), dbUserNewsletterSchedule{
		IANATimezone:     s.IANATimezone,
		HourIndex:        s.HourIndex,
		QuarterHourIndex: s.QuarterHourIndex,
	})
}
