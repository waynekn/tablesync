package utils

import "testing"

// TestGenerateId tests the GenerateID function to ensure it generates a unique ID
// of the expected length.
func TestGenerateId(t *testing.T) {
	id := GenerateID()
	if len(id) > maxIdLength {
		t.Errorf("Expected ID length lesser than or equal %d, got %d", maxIdLength, len(id))
	}

}

// TestReverseString tests the `reverseString` utility function to ensure it correctly reverses strings.
func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"", ""},
		{"a", "a"},
	}

	for _, test := range tests {
		result := reverseString(test.input)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}
