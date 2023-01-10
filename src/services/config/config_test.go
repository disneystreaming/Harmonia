package config

import (
	"os"
	"testing"
)

// TestIsLocal tests the IsLocal functionality
func TestIsLocal(t *testing.T) {
	testCases := []struct {
		setValue string
		expected bool
	}{
		{
			setValue: "true",
			expected: true,
		},
		{
			setValue: "false",
			expected: false,
		},
		{
			setValue: "junk",
			expected: false,
		},
	}

	for _, test := range testCases {
		os.Setenv("IS_LOCAL", test.setValue)
		actual := IsLocal()
		if actual != test.expected {
			t.Errorf("actual: %v is not equal to expected: %v", actual, test.expected)
		}
	}
}
