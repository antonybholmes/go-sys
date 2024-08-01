package sys

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	orig := []string{"z", "z", "a", "c", "b"}

	set := NewSet[string]()
	set.UpdateList(orig)

	fmt.Printf("set %v", set)
}
