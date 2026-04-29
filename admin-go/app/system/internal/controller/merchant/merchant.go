package merchant

import (
	"context"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
)

var Merchant = cMerchant{}

type cMerchant struct{}

func (c *cMerchant) Create(ctx context.Context, req *v1.MerchantCreateReq) (res *v1.MerchantCreateRes, err error) {
	err = service.Merchant().Create(ctx, &model.MerchantCreateInput{
		TenantID:      req.TenantID,
		Name:          req.Name,
		Code:          req.Code,
		ContactName:   req.ContactName,
		ContactPhone:  req.ContactPhone,
		Address:       req.Address,
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

func (c *cMerchant) Update(ctx context.Context, req *v1.MerchantUpdateReq) (res *v1.MerchantUpdateRes, err error) {
	err = service.Merchant().Update(ctx, &model.MerchantUpdateInput{
		ID:           req.ID,
		TenantID:     req.TenantID,
		Name:         req.Name,
		Code:         req.Code,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
		Status:       req.Status,
		Remark:       req.Remark,
	})
	return
}

func (c *cMerchant) Delete(ctx context.Context, req *v1.MerchantDeleteReq) (res *v1.MerchantDeleteRes, err error) {
	err = service.Merchant().Delete(ctx, req.ID)
	return
}

func (c *cMerchant) BatchDelete(ctx context.Context, req *v1.MerchantBatchDeleteReq) (res *v1.MerchantBatchDeleteRes, err error) {
	err = service.Merchant().BatchDelete(ctx, req.IDs)
	return
}

func (c *cMerchant) Detail(ctx context.Context, req *v1.MerchantDetailReq) (res *v1.MerchantDetailRes, err error) {
	res = &v1.MerchantDetailRes{}
	res.MerchantDetailOutput, err = service.Merchant().Detail(ctx, req.ID)
	return
}

func (c *cMerchant) List(ctx context.Context, req *v1.MerchantListReq) (res *v1.MerchantListRes, err error) {
	res = &v1.MerchantListRes{}
	res.List, res.Total, err = service.Merchant().List(ctx, &model.MerchantListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		TenantID: req.TenantID,
		Code:     req.Code,
		Status:   req.Status,
	})
	return
}
