package usernewsletterpreferences

import (
	"babblegraph/model/content"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/timeutils"
	"babblegraph/wordsmith"
	"time"
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

func (s *ScheduleWithMetadata) IsSendRequested(utcWeekday time.Weekday) bool {
	return len(s.userScheduleDays) == 0 || s.userScheduleDays[int(utcWeekday)].IsActive
}

func (s *ScheduleWithMetadata) GetNumberOfDocuments() int {
	return s.NumberOfArticlesPerEmail
}

func (s *ScheduleWithMetadata) ConvertUTCTimeToUserDate(c ctx.LogContext, utcTime time.Time) (*time.Time, error) {
	return resolveUTCMidnightWithNewsletterSchedule(c, timeutils.ConvertToMidnight(utcTime), dbUserNewsletterSchedule{
		IANATimezone:     s.IANATimezone,
		HourIndex:        s.HourIndex,
		QuarterHourIndex: s.QuarterHourIndex,
	})
}
