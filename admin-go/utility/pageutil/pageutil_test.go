package pageutil

import "testing"

func TestNormalize(t *testing.T) {
	pageNum, pageSize := Normalize(0, 0)
	if pageNum != DefaultPageNum || pageSize != DefaultPageSize {
		t.Fatalf("default normalize mismatch: %d %d", pageNum, pageSize)
	}

	pageNum, pageSize = Normalize(-1, 999)
	if pageNum != DefaultPageNum || pageSize != MaxPageSize {
		t.Fatalf("bounded normalize mismatch: %d %d", pageNum, pageSize)
	}

	pageNum, pageSize = Normalize(3, 50)
	if pageNum != 3 || pageSize != 50 {
		t.Fatalf("preserve normalize mismatch: %d %d", pageNum, pageSize)
	}
}
