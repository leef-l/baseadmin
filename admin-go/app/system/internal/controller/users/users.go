package users

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/snowflake"
)

var Users = cUsers{}

type cUsers struct{}

// Create 创建用��表
func (c *cUsers) Create(ctx context.Context, req *v1.UsersCreateReq) (res *v1.UsersCreateRes, err error) {
	createdBy := snowflake.JsonInt64(g.RequestFromCtx(ctx).GetCtxVar("jwt_user_id").Int64())
	err = service.Users().Create(ctx, &model.UsersCreateInput{
		Username:   req.Username,
		Password:   req.Password,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Avatar:     req.Avatar,
		Status:     req.Status,
		DeptID:     req.DeptID,
		TenantID:   req.TenantID,
		MerchantID: req.MerchantID,
		CreatedBy:  createdBy,
		RoleIDs:    req.RoleIDs,
	})
	return
}

// Update 更新用户表
func (c *cUsers) Update(ctx context.Context, req *v1.UsersUpdateReq) (res *v1.UsersUpdateRes, err error) {
	detail, err := service.Users().Detail(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	err = service.Users().Update(ctx, &model.UsersUpdateInput{
		ID:         req.ID,
		Username:   pickStringField(req.Username, detail.Username),
		Password:   req.Password,
		Nickname:   pickStringField(req.Nickname, detail.Nickname),
		Email:      pickStringField(req.Email, detail.Email),
		Avatar:     pickStringField(req.Avatar, detail.Avatar),
		Status:     pickIntField(req.Status, detail.Status),
		DeptID:     pickSnowflakeField(req.DeptID, detail.DeptID),
		TenantID:   pickSnowflakeField(req.TenantID, detail.TenantID),
		MerchantID: pickSnowflakeField(req.MerchantID, detail.MerchantID),
		RoleIDs:    req.RoleIDs,
	})
	return
}

func pickStringField(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return *value
}

func pickIntField(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func pickSnowflakeField(value *snowflake.JsonInt64, fallback snowflake.JsonInt64) snowflake.JsonInt64 {
	if value == nil {
		return fallback
	}
	return *value
}

// Delete 删除用户表
func (c *cUsers) Delete(ctx context.Context, req *v1.UsersDeleteReq) (res *v1.UsersDeleteRes, err error) {
	err = service.Users().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除用户表
func (c *cUsers) BatchDelete(ctx context.Context, req *v1.UsersBatchDeleteReq) (res *v1.UsersBatchDeleteRes, err error) {
	err = service.Users().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取用户表详情
func (c *cUsers) Detail(ctx context.Context, req *v1.UsersDetailReq) (res *v1.UsersDetailRes, err error) {
	res = &v1.UsersDetailRes{}
	res.UsersDetailOutput, err = service.Users().Detail(ctx, req.ID)
	return
}

// List 获取用户表列表
func (c *cUsers) List(ctx context.Context, req *v1.UsersListReq) (res *v1.UsersListRes, err error) {
	res = &v1.UsersListRes{}
	res.List, res.Total, err = service.Users().List(ctx, &model.UsersListInput{
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		Keyword:    req.Keyword,
		Username:   req.Username,
		Nickname:   req.Nickname,
		Email:      req.Email,
		DeptId:     req.DeptId,
		TenantId:   req.TenantId,
		MerchantId: req.MerchantId,
		Status:     req.Status,
	})
	return
}

// ResetPassword 重置用户密码
func (c *cUsers) ResetPassword(ctx context.Context, req *v1.UsersResetPasswordReq) (res *v1.UsersResetPasswordRes, err error) {
	err = service.Users().ResetPassword(ctx, &model.UsersResetPasswordInput{
		ID:       req.ID,
		Password: req.Password,
	})
	return
}
