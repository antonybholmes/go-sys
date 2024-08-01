package sys

import (
	"cmp"
	"slices"
	"testing"
)

func TestArgsort(t *testing.T) {
	orig := []string{"z", "a", "c", "b"}
	copy := append([]string{}, orig...)

	indices := Argsort(orig)

	sorted := copy
	slices.SortStableFunc(sorted, func(a, b string) int {
		return cmp.Compare(a, b)
	})

	for i := range orig {
		if got, want := orig[indices[i]], sorted[i]; got != want {
			t.Errorf("got: %q, want: %q", got, want)
		}
	}
}
