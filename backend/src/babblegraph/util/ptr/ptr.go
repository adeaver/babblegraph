package ptr

import "time"

func String(s string) *string {
	return &s
}

func Int64(i int64) *int64 {
	return &i
}

func Int(i int) *int {
	return &i
}

func Time(t time.Time) *time.Time {
	return &t
}

func Duration(d time.Duration) *time.Duration {
	return &d
}

func Bool(b bool) *bool {
	return &b
}
