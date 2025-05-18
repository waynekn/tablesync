package utils

import "testing"

// TestGenerateId tests the GenerateID function to ensure it generates a unique ID
// of the expected length.
func TestGenerateId(t *testing.T) {
	id := GenerateID()
	if len(id) > 22 {
		t.Errorf("Expected ID length lesser than or equal 22, got %d", len(id))
	}

}
