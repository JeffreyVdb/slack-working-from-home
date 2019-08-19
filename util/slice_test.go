package util

import (
	"testing"
)

func TestContainsString(t *testing.T) {
	abcSlice := []string{"a", "b", "c"}
	if !ContainsString(abcSlice, "a") {
		t.Errorf("a should be part of the array [a, b, c]")
	}

	if ContainsString(abcSlice, "k") {
		t.Errorf("k should not be a part of [a, b, c]")
	}
}
