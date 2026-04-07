package treeutil

import (
	"errors"
	"testing"

	"gbaseadmin/utility/snowflake"
)

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
