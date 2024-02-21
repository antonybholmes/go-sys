package sys

// panics on error or returns value, similar to google's
// must in uuid but this is generic
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
