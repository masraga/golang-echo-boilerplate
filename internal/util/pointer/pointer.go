package pointer

func String(s string) *string {
	return &s
}

func Int64(d int64) *int64 {
	return &d
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
