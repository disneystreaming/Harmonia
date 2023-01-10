package set

import (
	"fmt"
	"testing"
	"time"

	"math/rand"

	"github.com/stretchr/testify/assert"
)

var intSet Set[int]
var stringSet Set[string]

func setup() {
	intSet = NewSetOf(1, 2, 4, 8)
	stringSet = NewSetOf("1", "2", "3", "4")
}

func TestSetAdd(t *testing.T) {
	// arrange
	setup()
	expectedInts := []int{1, 2, 4, 8, 16}
	expectedStrings := []string{"1", "2", "3", "4", "5"}

	// act
	intSet.Add(16)
	stringSet.Add("5")

	// assert
	if !assert.ElementsMatch(t, expectedInts, intSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedInts, intSet.Values())
	}

	if !assert.ElementsMatch(t, expectedStrings, stringSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedStrings, stringSet.Values())
	}
}

func TestSetAddPresent(t *testing.T) {
	// arrange
	setup()
	expectedInts := []int{1, 2, 4, 8}
	var err error

	// act
	err = intSet.Add(8)

	// assert
	if !assert.ElementsMatch(t, expectedInts, intSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedInts, intSet.Values())
	}

	if err != nil {
		t.Errorf("unexpected error occurred when adding to set, expected nil")
	}
}

func TestSetDelete(t *testing.T) {
	// arrange
	setup()
	expectedInts := []int{1, 2, 4}
	expectedStrings := []string{"1", "2", "3"}

	// act
	intSet.Delete(8)
	stringSet.Delete("4")

	// assert
	if !assert.ElementsMatch(t, expectedInts, intSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedInts, intSet.Values())
	}

	if !assert.ElementsMatch(t, expectedStrings, stringSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedStrings, stringSet.Values())
	}
}

func TestSetDeleteNotPresent(t *testing.T) {
	// arrange
	setup()
	expectedInts := []int{1, 2, 4, 8}
	expectedStrings := []string{"1", "2", "3", "4"}

	// act
	intSet.Delete(32)
	stringSet.Delete("6")

	// assert
	if !assert.ElementsMatch(t, expectedInts, intSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedInts, intSet.Values())
	}

	if !assert.ElementsMatch(t, expectedStrings, stringSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedStrings, stringSet.Values())
	}
}

func TestSetContains(t *testing.T) {
	// arrange
	setup()

	// assert
	if !intSet.Contains(1) {
		t.Error("unexpected output. wanted true, got false", intSet, 1)
	}
	if intSet.Contains(-1) {
		t.Error("unexpected output. wanted false, got true", intSet, -1)
	}
}

func TestSetSize(t *testing.T) {
	// arrange
	setup()

	// assert
	if intSet.Size() != 4 {
		t.Errorf("unexpected value. wanted %v, got %v", 4, intSet.Size())
	}
}

func TestSetValues(t *testing.T) {
	// arrange
	setup()
	expectedInts := []int{1, 2, 4, 8}
	expectedStrings := []string{"1", "2", "3", "4"}

	// assert
	if !assert.ElementsMatch(t, expectedInts, intSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedInts, intSet.Values())
	}

	if !assert.ElementsMatch(t, expectedStrings, stringSet.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedStrings, stringSet.Values())
	}
}

func TestSetIntersect(t *testing.T) {
	// arrange
	setup()
	disjoint := NewSetOf(3, 9, 27, 81)
	intersecting := NewSetOf(1, 4, 16, 64)
	expectedDisjoint := []int{}
	expectedIntersection := []int{1, 4}

	// act
	actualDisjoint := intSet.Intersect(disjoint)
	actualIntersection := intSet.Intersect(intersecting)

	// assert
	if !assert.ElementsMatch(t, expectedDisjoint, actualDisjoint.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedDisjoint, actualDisjoint.Values())
	}

	if !assert.ElementsMatch(t, expectedIntersection, actualIntersection.Values()) {
		t.Errorf("unexpected values. wanted %v, got %v", expectedIntersection, actualIntersection.Values())
	}
}

func TestSetEquals(t *testing.T) {
	// arrange
	setup()
	var nilSet Set[int] = nil
	var emptySet Set[int] = NewSet[int]()
	var copy Set[int] = NewSetOf(1, 2, 4, 8)
	var superset Set[int] = NewSetOf(1, 2, 4, 8, 16)
	var different Set[int] = NewSetOf(1, 3, 9, 27)

	// assert
	if intSet.Equals(nilSet) {
		t.Errorf("unexpected output. %v should not equal %v", intSet, nilSet)
	}

	if intSet.Equals(emptySet) {
		t.Errorf("unexpected output. %v should not equal %v", intSet, emptySet)
	}

	if !intSet.Equals(copy) {
		t.Errorf("unexpected output. %v should equal %v", intSet, copy)
	}

	if intSet.Equals(superset) {
		t.Errorf("unexpected output. %v should not equal %v", intSet, superset)
	}

	if !intSet.Equals(intSet) {
		t.Errorf("unexpected output. %v should equal %v", intSet, intSet)
	}

	if intSet.Equals(different) {
		t.Errorf("unexpected output. %v should not equal %v", intSet, different)
	}
}

// Basic comparison test
// For 10000 trials with a space of arrays up to length 50000:
//	Set took on average 0.2901 microseconds, Array took on average 11.6131 microseconds
func TestSpeedVsArray(t *testing.T) {
	trials := 10000
	space := 50000
	rand.Seed(time.Now().UnixNano())

	var start int64
	var end int64
	var contains bool

	times := make([]int64, 2*trials)

	for i := 0; i < trials; i++ {
		n := rand.Intn(space-1) + 1            // represents the max length of the set/array
		numRange := rand.Intn((2*space)-1) + 1 // represents the max number generated
		toFind := rand.Intn(space)             // represents the number to find (may or may not exist)

		// generate array and set
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(numRange)
		}
		s := NewSetOf(arr...)

		// time
		start = time.Now().UnixNano()
		contains = arrayContains(arr, toFind)
		end = time.Now().UnixNano()

		times[2*i] = end - start

		start = time.Now().UnixNano()
		if contains != s.Contains(toFind) {
			t.Errorf("mismatch! arrayContains and set.Contains disagree: %v found in array: %v, in set: %v\n", toFind, contains, !contains)
		}
		end = time.Now().UnixNano()

		times[2*i+1] = end - start
	}

	var avgSetTime float64
	var avgArrayTime float64

	for i := 0; i < trials; i++ {
		avgArrayTime += float64(times[2*i])
		avgSetTime += float64(times[2*i+1])
	}

	avgArrayTime /= float64(trials) // average
	avgArrayTime /= 1e3             // convert nano to microseconds
	avgSetTime /= float64(trials)   // average
	avgSetTime /= 1e3               // convert nano to microseconds

	fmt.Printf("Set took on average %v microseconds, Array took on average %v microseconds", avgSetTime, avgArrayTime)
}

func arrayContains[K comparable](arr []K, toFind K) bool {
	for _, item := range arr {
		if item == toFind {
			return true
		}
	}

	return false
}
