package domain

import (
	"context"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
)

var Domain = cDomain{}

type cDomain struct{}

func (c *cDomain) Create(ctx context.Context, req *v1.DomainCreateReq) (res *v1.DomainCreateRes, err error) {
	err = service.Domain().Create(ctx, &model.DomainCreateInput{
		Domain:       req.Domain,
		OwnerType:    req.OwnerType,
		TenantID:     req.TenantID,
		MerchantID:   req.MerchantID,
		AppCode:      req.AppCode,
		VerifyStatus: req.VerifyStatus,
		SslStatus:    req.SslStatus,
		Status:       req.Status,
		Remark:       req.Remark,
	})
	return
}

func (c *cDomain) Update(ctx context.Context, req *v1.DomainUpdateReq) (res *v1.DomainUpdateRes, err error) {
	err = service.Domain().Update(ctx, &model.DomainUpdateInput{
		ID:           req.ID,
		Domain:       req.Domain,
		OwnerType:    req.OwnerType,
		TenantID:     req.TenantID,
		MerchantID:   req.MerchantID,
		AppCode:      req.AppCode,
		VerifyStatus: req.VerifyStatus,
		SslStatus:    req.SslStatus,
		Status:       req.Status,
		Remark:       req.Remark,
	})
	return
}

func (c *cDomain) Delete(ctx context.Context, req *v1.DomainDeleteReq) (res *v1.DomainDeleteRes, err error) {
	err = service.Domain().Delete(ctx, req.ID)
	return
}

func (c *cDomain) BatchDelete(ctx context.Context, req *v1.DomainBatchDeleteReq) (res *v1.DomainBatchDeleteRes, err error) {
	err = service.Domain().BatchDelete(ctx, req.IDs)
	return
}

func (c *cDomain) Detail(ctx context.Context, req *v1.DomainDetailReq) (res *v1.DomainDetailRes, err error) {
	res = &v1.DomainDetailRes{}
	res.DomainDetailOutput, err = service.Domain().Detail(ctx, req.ID)
	return
}

func (c *cDomain) List(ctx context.Context, req *v1.DomainListReq) (res *v1.DomainListRes, err error) {
	res = &v1.DomainListRes{}
	res.List, res.Total, err = service.Domain().List(ctx, &model.DomainListInput{
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		Keyword:    req.Keyword,
		Domain:     req.Domain,
		OwnerType:  req.OwnerType,
		TenantID:   req.TenantID,
		MerchantID: req.MerchantID,
		AppCode:    req.AppCode,
		Status:     req.Status,
	})
	return
}

func (c *cDomain) ApplyNginx(ctx context.Context, req *v1.DomainApplyNginxReq) (res *v1.DomainApplyNginxRes, err error) {
	res = &v1.DomainApplyNginxRes{}
	res.DomainApplyNginxOutput, err = service.Domain().ApplyNginx(ctx, req.ID)
	return
}

func (c *cDomain) ApplySSL(ctx context.Context, req *v1.DomainApplySSLReq) (res *v1.DomainApplySSLRes, err error) {
	res = &v1.DomainApplySSLRes{}
	res.DomainApplySSLOutput, err = service.Domain().ApplySSL(ctx, req.ID)
	return
}
