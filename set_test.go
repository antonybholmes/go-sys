package sys

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	orig := []string{"z", "z", "a", "c", "b"}

	set := NewStringSet()
	set.UpdateFromList(orig)

	fmt.Printf("set %v", set)
}
