package deref

import "time"

func String(s *string, d string) string {
	if s == nil {
		return d
	}
	return *s
}

func Int(i *int, d int) int {
	if i == nil {
		return d
	}
	return *i
}

func Int64(i *int64, d int64) int64 {
	if i == nil {
		return d
	}
	return *i
}

func Bool(b *bool, d bool) bool {
	if b == nil {
		return d
	}
	return *b
}

func Time(t *time.Time, d time.Time) time.Time {
	if t == nil {
		return d
	}
	return *t
}
