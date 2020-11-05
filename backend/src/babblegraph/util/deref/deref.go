package deref

func String(s *string, d string) string {
	if s == nil {
		return d
	}
	return *s
}
