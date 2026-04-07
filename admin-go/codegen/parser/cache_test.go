package parser

import (
	"errors"
	"testing"
)

func TestGetTableColumnsCachesResults(t *testing.T) {
	calls := 0
	p := &Parser{
		tableColumnsCache: make(map[string]map[string]struct{}),
		tableColumnsLoader: func(tableName string) (map[string]struct{}, error) {
			calls++
			if tableName != "system_role" {
				t.Fatalf("unexpected table name: %s", tableName)
			}
			return map[string]struct{}{
				"id":        {},
				"name":      {},
				"parent_id": {},
			}, nil
		},
	}

	if _, err := p.getTableColumns("system_role"); err != nil {
		t.Fatalf("getTableColumns failed: %v", err)
	}
	if _, err := p.getTableColumns("system_role"); err != nil {
		t.Fatalf("getTableColumns second pass failed: %v", err)
	}
	if !p.tableHasColumn("system_role", "parent_id") {
		t.Fatal("expected cached column lookup to succeed")
	}
	if got := p.findDisplayField("system_role"); got != "name" {
		t.Fatalf("findDisplayField mismatch: %q", got)
	}
	if calls != 1 {
		t.Fatalf("expected one loader call, got %d", calls)
	}
}

func TestGetTableColumnsCachesMissingTables(t *testing.T) {
	calls := 0
	p := &Parser{
		tableColumnsCache: make(map[string]map[string]struct{}),
		tableColumnsLoader: func(tableName string) (map[string]struct{}, error) {
			calls++
			return map[string]struct{}{}, nil
		},
	}

	if got := p.findDisplayField("missing_table"); got != "" {
		t.Fatalf("missing table should not resolve display field, got %q", got)
	}
	if p.tableHasColumn("missing_table", "parent_id") {
		t.Fatal("missing table should not report existing columns")
	}
	if calls != 1 {
		t.Fatalf("expected missing table lookup to be cached, got %d loader calls", calls)
	}
}

func TestGetTableColumnsDoesNotCacheErrors(t *testing.T) {
	calls := 0
	p := &Parser{
		tableColumnsCache: make(map[string]map[string]struct{}),
		tableColumnsLoader: func(tableName string) (map[string]struct{}, error) {
			calls++
			return nil, errors.New("boom")
		},
	}

	if _, err := p.getTableColumns("broken_table"); err == nil {
		t.Fatal("expected loader error on first call")
	}
	if _, err := p.getTableColumns("broken_table"); err == nil {
		t.Fatal("expected loader error on second call")
	}
	if calls != 2 {
		t.Fatalf("loader errors should not be cached, got %d calls", calls)
	}
}

func TestFindDisplayFieldUsesPriorityOrder(t *testing.T) {
	p := &Parser{
		tableColumnsCache: map[string]map[string]struct{}{
			"system_users": {
				"email": {},
				"title": {},
				"name":  {},
			},
		},
	}

	if got := p.findDisplayField("system_users"); got != "title" {
		t.Fatalf("expected highest-priority display field, got %q", got)
	}
}

func TestTableHasColumnReflectsDeletedAtPresence(t *testing.T) {
	p := &Parser{
		tableColumnsCache: map[string]map[string]struct{}{
			"system_users": {
				"id":         {},
				"username":   {},
				"deleted_at": {},
			},
			"legacy_member": {
				"id":       {},
				"nickname": {},
			},
		},
	}

	if !p.tableHasColumn("system_users", "deleted_at") {
		t.Fatal("expected deleted_at on system_users")
	}
	if p.tableHasColumn("legacy_member", "deleted_at") {
		t.Fatal("did not expect deleted_at on legacy_member")
	}
}
