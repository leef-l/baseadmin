package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Survey DTO 模型

// SurveyCreateInput 创建体验问卷输入
type SurveyCreateInput struct {
	SurveyNo string `json:"surveyNo"`
	Title string `json:"title"`
	Poster string `json:"poster"`
	QuestionJSON string `json:"questionJSON"`
	IntroContent string `json:"introContent"`
	PublishAt *gtime.Time `json:"publishAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	IsAnonymous int `json:"isAnonymous"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// SurveyUpdateInput 更新体验问卷输入
type SurveyUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	SurveyNo string `json:"surveyNo"`
	Title string `json:"title"`
	Poster string `json:"poster"`
	QuestionJSON string `json:"questionJSON"`
	IntroContent string `json:"introContent"`
	PublishAt *gtime.Time `json:"publishAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	IsAnonymous int `json:"isAnonymous"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// SurveyDetailOutput 体验问卷详情输出
type SurveyDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	SurveyNo string `json:"surveyNo"`
	Title string `json:"title"`
	Poster string `json:"poster"`
	QuestionJSON string `json:"questionJSON"`
	IntroContent string `json:"introContent"`
	PublishAt *gtime.Time `json:"publishAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	IsAnonymous int `json:"isAnonymous"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// SurveyListOutput 体验问卷列表输出
type SurveyListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	SurveyNo string `json:"surveyNo"`
	Title string `json:"title"`
	Poster string `json:"poster"`
	QuestionJSON string `json:"questionJSON"`
	IntroContent string `json:"introContent"`
	PublishAt *gtime.Time `json:"publishAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	IsAnonymous int `json:"isAnonymous"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// SurveyListInput 体验问卷列表查询输入
type SurveyListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	SurveyNo string `json:"surveyNo"`
	Title string `json:"title"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	IsAnonymous *int `json:"isAnonymous"`
	Status *int `json:"status"`
	PublishAtStart string `json:"publishAtStart"`
	PublishAtEnd string `json:"publishAtEnd"`
	ExpireAtStart string `json:"expireAtStart"`
	ExpireAtEnd string `json:"expireAtEnd"`
}

// SurveyBatchUpdateInput 批量编辑体验问卷输入
type SurveyBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	IsAnonymous *int `json:"isAnonymous"`
	Status *int `json:"status"`
}

