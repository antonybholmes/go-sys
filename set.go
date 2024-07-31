package sys

// Implementation of a set
type Set[T comparable] map[T]struct{}

// Adds an animal to the set
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

// Removes an item from the set
func (s Set[T]) Remove(v T) {
	delete(s, v)
}

// Returns a boolean value describing if the value exists in the set
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Returns the insection of a map with another
func (s Set[T]) Intersect(s2 Set[T]) Set[T] {
	ret := Set[T]{}

	for k := range s {
		if s2.Has(k) {
			ret.Add(k)
		}
	}

	return ret
}

func (s Set[T]) Union(s2 Set[T]) Set[T] {
	ret := Set[T]{}

	for k := range s {
		ret.Add(k)
	}

	for k := range s2 {
		ret.Add(k)
	}

	return ret
}

func (s Set[T]) UpdateSet(values Set[T]) Set[T] {
	for v := range values {
		s.Add(v)
	}

	return s
}

func (s Set[T]) UpdateList(values []T) Set[T] {
	for _, v := range values {
		s.Add(v)
	}

	return s
}
