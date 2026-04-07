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
