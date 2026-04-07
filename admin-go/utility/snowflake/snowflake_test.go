package snowflake

import "testing"

func TestSnowflakeGeneratorRemainsMonotonicWhenClockMovesBackwards(t *testing.T) {
	timestamps := []int64{epoch + 1000, epoch + 999, epoch + 1001}
	idx := 0
	gen := &snowflakeGen{
		workerID: 1,
		nowMillis: func() int64 {
			value := timestamps[idx]
			if idx < len(timestamps)-1 {
				idx++
			}
			return value
		},
	}

	first := gen.generate()
	second := gen.generate()
	third := gen.generate()

	if !(first < second && second < third) {
		t.Fatalf("IDs should remain monotonic, got %d, %d, %d", first, second, third)
	}
}

func TestJsonInt64JSONRoundTrip(t *testing.T) {
	value := JsonInt64(1234567890123456789)
	data, err := value.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(data) != `"1234567890123456789"` {
		t.Fatalf("unexpected JSON: %s", string(data))
	}

	var decoded JsonInt64
	if err := decoded.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if decoded != value {
		t.Fatalf("round trip mismatch: got=%d want=%d", decoded, value)
	}
}

func TestJsonInt64UnmarshalJSONNull(t *testing.T) {
	var decoded JsonInt64 = 99
	if err := decoded.UnmarshalJSON([]byte("null")); err != nil {
		t.Fatalf("UnmarshalJSON(null) failed: %v", err)
	}
	if decoded != 0 {
		t.Fatalf("null should reset JsonInt64 to zero, got %d", decoded)
	}
}

func TestJsonInt64UnmarshalJSONTrimmedString(t *testing.T) {
	var decoded JsonInt64
	if err := decoded.UnmarshalJSON([]byte("  \"42\"  ")); err != nil {
		t.Fatalf("UnmarshalJSON(trimmed string) failed: %v", err)
	}
	if decoded != 42 {
		t.Fatalf("trimmed string should decode to 42, got %d", decoded)
	}
}

func TestJsonInt64ScanNil(t *testing.T) {
	var decoded JsonInt64 = 99
	if err := decoded.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) failed: %v", err)
	}
	if decoded != 0 {
		t.Fatalf("nil should reset JsonInt64 to zero, got %d", decoded)
	}
}

func TestJsonInt64ScanTrimmedText(t *testing.T) {
	var decoded JsonInt64
	if err := decoded.Scan([]byte(" 77 ")); err != nil {
		t.Fatalf("Scan(trimmed bytes) failed: %v", err)
	}
	if decoded != 77 {
		t.Fatalf("trimmed bytes should decode to 77, got %d", decoded)
	}
}

func TestTrySetWorkerID(t *testing.T) {
	original := defaultGen.workerID
	t.Cleanup(func() {
		defaultGen.workerID = original
	})

	if err := TrySetWorkerID(7); err != nil {
		t.Fatalf("TrySetWorkerID(valid) failed: %v", err)
	}
	if defaultGen.workerID != 7 {
		t.Fatalf("worker id mismatch after valid set: %d", defaultGen.workerID)
	}

	if err := TrySetWorkerID(workerMax + 1); err == nil {
		t.Fatal("TrySetWorkerID should reject out-of-range value")
	}
	if defaultGen.workerID != 7 {
		t.Fatalf("worker id should remain unchanged after invalid set: %d", defaultGen.workerID)
	}
}
