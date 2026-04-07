package treeutil

import (
	"gbaseadmin/utility/snowflake"

	"github.com/gogf/gf/v2/errors/gerror"
)

type Messages struct {
	Self         string
	Missing      string
	ChildLoop    string
	Cycle        string
	InvalidChain string
}

type ParentLoader func(id int64) (nodeID int64, parentID int64, err error)

type TreeNodeAccessor[T any] struct {
	ID       func(item T) int64
	ParentID func(item T) int64
	Init     func(item T)
	Append   func(parent T, child T)
}

func ValidateParent(parentID, currentID snowflake.JsonInt64, load ParentLoader, messages Messages) error {
	if parentID == 0 {
		return nil
	}
	if currentID != 0 && parentID == currentID {
		return gerror.New(messages.Self)
	}
	if load == nil {
		return gerror.New(messages.Missing)
	}

	nodeID, nextParentID, err := load(int64(parentID))
	if err != nil {
		return err
	}
	if nodeID == 0 {
		return gerror.New(messages.Missing)
	}

	seen := map[int64]struct{}{int64(parentID): {}}
	for nextParentID != 0 {
		if currentID != 0 && nextParentID == int64(currentID) {
			return gerror.New(messages.ChildLoop)
		}
		if _, ok := seen[nextParentID]; ok {
			return gerror.New(messages.Cycle)
		}
		seen[nextParentID] = struct{}{}

		nodeID, nextParentID, err = load(nextParentID)
		if err != nil {
			return err
		}
		if nodeID == 0 {
			return gerror.New(messages.InvalidChain)
		}
	}
	return nil
}

func BuildForest[T any](list []T, accessor TreeNodeAccessor[T]) []T {
	if len(list) == 0 {
		return make([]T, 0)
	}
	if accessor.ID == nil || accessor.ParentID == nil || accessor.Append == nil {
		return list
	}

	nodeMap := make(map[int64]T, len(list))
	for _, item := range list {
		if accessor.Init != nil {
			accessor.Init(item)
		}
		nodeMap[accessor.ID(item)] = item
	}

	tree := make([]T, 0, len(list))
	for _, item := range list {
		parentID := accessor.ParentID(item)
		if parentID == 0 {
			tree = append(tree, item)
			continue
		}
		parent, ok := nodeMap[parentID]
		if !ok {
			tree = append(tree, item)
			continue
		}
		accessor.Append(parent, item)
	}
	return tree
}
