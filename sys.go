package sys

// panics on error or returns value, similar to google's
// must in uuid but this is generic
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
