package batchutil

import (
	"reflect"
	"testing"

	"gbaseadmin/utility/snowflake"
)

func TestCompactIDs(t *testing.T) {
	input := []snowflake.JsonInt64{0, -1, 3, 3, 2, 1}
	want := []snowflake.JsonInt64{3, 2, 1}
	if got := CompactIDs(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("CompactIDs mismatch: got=%v want=%v", got, want)
	}
}

func TestExpandTreeDeleteOrder(t *testing.T) {
	rows := []TreeRow{
		{ID: 1, ParentID: 0},
		{ID: 2, ParentID: 1},
		{ID: 3, ParentID: 2},
		{ID: 4, ParentID: 1},
		{ID: 5, ParentID: 0},
	}
	got := ExpandTreeDeleteOrder([]snowflake.JsonInt64{1, 3, 99}, rows)
	want := []snowflake.JsonInt64{3, 2, 4, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpandTreeDeleteOrder mismatch: got=%v want=%v", got, want)
	}
}

func TestToInt64s(t *testing.T) {
	input := []snowflake.JsonInt64{0, 4, 2}
	want := []int64{4, 2}
	if got := ToInt64s(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("ToInt64s mismatch: got=%v want=%v", got, want)
	}
}
