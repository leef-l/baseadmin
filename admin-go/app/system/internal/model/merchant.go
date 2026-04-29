package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// MerchantCreateInput 创建商户输入
type MerchantCreateInput struct {
	TenantID      snowflake.JsonInt64 `json:"tenantId"`
	Name          string              `json:"name"`
	Code          string              `json:"code"`
	ContactName   string              `json:"contactName"`
	ContactPhone  string              `json:"contactPhone"`
	Address       string              `json:"address"`
	Status        int                 `json:"status"`
	Remark        string              `json:"remark"`
	CreateAdmin   int                 `json:"createAdmin"`
	AdminUsername string              `json:"adminUsername"`
	AdminPassword string              `json:"adminPassword"`
	AdminNickname string              `json:"adminNickname"`
	AdminEmail    string              `json:"adminEmail"`
}

// MerchantUpdateInput 更新商户输入
type MerchantUpdateInput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	Name         string              `json:"name"`
	Code         string              `json:"code"`
	ContactName  string              `json:"contactName"`
	ContactPhone string              `json:"contactPhone"`
	Address      string              `json:"address"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
}

// MerchantDetailOutput 商户详情输出
type MerchantDetailOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	TenantName   string              `json:"tenantName"`
	Name         string              `json:"name"`
	Code         string              `json:"code"`
	ContactName  string              `json:"contactName"`
	ContactPhone string              `json:"contactPhone"`
	Address      string              `json:"address"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// MerchantListOutput 商户列表输出
type MerchantListOutput = MerchantDetailOutput

// MerchantListInput 商户列表查询输入
type MerchantListInput struct {
	PageNum  int                 `json:"pageNum"`
	PageSize int                 `json:"pageSize"`
	Keyword  string              `json:"keyword"`
	TenantID snowflake.JsonInt64 `json:"tenantId"`
	Code     string              `json:"code"`
	Status   *int                `json:"status"`
}
