package shared

import (
	"reflect"
	"testing"
)

func TestCompactDirIDs(t *testing.T) {
	input := []int64{0, 7, 7, -1, 3, 0, 5, 3}
	want := []int64{7, 3, 5}
	if got := compactDirIDs(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactDirIDs mismatch: got=%v want=%v", got, want)
	}
}
