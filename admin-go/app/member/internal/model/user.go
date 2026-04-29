package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// User DTO 模型

// UserCreateInput 创建会员用户输入
type UserCreateInput struct {
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	RealName string `json:"realName"`
	LevelID snowflake.JsonInt64 `json:"levelID"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"`
	TeamCount int `json:"teamCount"`
	DirectCount int `json:"directCount"`
	ActiveCount int `json:"activeCount"`
	TeamTurnover int64 `json:"teamTurnover"`
	IsActive int `json:"isActive"`
	IsQualified int `json:"isQualified"`
	InviteCode string `json:"inviteCode"`
	RegisterIP string `json:"registerIP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// UserUpdateInput 更新会员用户输入
type UserUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	RealName string `json:"realName"`
	LevelID snowflake.JsonInt64 `json:"levelID"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"`
	TeamCount int `json:"teamCount"`
	DirectCount int `json:"directCount"`
	ActiveCount int `json:"activeCount"`
	TeamTurnover int64 `json:"teamTurnover"`
	IsActive int `json:"isActive"`
	IsQualified int `json:"isQualified"`
	InviteCode string `json:"inviteCode"`
	RegisterIP string `json:"registerIP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// UserDetailOutput 会员用户详情输出
type UserDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	UserUsername string `json:"userUsername"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	RealName string `json:"realName"`
	LevelID snowflake.JsonInt64 `json:"levelID"`
	LevelName string `json:"levelName"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"`
	TeamCount int `json:"teamCount"`
	DirectCount int `json:"directCount"`
	ActiveCount int `json:"activeCount"`
	TeamTurnover int64 `json:"teamTurnover"`
	IsActive int `json:"isActive"`
	IsQualified int `json:"isQualified"`
	InviteCode string `json:"inviteCode"`
	RegisterIP string `json:"registerIP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// UserListOutput 会员用户列表输出
type UserListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	UserUsername string `json:"userUsername"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	RealName string `json:"realName"`
	LevelID snowflake.JsonInt64 `json:"levelID"`
	LevelName string `json:"levelName"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"`
	TeamCount int `json:"teamCount"`
	DirectCount int `json:"directCount"`
	ActiveCount int `json:"activeCount"`
	TeamTurnover int64 `json:"teamTurnover"`
	IsActive int `json:"isActive"`
	IsQualified int `json:"isQualified"`
	InviteCode string `json:"inviteCode"`
	RegisterIP string `json:"registerIP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// UserListInput 会员用户列表查询输入
type UserListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	Username string `json:"username"`
	InviteCode string `json:"inviteCode"`
	Nickname string `json:"nickname"`
	RealName string `json:"realName"`
	Phone string `json:"phone"`
	ParentID *snowflake.JsonInt64 `json:"parentID"`
	LevelID *snowflake.JsonInt64 `json:"levelID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	IsActive *int `json:"isActive"`
	IsQualified *int `json:"isQualified"`
	Status *int `json:"status"`
	LevelExpireAtStart string `json:"levelExpireAtStart"`
	LevelExpireAtEnd string `json:"levelExpireAtEnd"`
	LastLoginAtStart string `json:"lastLoginAtStart"`
	LastLoginAtEnd string `json:"lastLoginAtEnd"`
}

// UserTreeInput 会员用户树形查询输入
type UserTreeInput struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	Username string `json:"username"`
	InviteCode string `json:"inviteCode"`
	Nickname string `json:"nickname"`
	RealName string `json:"realName"`
	Phone string `json:"phone"`
	ParentID *snowflake.JsonInt64 `json:"parentID"`
	LevelID *snowflake.JsonInt64 `json:"levelID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	IsActive *int `json:"isActive"`
	IsQualified *int `json:"isQualified"`
	Status *int `json:"status"`
	LevelExpireAtStart string `json:"levelExpireAtStart"`
	LevelExpireAtEnd string `json:"levelExpireAtEnd"`
	LastLoginAtStart string `json:"lastLoginAtStart"`
	LastLoginAtEnd string `json:"lastLoginAtEnd"`
}

// UserTreeOutput 会员用户树形输出
type UserTreeOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	UserUsername string `json:"userUsername"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	RealName string `json:"realName"`
	LevelID snowflake.JsonInt64 `json:"levelID"`
	LevelName string `json:"levelName"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"`
	TeamCount int `json:"teamCount"`
	DirectCount int `json:"directCount"`
	ActiveCount int `json:"activeCount"`
	TeamTurnover int64 `json:"teamTurnover"`
	IsActive int `json:"isActive"`
	IsQualified int `json:"isQualified"`
	InviteCode string `json:"inviteCode"`
	RegisterIP string `json:"registerIP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time                `json:"createdAt"`
	UpdatedAt *gtime.Time                `json:"updatedAt"`
	Children  []*UserTreeOutput `json:"children"`
}

// UserBatchUpdateInput 批量编辑会员用户输入
type UserBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	IsActive *int `json:"isActive"`
	IsQualified *int `json:"isQualified"`
	Status *int `json:"status"`
}

