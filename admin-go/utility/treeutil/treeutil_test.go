package treeutil

import (
	"errors"
	"testing"

	"gbaseadmin/utility/snowflake"
)

type testNode struct {
	ID       int64
	ParentID int64
	Children []*testNode
}

func TestValidateParent(t *testing.T) {
	msgs := Messages{
		Self:         "self",
		Missing:      "missing",
		ChildLoop:    "child",
		Cycle:        "cycle",
		InvalidChain: "invalid",
	}

	t.Run("self", func(t *testing.T) {
		err := ValidateParent(1, 1, nil, msgs)
		if err == nil || err.Error() != "self" {
			t.Fatalf("expected self error, got %v", err)
		}
	})

	t.Run("missing", func(t *testing.T) {
		err := ValidateParent(2, 0, func(id int64) (int64, int64, error) { return 0, 0, nil }, msgs)
		if err == nil || err.Error() != "missing" {
			t.Fatalf("expected missing error, got %v", err)
		}
	})

	t.Run("child loop", func(t *testing.T) {
		load := func(id int64) (int64, int64, error) {
			if id == 5 {
				return 5, 9, nil
			}
			return 0, 0, nil
		}
		err := ValidateParent(snowflake.JsonInt64(5), snowflake.JsonInt64(9), load, msgs)
		if err == nil || err.Error() != "child" {
			t.Fatalf("expected child loop error, got %v", err)
		}
	})

	t.Run("cycle", func(t *testing.T) {
		load := func(id int64) (int64, int64, error) {
			switch id {
			case 2:
				return 2, 3, nil
			case 3:
				return 3, 2, nil
			default:
				return 0, 0, nil
			}
		}
		err := ValidateParent(2, 0, load, msgs)
		if err == nil || err.Error() != "cycle" {
			t.Fatalf("expected cycle error, got %v", err)
		}
	})

	t.Run("invalid chain", func(t *testing.T) {
		load := func(id int64) (int64, int64, error) {
			if id == 2 {
				return 2, 3, nil
			}
			return 0, 0, nil
		}
		err := ValidateParent(2, 0, load, msgs)
		if err == nil || err.Error() != "invalid" {
			t.Fatalf("expected invalid chain error, got %v", err)
		}
	})

	t.Run("loader error", func(t *testing.T) {
		load := func(id int64) (int64, int64, error) {
			return 0, 0, errors.New("db error")
		}
		err := ValidateParent(2, 0, load, msgs)
		if err == nil || err.Error() != "db error" {
			t.Fatalf("expected loader error, got %v", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		load := func(id int64) (int64, int64, error) {
			switch id {
			case 2:
				return 2, 3, nil
			case 3:
				return 3, 0, nil
			default:
				return 0, 0, nil
			}
		}
		if err := ValidateParent(2, 9, load, msgs); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
	})
}

func TestBuildForest(t *testing.T) {
	root := &testNode{ID: 1}
	child := &testNode{ID: 2, ParentID: 1}
	orphan := &testNode{ID: 3, ParentID: 99}

	tree := BuildForest([]*testNode{root, child, orphan}, TreeNodeAccessor[*testNode]{
		ID:       func(item *testNode) int64 { return item.ID },
		ParentID: func(item *testNode) int64 { return item.ParentID },
		Init: func(item *testNode) {
			item.Children = make([]*testNode, 0)
		},
		Append: func(parent *testNode, child *testNode) {
			parent.Children = append(parent.Children, child)
		},
	})

	if len(tree) != 2 {
		t.Fatalf("BuildForest tree size mismatch: %d", len(tree))
	}
	if len(root.Children) != 1 || root.Children[0] != child {
		t.Fatalf("BuildForest should append child under root: %+v", root.Children)
	}
	if tree[1] != orphan {
		t.Fatalf("BuildForest should keep orphan on root level: %+v", tree[1])
	}
}
