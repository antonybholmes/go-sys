package sys

import "slices"

// Implementation of a set using generics and maps
// since Go does not have sets built in
type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Returns the number of items in the set
func (s *Set[T]) Len() int {
	return len(s.items)
}

// Add an item to the set
func (s *Set[T]) Add(v T) *Set[T] {
	s.items[v] = struct{}{}
	return s
}

// Removes an item from the set
func (s *Set[T]) Remove(v T) *Set[T] {
	delete(s.items, v)
	return s
}

// Returns a boolean value describing if the value exists in the set
func (s *Set[T]) Has(v T) bool {
	_, ok := s.items[v]
	return ok
}

// Returns a slice of the keys in the set
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

// Returns the union of a map with another. The
// original sets are not modified.
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

// Add values from a list to the set
func (s *Set[T]) ListUpdate(values []T) *Set[T] {
	for _, v := range values {
		s.Add(v)
	}

	return s
}

// Check if any of the values in the list are in this set
func (s *Set[T]) ListContains(values []T) bool {
	return slices.ContainsFunc(values, s.Has)
}

// Check if any the values in the s2 are in this set
func (s *Set[T]) Contains(s2 *Set[T]) bool {
	return s.ListContains(s2.Keys())
}

type StringSet struct {
	// functions and fields of Set[string]
	// are promoted to StringSet, but we can
	// override them if needed
	*Set[string]
}

// Creates a new Set of strings
func NewStringSet() *StringSet {
	return &StringSet{
		Set: NewSet[string](), // assuming NewSet initializes the map
	}
}

// Shadowed method - StringSet keys are returned sorted
func (s *StringSet) Keys() []string {
	keys := make([]string, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}

func (s *StringSet) Add(v string) *StringSet {
	s.items[v] = struct{}{}
	return s
}

func (s *StringSet) Remove(v string) *StringSet {
	delete(s.items, v)
	return s
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

func (s *StringSet) Union(s2 *StringSet) *StringSet {
	ret := NewStringSet()

	for k := range s.items {
		ret.Add(k)
	}

	for k := range s2.items {
		ret.Add(k)
	}

	return ret
}

func (s *StringSet) Update(values *StringSet) *StringSet {
	for v := range values.items {
		s.Add(v)
	}

	return s
}

// Add values from a list to the set and returns a
// pointer to the updated set. Suitable for chaining.
func (s *StringSet) ListUpdate(values []string) *StringSet {
	// so we don't infinite loop
	s.Set.ListUpdate(values)

	return s
}

func (s *StringSet) Contains(s2 *StringSet) bool {
	return s.Set.ListContains(s2.Keys())
}

// func StringSetSort(s *Set[string]) []string {

// 	sortedGenes := make([]string, 0, len(s.items))

// 	for key := range s.items {
// 		sortedGenes = append(sortedGenes, key)
// 	}

// 	slices.Sort(sortedGenes)

// 	return sortedGenes
// }
