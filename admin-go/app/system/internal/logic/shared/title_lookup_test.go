package shared

import (
	"reflect"
	"testing"
)

func TestCompactPositiveIDs(t *testing.T) {
	input := []int64{0, 9, 9, -2, 4, 0, 3, 4}
	want := []int64{9, 4, 3}
	if got := compactPositiveIDs(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactPositiveIDs mismatch: got=%v want=%v", got, want)
	}
}

func TestCompactPositiveIDsNilWhenEmpty(t *testing.T) {
	if got := compactPositiveIDs(nil); got != nil {
		t.Fatalf("compactPositiveIDs(nil) should return nil, got=%v", got)
	}
	if got := compactPositiveIDs([]int64{0, -1}); got != nil {
		t.Fatalf("compactPositiveIDs should drop invalid values, got=%v", got)
	}
}
