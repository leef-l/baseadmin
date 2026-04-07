package inpututil

import "testing"

func TestRequire(t *testing.T) {
	if err := Require(nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Require(nil) mismatch: %v", err)
	}
	if err := Require(struct{}{}); err != nil {
		t.Fatalf("Require(non-nil) should succeed: %v", err)
	}
}
