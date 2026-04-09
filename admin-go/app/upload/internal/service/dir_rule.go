package service

import (
	"context"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IDirRule interface {
	Create(ctx context.Context, in *model.DirRuleCreateInput) error
	Update(ctx context.Context, in *model.DirRuleUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirRuleDetailOutput, err error)
	List(ctx context.Context, in *model.DirRuleListInput) (list []*model.DirRuleListOutput, total int, err error)
}

var localDirRule IDirRule

func DirRule() IDirRule {
	return localDirRule
}

func RegisterDirRule(i IDirRule) {
	localDirRule = i
}
