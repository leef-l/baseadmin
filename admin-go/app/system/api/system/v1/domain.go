package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// DomainCreateReq 创建域名绑定请求
type DomainCreateReq struct {
	g.Meta       `path:"/domain/create" method:"post" tags:"域名绑定" summary:"创建域名绑定"`
	Domain       string              `json:"domain" v:"required#域名不能为空" dc:"绑定域名"`
	OwnerType    int                 `json:"ownerType" dc:"主体类型:1=租户,2=商户"`
	TenantID     snowflake.JsonInt64 `json:"tenantId" dc:"租户ID"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId" dc:"商户ID"`
	AppCode      string              `json:"appCode" dc:"应用编码"`
	VerifyStatus int                 `json:"verifyStatus" dc:"校验状态"`
	SslStatus    int                 `json:"sslStatus" dc:"SSL状态"`
	Status       int                 `json:"status" dc:"状态"`
	Remark       string              `json:"remark" dc:"备注"`
}

type DomainCreateRes struct {
	g.Meta `mime:"application/json"`
}

// DomainUpdateReq 更新域名绑定请求
type DomainUpdateReq struct {
	g.Meta       `path:"/domain/update" method:"put" tags:"域名绑定" summary:"更新域名绑定"`
	ID           snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"域名ID"`
	Domain       string              `json:"domain" v:"required#域名不能为空" dc:"绑定域名"`
	OwnerType    int                 `json:"ownerType" dc:"主体类型:1=租户,2=商户"`
	TenantID     snowflake.JsonInt64 `json:"tenantId" dc:"租户ID"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId" dc:"商户ID"`
	AppCode      string              `json:"appCode" dc:"应用编码"`
	VerifyStatus int                 `json:"verifyStatus" dc:"校验状态"`
	SslStatus    int                 `json:"sslStatus" dc:"SSL状态"`
	Status       int                 `json:"status" dc:"状态"`
	Remark       string              `json:"remark" dc:"备注"`
}

type DomainUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// DomainDeleteReq 删除域名绑定请求
type DomainDeleteReq struct {
	g.Meta `path:"/domain/delete" method:"delete" tags:"域名绑定" summary:"删除域名绑定"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"域名ID"`
}

type DomainDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DomainBatchDeleteReq 批量删除域名绑定请求
type DomainBatchDeleteReq struct {
	g.Meta `path:"/domain/batch-delete" method:"delete" tags:"域名绑定" summary:"批量删除域名绑定"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"域名ID列表"`
}

type DomainBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DomainDetailReq 获取域名绑定详情请求
type DomainDetailReq struct {
	g.Meta `path:"/domain/detail" method:"get" tags:"域名绑定" summary:"获取域名绑定详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"域名ID"`
}

type DomainDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.DomainDetailOutput
}

// DomainListReq 获取域名绑定列表请求
type DomainListReq struct {
	g.Meta     `path:"/domain/list" method:"get" tags:"域名绑定" summary:"获取域名绑定列表"`
	PageNum    int                 `json:"pageNum" d:"1" dc:"页码"`
	PageSize   int                 `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword    string              `json:"keyword" dc:"关键词"`
	Domain     string              `json:"domain" dc:"绑定域名"`
	OwnerType  int                 `json:"ownerType" dc:"主体类型"`
	TenantID   snowflake.JsonInt64 `json:"tenantId" dc:"租户ID"`
	MerchantID snowflake.JsonInt64 `json:"merchantId" dc:"商户ID"`
	AppCode    string              `json:"appCode" dc:"应用编码"`
	Status     *int                `json:"status" dc:"状态"`
}

type DomainListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.DomainListOutput `json:"list" dc:"列表数据"`
	Total  int                       `json:"total" dc:"总数"`
}

// DomainApplyNginxReq 应用 Nginx 配置请求
type DomainApplyNginxReq struct {
	g.Meta `path:"/domain/apply-nginx" method:"post" tags:"域名绑定" summary:"应用Nginx配置"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"域名ID"`
}

type DomainApplyNginxRes struct {
	g.Meta `mime:"application/json"`
	*model.DomainApplyNginxOutput
}

// DomainApplySSLReq 申请 SSL 证书请求
type DomainApplySSLReq struct {
	g.Meta `path:"/domain/apply-ssl" method:"post" tags:"域名绑定" summary:"申请SSL证书"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"域名ID"`
}

type DomainApplySSLRes struct {
	g.Meta `mime:"application/json"`
	*model.DomainApplySSLOutput
}
