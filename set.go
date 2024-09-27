package sys

import "slices"

// Implementation of a set
type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() *Set[T] {
	s := make(Set[T])
	return &s
}

// Adds length function
func (s *Set[T]) Len() int {
	return len(*s)
}

// Adds to the set
func (s *Set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

// Removes an item from the set
func (s *Set[T]) Remove(v T) {
	delete(*s, v)
}

// Returns a boolean value describing if the value exists in the set
func (s *Set[T]) Has(v T) bool {
	_, ok := (*s)[v]
	return ok
}

// Returns the insection of a map with another
func (s *Set[T]) Intersect(s2 *Set[T]) *Set[T] {
	ret := NewSet[T]()

	for k := range *s {
		if s2.Has(k) {
			ret.Add(k)
		}
	}

	return ret
}

func (s *Set[T]) Union(s2 *Set[T]) *Set[T] {
	ret := NewSet[T]()

	for k := range *s {
		ret.Add(k)
	}

	for k := range *s2 {
		ret.Add(k)
	}

	return ret
}

func (s *Set[T]) Update(values *Set[T]) *Set[T] {
	for v := range *values {
		s.Add(v)
	}

	return s
}

func (s *Set[T]) UpdateList(values []T) *Set[T] {
	for _, v := range values {
		s.Add(v)
	}

	return s
}

func StringSetSort(s *Set[string]) []string {

	sortedGenes := make([]string, 0, len(*s))

	for key := range *s {
		sortedGenes = append(sortedGenes, key)
	}

	slices.Sort(sortedGenes)

	return sortedGenes
}
