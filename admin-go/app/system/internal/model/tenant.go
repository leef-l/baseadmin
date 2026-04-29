package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// TenantCreateInput 创建租户输入
type TenantCreateInput struct {
	Name          string      `json:"name"`
	Code          string      `json:"code"`
	ContactName   string      `json:"contactName"`
	ContactPhone  string      `json:"contactPhone"`
	Domain        string      `json:"domain"`
	ExpireAt      *gtime.Time `json:"expireAt"`
	Status        int         `json:"status"`
	Remark        string      `json:"remark"`
	CreateAdmin   int         `json:"createAdmin"`
	AdminUsername string      `json:"adminUsername"`
	AdminPassword string      `json:"adminPassword"`
	AdminNickname string      `json:"adminNickname"`
	AdminEmail    string      `json:"adminEmail"`
}

// TenantUpdateInput 更新租户输入
type TenantUpdateInput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Code         string              `json:"code"`
	ContactName  string              `json:"contactName"`
	ContactPhone string              `json:"contactPhone"`
	Domain       string              `json:"domain"`
	ExpireAt     *gtime.Time         `json:"expireAt"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
}

// TenantDetailOutput 租户详情输出
type TenantDetailOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Code         string              `json:"code"`
	ContactName  string              `json:"contactName"`
	ContactPhone string              `json:"contactPhone"`
	Domain       string              `json:"domain"`
	ExpireAt     *gtime.Time         `json:"expireAt"`
	Status       int                 `json:"status"`
	Remark       string              `json:"remark"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// TenantListOutput 租户列表输出
type TenantListOutput = TenantDetailOutput

// TenantListInput 租户列表查询输入
type TenantListInput struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
	Code     string `json:"code"`
	Status   *int   `json:"status"`
}
