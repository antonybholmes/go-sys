package sys

import "slices"

// Implementation of a set
type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Adds length function
func (s *Set[T]) Len() int {
	return len(s.items)
}

// Adds to the set
func (s *Set[T]) Add(v T) {
	s.items[v] = struct{}{}
}

// Removes an item from the set
func (s *Set[T]) Remove(v T) {
	delete(s.items, v)
}

// Returns a boolean value describing if the value exists in the set
func (s *Set[T]) Has(v T) bool {
	_, ok := s.items[v]
	return ok
}

func (s *Set[T]) Keys() []T {
	keys := make([]T, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	return keys
}

// Returns the insection of a map with another
func (s *Set[T]) Intersect(s2 *Set[T]) *Set[T] {
	ret := NewSet[T]()

	for k := range s.items {
		if s2.Has(k) {
			ret.Add(k)
		}
	}

	return ret
}

func (s *Set[T]) Union(s2 *Set[T]) *Set[T] {
	ret := NewSet[T]()

	for k := range s.items {
		ret.Add(k)
	}

	for k := range s2.items {
		ret.Add(k)
	}

	return ret
}

func (s *Set[T]) Update(values *Set[T]) *Set[T] {
	for v := range values.items {
		s.Add(v)
	}

	return s
}

func (s *Set[T]) UpdateFromList(values []T) *Set[T] {
	for _, v := range values {
		s.Add(v)
	}

	return s
}

type StringSet struct {
	*Set[string]
}

func NewStringSet() *StringSet {
	return &StringSet{
		Set: NewSet[string](), // assuming NewSet initializes the map
	}
}

// Specialized keys that returns sorted keys
func (s *StringSet) Keys() []string {
	keys := make([]string, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}

func (s *StringSet) Intersect(s2 *StringSet) *StringSet {
	ret := NewStringSet()

	for k := range s.items {
		if s2.Has(k) {
			ret.Add(k)
		}
	}

	return ret
}

func (s *StringSet) Update(values *StringSet) *StringSet {
	for v := range values.items {
		s.Add(v)
	}

	return s
}

func (s *StringSet) UpdateFromList(values []string) *StringSet {
	// so we don't infinite loop
	s.Set.UpdateFromList(values)

	return s
}

// func StringSetSort(s *Set[string]) []string {

// 	sortedGenes := make([]string, 0, len(s.items))

// 	for key := range s.items {
// 		sortedGenes = append(sortedGenes, key)
// 	}

// 	slices.Sort(sortedGenes)

// 	return sortedGenes
// }
