package dir_rule

import (
	"context"

	v1 "gbaseadmin/app/upload/api/upload/v1"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
)

var DirRule = cDirRule{}

type cDirRule struct{}

// Create 创建文件目录规则
func (c *cDirRule) Create(ctx context.Context, req *v1.DirRuleCreateReq) (res *v1.DirRuleCreateRes, err error) {
	err = service.DirRule().Create(ctx, &model.DirRuleCreateInput{
		DirID:    req.DirID,
		Category: req.Category,
		SavePath: req.SavePath,
		Status:   req.Status,
	})
	return
}

// Update 更新文件目录规则
func (c *cDirRule) Update(ctx context.Context, req *v1.DirRuleUpdateReq) (res *v1.DirRuleUpdateRes, err error) {
	err = service.DirRule().Update(ctx, &model.DirRuleUpdateInput{
		ID:       req.ID,
		DirID:    req.DirID,
		Category: req.Category,
		SavePath: req.SavePath,
		Status:   req.Status,
	})
	return
}

// Delete 删除文件目录规则
func (c *cDirRule) Delete(ctx context.Context, req *v1.DirRuleDeleteReq) (res *v1.DirRuleDeleteRes, err error) {
	err = service.DirRule().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除文件目录规则
func (c *cDirRule) BatchDelete(ctx context.Context, req *v1.DirRuleBatchDeleteReq) (res *v1.DirRuleBatchDeleteRes, err error) {
	err = service.DirRule().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取文件目录规则详情
func (c *cDirRule) Detail(ctx context.Context, req *v1.DirRuleDetailReq) (res *v1.DirRuleDetailRes, err error) {
	res = &v1.DirRuleDetailRes{}
	res.DirRuleDetailOutput, err = service.DirRule().Detail(ctx, req.ID)
	return
}

// List 获取文件目录规则列表
func (c *cDirRule) List(ctx context.Context, req *v1.DirRuleListReq) (res *v1.DirRuleListRes, err error) {
	res = &v1.DirRuleListRes{}
	res.List, res.Total, err = service.DirRule().List(ctx, &model.DirRuleListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Category: req.Category,
		Status:   req.Status,
	})
	return
}
