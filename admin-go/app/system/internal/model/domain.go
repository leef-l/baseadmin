package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// DomainCreateInput 创建域名绑定输入
type DomainCreateInput struct {
	Domain       string              `json:"domain"`
	OwnerType    int                 `json:"ownerType"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId"`
	AppCode      string              `json:"appCode"`
	VerifyStatus int                 `json:"verifyStatus"`
	SslStatus    int                 `json:"sslStatus"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
}

// DomainUpdateInput 更新域名绑定输入
type DomainUpdateInput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Domain       string              `json:"domain"`
	OwnerType    int                 `json:"ownerType"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId"`
	AppCode      string              `json:"appCode"`
	VerifyStatus int                 `json:"verifyStatus"`
	SslStatus    int                 `json:"sslStatus"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
}

// DomainDetailOutput 域名绑定详情输出
type DomainDetailOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Domain       string              `json:"domain"`
	OwnerType    int                 `json:"ownerType"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	TenantName   string              `json:"tenantName"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId"`
	MerchantName string              `json:"merchantName"`
	AppCode      string              `json:"appCode"`
	VerifyToken  string              `json:"verifyToken"`
	VerifyStatus int                 `json:"verifyStatus"`
	SslStatus    int                 `json:"sslStatus"`
	NginxStatus  int                 `json:"nginxStatus"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// DomainListOutput 域名绑定列表输出（不含 VerifyToken）
type DomainListOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Domain       string              `json:"domain"`
	OwnerType    int                 `json:"ownerType"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	TenantName   string              `json:"tenantName"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId"`
	MerchantName string              `json:"merchantName"`
	AppCode      string              `json:"appCode"`
	VerifyStatus int                 `json:"verifyStatus"`
	SslStatus    int                 `json:"sslStatus"`
	NginxStatus  int                 `json:"nginxStatus"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// DomainListInput 域名绑定列表查询输入
type DomainListInput struct {
	PageNum    int                 `json:"pageNum"`
	PageSize   int                 `json:"pageSize"`
	Keyword    string              `json:"keyword"`
	Domain     string              `json:"domain"`
	OwnerType  int                 `json:"ownerType"`
	TenantID   snowflake.JsonInt64 `json:"tenantId"`
	MerchantID snowflake.JsonInt64 `json:"merchantId"`
	AppCode    string              `json:"appCode"`
	Status     *int                `json:"status"`
}

// DomainApplyNginxOutput 应用 Nginx 配置输出
type DomainApplyNginxOutput struct {
	ConfigPath  string `json:"configPath"`
	NginxStatus int    `json:"nginxStatus"`
	SslStatus   int    `json:"sslStatus"`
}

// DomainApplySSLOutput 申请 SSL 证书输出
type DomainApplySSLOutput struct {
	ConfigPath  string `json:"configPath"`
	CertPath    string `json:"certPath"`
	NginxStatus int    `json:"nginxStatus"`
	SslStatus   int    `json:"sslStatus"`
}
