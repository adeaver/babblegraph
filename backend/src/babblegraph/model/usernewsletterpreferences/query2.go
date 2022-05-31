package usernewsletterpreferences

import (
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/timeutils"
	"babblegraph/wordsmith"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

func GetUserNewsletterSchedule(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, utcMidnight *time.Time) (*ScheduleWithMetadata, error) {
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
		c.Debugf("Input %+v, output %+v", todayUTCMidnight, userSendTime)
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
			c.Debugf("Offset %d", offset)
			sort.SliceStable(userScheduleDays, func(i, j int) bool {
				return (userScheduleDays[i].DayOfWeekIndex+offset)%7 < (userScheduleDays[j].DayOfWeekIndex+offset)%7
			})
			c.Debugf("%+v", userScheduleDays)
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

type UpsertUserNewsletterScheduleInput struct {
	UserID                   users.UserID
	LanguageCode             wordsmith.LanguageCode
	IANATimezone             *time.Location
	HourIndex                int
	QuarterHourIndex         int
	NumberOfArticlesPerEmail int
	IsActiveForDays          []bool
}

func UpsertUserNewsletterSchedule(c ctx.LogContext, tx *sqlx.Tx, input UpsertUserNewsletterScheduleInput) error {
	c.Debugf("Inserting days")
	for idx, isActive := range input.IsActiveForDays {
		if err := upsertNewsletterDayMetadataForUser(tx, upsertNewsletterDayMetadataForUserInput{
			UserID:         input.UserID,
			LanguageCode:   input.LanguageCode,
			DayOfWeekIndex: idx,
			IsActive:       isActive,
		}); err != nil {
			return err
		}
	}
	c.Debugf("Inserting schedule")
	return upsertUserNewsletterSchedule(tx, upsertUserNewsletterScheduleInput{
		UserID:                   input.UserID,
		LanguageCode:             input.LanguageCode,
		IANATimezone:             input.IANATimezone,
		HourIndex:                input.HourIndex,
		QuarterHourIndex:         input.QuarterHourIndex,
		NumberOfArticlesPerEmail: input.NumberOfArticlesPerEmail,
	})
}
