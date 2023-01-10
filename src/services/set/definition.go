package set

type Set[K comparable] interface {
	// Add adds the given values to the set
	Add(vals ...K) error
	// Delete removes the values from the set
	Delete(vals ...K) error
	// Contains returns true if the given value is contained within the set
	Contains(val K) bool
	// Size returns the size of the set
	Size() int
	// Values returns an iterable slice containing the same values of the set
	Values() []K
	// Intersect returns the intersection of the set with the given other set
	// the underlying set will be empty if there is no intersection
	Intersect(Set[K]) Set[K]
	// Equals returns true if the set is equal to the given other set
	Equals(Set[K]) bool
}

var exists struct{}
