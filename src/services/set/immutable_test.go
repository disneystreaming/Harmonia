package set

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var intImmutableSet Set[int]
var stringImmutableSet Set[string]

func setupImmutable() {
	intImmutableSet = NewImmutableOf(1, 2, 4, 8)
	stringImmutableSet = NewImmutableOf("1", "2", "3", "4")
}

func TestImmutableAdd(t *testing.T) {
	// arrange
	setupImmutable()
	expected := fmt.Errorf("unsupported operation: Add. cannot modify an immutable set")
	var err error

	// act
	err = intImmutableSet.Add(16)

	// assert
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("unexpected return value. expected %v, got %v", expected, err)
	}
}

func TestImmutableDelete(t *testing.T) {
	// arrange
	setupImmutable()
	expected := fmt.Errorf("unsupported operation: Delete. cannot modify an immutable set")
	var err error

	// act
	err = stringImmutableSet.Delete("4")

	// assert
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("unexpected return value. expected %v, got %v", expected, err)
	}
}

func TestImmutableContains(t *testing.T) {
	// arrange
	setupImmutable()

	// assert
	if !intImmutableSet.Contains(1) {
		t.Error("unexpected output. wanted true, got false", intImmutableSet, 1)
	}
	if intImmutableSet.Contains(-1) {
		t.Error("unexpected output. wanted false, got true", intImmutableSet, -1)
	}
}

func TestImmutableSize(t *testing.T) {
	// arrange
	setupImmutable()

	// assert
	if intImmutableSet.Size() != 4 {
		t.Errorf("unexpected value. wanted %v, got %v", 4, intImmutableSet.Size())
	}
}

func TestImmutableValues(t *testing.T) {
	// arrange
	setupImmutable()
	expectedInts := []int{1, 2, 4, 8}
	expectedStrings := []string{"1", "2", "3", "4"}

	// assert
	if !assert.ElementsMatch(t, expectedInts, intImmutableSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedInts, intImmutableSet.Values())
	}

	if !assert.ElementsMatch(t, expectedStrings, stringImmutableSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedStrings, stringImmutableSet.Values())
	}
}

func TestImmutableIntersect(t *testing.T) {
	// arrange
	setupImmutable()
	disjoint := NewImmutableOf(3, 9, 27, 81)
	intersecting := NewImmutableOf(1, 4, 16, 64)
	expectedDisjoint := []int{}
	expectedIntersection := []int{1, 4}

	// act
	actualDisjoint := intImmutableSet.Intersect(disjoint)
	actualIntersection := intImmutableSet.Intersect(intersecting)

	// assert
	if !assert.ElementsMatch(t, expectedDisjoint, actualDisjoint.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedDisjoint, actualDisjoint.Values())
	}

	if !assert.ElementsMatch(t, expectedIntersection, actualIntersection.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedIntersection, actualIntersection.Values())
	}
}

func TestImmutableEquals(t *testing.T) {
	// arrange
	setupImmutable()
	var nilSet Set[int] = nil
	var emptySet Set[int] = NewSet[int]()
	var copy Set[int] = NewImmutableOf(1, 2, 4, 8)
	var superset Set[int] = NewImmutableOf(1, 2, 4, 8, 16)
	var different Set[int] = NewImmutableOf(1, 3, 9, 27)

	// assert
	if intImmutableSet.Equals(nilSet) {
		t.Errorf("unexpected output. %v should not equal %v", intImmutableSet, nilSet)
	}

	if intImmutableSet.Equals(emptySet) {
		t.Errorf("unexpected output. %v should not equal %v", intImmutableSet, emptySet)
	}

	if !intImmutableSet.Equals(copy) {
		t.Errorf("unexpected output. %v should equal %v", intImmutableSet, copy)
	}

	if intImmutableSet.Equals(superset) {
		t.Errorf("unexpected output. %v should not equal %v", intImmutableSet, superset)
	}

	if !intImmutableSet.Equals(intImmutableSet) {
		t.Errorf("unexpected output. %v should equal %v", intImmutableSet, intImmutableSet)
	}

	if intImmutableSet.Equals(different) {
		t.Errorf("unexpected output. %v should not equal %v", intImmutableSet, different)
	}
}
