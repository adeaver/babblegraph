package usernewsletterschedule

import "testing"

func TestGetClosetSendTimeInUTC(t *testing.T) {
	type testCase struct {
		InputDayIndex       int
		IANATimezone        string
		ExpectedDayIndexUTC int
	}
	tcs := []testCase{
		{
			InputDayIndex:       0,
			IANATimezone:        "Pacific/Kiritimati", // UTC+14, no daylight savings time
			ExpectedDayIndexUTC: 6,
		}, {
			InputDayIndex:       0,
			IANATimezone:        "Etc/GMT+5",
			ExpectedDayIndexUTC: 0,
		}, {
			InputDayIndex:       2,
			IANATimezone:        "Pacific/Kiritimati", // UTC+14, no daylight savings time
			ExpectedDayIndexUTC: 1,
		}, {
			InputDayIndex:       2,
			IANATimezone:        "Pacific/Pago_Pago", // UTC-11
			ExpectedDayIndexUTC: 2,
		},
	}
	for idx, tc := range tcs {
		result, _, _, err := GetClosetSendTimeInUTC(tc.InputDayIndex, tc.IANATimezone)
		if err != nil {
			t.Fatalf("Got error on test case %d: %s", idx+1, err.Error())
		}
		if *result != tc.ExpectedDayIndexUTC {
			t.Errorf("Error on test case %d. Got %d, but expected %d", idx+1, *result, tc.ExpectedDayIndexUTC)
		}
	}
}

func TestConvertIndexedTimeUTCToUserTimezone(t *testing.T) {
	type simpleDatetime struct {
		DayIndex         int
		HourIndex        int
		QuarterHourIndex int
	}
	type testCase struct {
		InputUTC          simpleDatetime
		InputIANATimezone string
		ExpectedDatetime  simpleDatetime
	}
	tcs := []testCase{
		{
			InputUTC: simpleDatetime{
				DayIndex:         5,
				HourIndex:        18,
				QuarterHourIndex: 0,
			},
			InputIANATimezone: "Pacific/Kiritimati",
			ExpectedDatetime: simpleDatetime{
				DayIndex:         6,
				HourIndex:        8,
				QuarterHourIndex: 0,
			},
		}, {
			InputUTC: simpleDatetime{
				DayIndex:         5,
				HourIndex:        18,
				QuarterHourIndex: 0,
			},
			InputIANATimezone: "Pacific/Pago_Pago",
			ExpectedDatetime: simpleDatetime{
				DayIndex:         5,
				HourIndex:        7,
				QuarterHourIndex: 0,
			},
		}, {
			InputUTC: simpleDatetime{
				DayIndex:         5,
				HourIndex:        20,
				QuarterHourIndex: 0,
			},
			InputIANATimezone: "Asia/Kathmandu",
			ExpectedDatetime: simpleDatetime{
				DayIndex:         6,
				HourIndex:        1,
				QuarterHourIndex: 3,
			},
		},
	}
	for idx, tc := range tcs {
		resultDayIndex, resultHourIndex, resultQuarterHourIndex, err := ConvertIndexedTimeUTCToUserTimezone(tc.InputUTC.DayIndex, tc.InputUTC.HourIndex, tc.InputUTC.QuarterHourIndex, tc.InputIANATimezone)
		if err != nil {
			t.Fatalf("Got error on test case %d: %s", idx+1, err.Error())
		}
		switch {
		case *resultDayIndex != tc.ExpectedDatetime.DayIndex:
			t.Errorf("Error on test case %d: Got day %d, but expected %d", idx+1, *resultDayIndex, tc.ExpectedDatetime.DayIndex)
		case *resultHourIndex != tc.ExpectedDatetime.HourIndex:
			t.Errorf("Error on test case %d: Got hour %d, but expected %d", idx+1, *resultHourIndex, tc.ExpectedDatetime.HourIndex)
		case *resultQuarterHourIndex != tc.ExpectedDatetime.QuarterHourIndex:
			t.Errorf("Error on test case %d: Got quarter hour %d, but expected %d", idx+1, *resultQuarterHourIndex, tc.ExpectedDatetime.QuarterHourIndex)
		default:
			// no-op
		}
	}
}
