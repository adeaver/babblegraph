package usernewsletterpreferences

import (
	"babblegraph/util/ctx"
	"testing"
	"time"
)

func TestResolveUTCMidnightWithNewsletterSchedule(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	testCases := []struct {
		utcYear              int
		utcMonth             time.Month
		utcDayIndex          int
		userHourIndex        int
		userQuarterHourIndex int
		userTimezone         string
		expectedDayIndex     int
	}{
		// Christmas 2020 is a Friday
		{
			utcYear:              2020,
			utcMonth:             time.December,
			utcDayIndex:          25,
			userHourIndex:        21,
			userQuarterHourIndex: 0,
			userTimezone:         "America/Los_Angeles",
			expectedDayIndex:     4,
		}, {
			utcYear:              2020,
			utcMonth:             time.December,
			utcDayIndex:          25,
			userHourIndex:        5,
			userQuarterHourIndex: 0,
			userTimezone:         "America/Los_Angeles",
			expectedDayIndex:     5,
		}, {
			utcYear:              2020,
			utcMonth:             time.December,
			utcDayIndex:          25,
			userHourIndex:        12,
			userQuarterHourIndex: 0,
			userTimezone:         "Pacific/Tarawa",
			expectedDayIndex:     5,
		}, {
			utcYear:              2020,
			utcMonth:             time.December,
			utcDayIndex:          25,
			userHourIndex:        11,
			userQuarterHourIndex: 0,
			userTimezone:         "Pacific/Tarawa",
			expectedDayIndex:     6,
		},
	}
	for idx, tc := range testCases {
		utcMidnight := time.Date(tc.utcYear, tc.utcMonth, tc.utcDayIndex, 0, 0, 0, 0, time.UTC)
		result, err := resolveUTCMidnightWithNewsletterSchedule(c, utcMidnight, dbUserNewsletterSchedule{
			IANATimezone:     tc.userTimezone,
			HourIndex:        tc.userHourIndex,
			QuarterHourIndex: tc.userQuarterHourIndex,
		})
		switch {
		case err != nil:
			t.Errorf("Error on test %d: %s", idx, err.Error())
		case result != nil:
			if tc.expectedDayIndex != int(result.Weekday()) {
				t.Errorf("Error on test %d: expected %d, but got %d", idx, tc.expectedDayIndex, int(result.Weekday()))
			}
			if result.UTC().Weekday() != utcMidnight.Weekday() {
				t.Errorf("Error on test %d: expected converted utc day to be %d, but got %d", idx, utcMidnight.Weekday(), result.UTC().Weekday())
			}
		default:
			t.Errorf("Error on test %d: no result, but no err", idx)
		}

	}
}
