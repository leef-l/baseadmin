package fieldvalid

import "testing"

func TestEnum(t *testing.T) {
	if err := Enum("菜单类型", 2, 1, 2, 3); err != nil {
		t.Fatalf("Enum should allow valid value: %v", err)
	}
	if err := Enum("菜单类型", 9, 1, 2, 3); err == nil || err.Error() != "菜单类型值不合法" {
		t.Fatalf("Enum invalid mismatch: %v", err)
	}
}

func TestBinary(t *testing.T) {
	if err := Binary("状态", 1); err != nil {
		t.Fatalf("Binary should allow valid value: %v", err)
	}
	if err := Binary("状态", 2); err == nil || err.Error() != "状态值不合法" {
		t.Fatalf("Binary invalid mismatch: %v", err)
	}
}

func TestNonNegative(t *testing.T) {
	if err := NonNegative("排序", 0); err != nil {
		t.Fatalf("NonNegative should allow zero: %v", err)
	}
	if err := NonNegative("排序", -1); err == nil || err.Error() != "排序不能小于0" {
		t.Fatalf("NonNegative invalid mismatch: %v", err)
	}
}

func TestNonNegative64(t *testing.T) {
	if err := NonNegative64("文件大小", 12); err != nil {
		t.Fatalf("NonNegative64 should allow positive value: %v", err)
	}
	if err := NonNegative64("文件大小", -3); err == nil || err.Error() != "文件大小不能小于0" {
		t.Fatalf("NonNegative64 invalid mismatch: %v", err)
	}
}
