package dir

import (
	"context"

	v1 "gbaseadmin/app/upload/api/upload/v1"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
)

var Dir = cDir{}

type cDir struct{}

// Create 创建文件目录
func (c *cDir) Create(ctx context.Context, req *v1.DirCreateReq) (res *v1.DirCreateRes, err error) {
	err = service.Dir().Create(ctx, &model.DirCreateInput{
		ParentID: req.ParentID,
		Name:     req.Name,
		Path:     req.Path,
		KeepName: req.KeepName,
		Sort:     req.Sort,
		Status:   req.Status,
	})
	return
}

// Update 更新文件目录
func (c *cDir) Update(ctx context.Context, req *v1.DirUpdateReq) (res *v1.DirUpdateRes, err error) {
	err = service.Dir().Update(ctx, &model.DirUpdateInput{
		ID:       req.ID,
		ParentID: req.ParentID,
		Name:     req.Name,
		Path:     req.Path,
		KeepName: req.KeepName,
		Sort:     req.Sort,
		Status:   req.Status,
	})
	return
}

// Delete 删除文件目录
func (c *cDir) Delete(ctx context.Context, req *v1.DirDeleteReq) (res *v1.DirDeleteRes, err error) {
	err = service.Dir().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除文件目录
func (c *cDir) BatchDelete(ctx context.Context, req *v1.DirBatchDeleteReq) (res *v1.DirBatchDeleteRes, err error) {
	err = service.Dir().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取文件目录详情
func (c *cDir) Detail(ctx context.Context, req *v1.DirDetailReq) (res *v1.DirDetailRes, err error) {
	res = &v1.DirDetailRes{}
	res.DirDetailOutput, err = service.Dir().Detail(ctx, req.ID)
	return
}

// List 获取文件目录列表
func (c *cDir) List(ctx context.Context, req *v1.DirListReq) (res *v1.DirListRes, err error) {
	res = &v1.DirListRes{}
	res.List, res.Total, err = service.Dir().List(ctx, &model.DirListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   req.Status,
	})
	return
}

// Tree 获取文件目录树形结构
func (c *cDir) Tree(ctx context.Context, req *v1.DirTreeReq) (res *v1.DirTreeRes, err error) {
	res = &v1.DirTreeRes{}
	res.List, err = service.Dir().Tree(ctx, &model.DirTreeInput{
		Keyword: req.Keyword,
		Status:  req.Status,
	})
	return
}
