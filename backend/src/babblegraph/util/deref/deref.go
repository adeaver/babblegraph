package deref

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
