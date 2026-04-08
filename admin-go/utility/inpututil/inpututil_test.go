package inpututil

import "testing"

type sampleInput struct{}

func TestRequire(t *testing.T) {
	if err := Require(nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Require(nil) mismatch: %v", err)
	}
	var typedNil *sampleInput
	if err := Require(typedNil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Require(typed nil) mismatch: %v", err)
	}
	if err := Require(struct{}{}); err != nil {
		t.Fatalf("Require(non-nil) should succeed: %v", err)
	}
	if err := Require(&sampleInput{}); err != nil {
		t.Fatalf("Require(pointer) should succeed: %v", err)
	}
}
