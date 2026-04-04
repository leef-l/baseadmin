package file

import (
	"context"

	v1 "gbaseadmin/app/upload/api/upload/v1"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
)

var File = cFile{}

type cFile struct{}

// Create 创建文件记录
func (c *cFile) Create(ctx context.Context, req *v1.FileCreateReq) (res *v1.FileCreateRes, err error) {
	err = service.File().Create(ctx, &model.FileCreateInput{
		DirID:   req.DirID,
		Name:    req.Name,
		URL:     req.URL,
		Ext:     req.Ext,
		Size:    req.Size,
		Mime:    req.Mime,
		Storage: req.Storage,
		IsImage: req.IsImage,
	})
	return
}

// Update 更新文件记录
func (c *cFile) Update(ctx context.Context, req *v1.FileUpdateReq) (res *v1.FileUpdateRes, err error) {
	err = service.File().Update(ctx, &model.FileUpdateInput{
		ID:      req.ID,
		DirID:   req.DirID,
		Name:    req.Name,
		URL:     req.URL,
		Ext:     req.Ext,
		Size:    req.Size,
		Mime:    req.Mime,
		Storage: req.Storage,
		IsImage: req.IsImage,
	})
	return
}

// Delete 删除文件记录
func (c *cFile) Delete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error) {
	err = service.File().Delete(ctx, req.ID)
	return
}

// Detail 获取文件记录详情
func (c *cFile) Detail(ctx context.Context, req *v1.FileDetailReq) (res *v1.FileDetailRes, err error) {
	res = &v1.FileDetailRes{}
	res.FileDetailOutput, err = service.File().Detail(ctx, req.ID)
	return
}

// List 获取文件记录列表
func (c *cFile) List(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error) {
	res = &v1.FileListRes{}
	res.List, res.Total, err = service.File().List(ctx, &model.FileListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		DirID:    req.DirID,
		Name:     req.Name,
		Storage:  req.Storage,
		IsImage:  req.IsImage,
	})
	return
}
