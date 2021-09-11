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
