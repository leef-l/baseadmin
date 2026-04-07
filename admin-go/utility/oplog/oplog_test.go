package oplog

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func TestRecordDispatchesAsyncInsert(t *testing.T) {
	original := insertOperationLog
	defer func() {
		insertOperationLog = original
	}()

	ch := make(chan g.Map, 1)
	insertOperationLog = func(ctx context.Context, data g.Map) {
		ch <- data
	}

	Record(context.Background(), "order", "create", "1001", "demo")

	select {
	case data := <-ch:
		if data["module"] != "order" || data["action"] != "create" || data["target_id"] != "1001" || data["detail"] != "demo" {
			t.Fatalf("unexpected oplog payload: %+v", data)
		}
		if data["created_at"] == nil {
			t.Fatalf("created_at should be populated: %+v", data)
		}
	case <-time.After(time.Second):
		t.Fatal("Record did not dispatch async insert in time")
	}
}

func TestRecordSkipsBlankModuleOrAction(t *testing.T) {
	original := insertOperationLog
	defer func() {
		insertOperationLog = original
	}()

	called := false
	insertOperationLog = func(ctx context.Context, data g.Map) {
		called = true
	}

	Record(context.Background(), " ", "create", "1001", "demo")
	Record(context.Background(), "order", " ", "1001", "demo")
	time.Sleep(20 * time.Millisecond)
	if called {
		t.Fatal("Record should skip blank module/action")
	}
}
