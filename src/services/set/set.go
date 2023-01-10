package set

import (
	"encoding/json"
	"fmt"
)

type set[K comparable] struct {
	vals map[K]struct{}
}

// NewSet creates a new empty, mutable set
// TODO: return set struct @rsharan
func NewSet[K comparable]() Set[K] {
	s := &set[K]{
		vals: make(map[K]struct{}),
	}

	return s
}

// NewSetOf creates and initializes a mutable set with the given values
func NewSetOf[K comparable](vals ...K) Set[K] {
	s := &set[K]{
		vals: make(map[K]struct{}),
	}

	for _, val := range vals {
		s.vals[val] = exists
	}

	return s
}

// Add adds the given values to the set
func (s *set[K]) Add(vals ...K) error {
	for _, val := range vals {
		s.vals[val] = exists
	}

	return nil
}

// Delete removes the values from the set
func (s *set[K]) Delete(vals ...K) error {
	for _, val := range vals {
		delete(s.vals, val)
	}

	return nil
}

// Contains returns true if the given value is contained within the set
func (s *set[K]) Contains(val K) bool {
	_, c := s.vals[val]
	return c
}

// Size returns the size of the set
func (s *set[K]) Size() int {
	return len(s.vals)
}

// Values returns an iterable slice containing the same values of the set
func (s *set[K]) Values() []K {
	var values []K

	for val := range s.vals {
		values = append(values, val)
	}

	return values
}

// Intersect returns the intersection of the set with the given other set
// the underlying set will be mutable and empty if there is no intersection
func (s *set[K]) Intersect(other Set[K]) Set[K] {
	var intersection []K

	for _, val := range s.Values() {
		if other.Contains(val) {
			intersection = append(intersection, val)
		}
	}

	return NewSetOf(intersection...)
}

// Equals returns true if the set is equal to the given other set
// Equality is defined as:
//	The receiver pointer and given pointer point to the same memory address OR
//	The set pointed to by the receiver pointer and the set pointed to by the given pointer:
//		Are the same size AND
//		Every value in one set is contained in the other, with == being the qualifier for "contained"
func (s *set[K]) Equals(other Set[K]) bool {
	if s == other {
		return true
	}

	if s == nil && other != nil || s != nil && other == nil {
		return false
	}

	if s.Size() != other.Size() {
		return false
	}

	for val := range s.vals {
		if !other.Contains(val) {
			return false
		}
	}

	return true
}

// MarshalJSON implements the Marshaler interface and simply returns the JSON representation of the values in the set
func (s *set[K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// String implements the Stringer interface and returns the string representation of the values in the set
func (s *set[K]) String() string {
	return fmt.Sprint(s.Values())
}
