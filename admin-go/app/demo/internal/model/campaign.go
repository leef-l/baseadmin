package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Campaign DTO 模型

// CampaignCreateInput 创建体验活动输入
type CampaignCreateInput struct {
	CampaignNo string `json:"campaignNo"`
	Title string `json:"title"`
	Banner string `json:"banner"`
	Type int `json:"type"`
	Channel int `json:"channel"`
	BudgetAmount int `json:"budgetAmount"`
	LandingURL string `json:"landingURL"`
	RuleJSON string `json:"ruleJSON"`
	IntroContent string `json:"introContent"`
	StartAt *gtime.Time `json:"startAt"`
	EndAt *gtime.Time `json:"endAt"`
	IsPublic int `json:"isPublic"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CampaignUpdateInput 更新体验活动输入
type CampaignUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CampaignNo string `json:"campaignNo"`
	Title string `json:"title"`
	Banner string `json:"banner"`
	Type int `json:"type"`
	Channel int `json:"channel"`
	BudgetAmount int `json:"budgetAmount"`
	LandingURL string `json:"landingURL"`
	RuleJSON string `json:"ruleJSON"`
	IntroContent string `json:"introContent"`
	StartAt *gtime.Time `json:"startAt"`
	EndAt *gtime.Time `json:"endAt"`
	IsPublic int `json:"isPublic"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CampaignDetailOutput 体验活动详情输出
type CampaignDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CampaignNo string `json:"campaignNo"`
	Title string `json:"title"`
	Banner string `json:"banner"`
	Type int `json:"type"`
	Channel int `json:"channel"`
	BudgetAmount int `json:"budgetAmount"`
	LandingURL string `json:"landingURL"`
	RuleJSON string `json:"ruleJSON"`
	IntroContent string `json:"introContent"`
	StartAt *gtime.Time `json:"startAt"`
	EndAt *gtime.Time `json:"endAt"`
	IsPublic int `json:"isPublic"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// CampaignListOutput 体验活动列表输出
type CampaignListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CampaignNo string `json:"campaignNo"`
	Title string `json:"title"`
	Banner string `json:"banner"`
	Type int `json:"type"`
	Channel int `json:"channel"`
	BudgetAmount int `json:"budgetAmount"`
	LandingURL string `json:"landingURL"`
	RuleJSON string `json:"ruleJSON"`
	IntroContent string `json:"introContent"`
	StartAt *gtime.Time `json:"startAt"`
	EndAt *gtime.Time `json:"endAt"`
	IsPublic int `json:"isPublic"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// CampaignListInput 体验活动列表查询输入
type CampaignListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	CampaignNo string `json:"campaignNo"`
	Title string `json:"title"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Type *int `json:"type"`
	Channel *int `json:"channel"`
	IsPublic *int `json:"isPublic"`
	Status *int `json:"status"`
	StartAtStart string `json:"startAtStart"`
	StartAtEnd string `json:"startAtEnd"`
	EndAtStart string `json:"endAtStart"`
	EndAtEnd string `json:"endAtEnd"`
}

// CampaignBatchUpdateInput 批量编辑体验活动输入
type CampaignBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Type *int `json:"type"`
	Channel *int `json:"channel"`
	IsPublic *int `json:"isPublic"`
	Status *int `json:"status"`
}

