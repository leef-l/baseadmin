package tenant

import (
	"context"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
)

var Tenant = cTenant{}

type cTenant struct{}

func (c *cTenant) Create(ctx context.Context, req *v1.TenantCreateReq) (res *v1.TenantCreateRes, err error) {
	err = service.Tenant().Create(ctx, &model.TenantCreateInput{
		Name:          req.Name,
		Code:          req.Code,
		ContactName:   req.ContactName,
		ContactPhone:  req.ContactPhone,
		Domain:        req.Domain,
		ExpireAt:      req.ExpireAt,
		Status:        req.Status,
		Remark:        req.Remark,
		CreateAdmin:   req.CreateAdmin,
		AdminUsername: req.AdminUsername,
		AdminPassword: req.AdminPassword,
		AdminNickname: req.AdminNickname,
		AdminEmail:    req.AdminEmail,
	})
	return
}

func (c *cTenant) Update(ctx context.Context, req *v1.TenantUpdateReq) (res *v1.TenantUpdateRes, err error) {
	err = service.Tenant().Update(ctx, &model.TenantUpdateInput{
		ID:           req.ID,
		Name:         req.Name,
		Code:         req.Code,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Domain:       req.Domain,
		ExpireAt:     req.ExpireAt,
		Status:       req.Status,
		Remark:       req.Remark,
	})
	return
}

func (c *cTenant) Delete(ctx context.Context, req *v1.TenantDeleteReq) (res *v1.TenantDeleteRes, err error) {
	err = service.Tenant().Delete(ctx, req.ID)
	return
}

func (c *cTenant) BatchDelete(ctx context.Context, req *v1.TenantBatchDeleteReq) (res *v1.TenantBatchDeleteRes, err error) {
	err = service.Tenant().BatchDelete(ctx, req.IDs)
	return
}

func (c *cTenant) Detail(ctx context.Context, req *v1.TenantDetailReq) (res *v1.TenantDetailRes, err error) {
	res = &v1.TenantDetailRes{}
	res.TenantDetailOutput, err = service.Tenant().Detail(ctx, req.ID)
	return
}

func (c *cTenant) List(ctx context.Context, req *v1.TenantListReq) (res *v1.TenantListRes, err error) {
	res = &v1.TenantListRes{}
	res.List, res.Total, err = service.Tenant().List(ctx, &model.TenantListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Code:     req.Code,
		Status:   req.Status,
	})
	return
}
