package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// ContractTemplate DTO 模型

// ContractTemplateCreateInput 创建会员合同模板输入
type ContractTemplateCreateInput struct {
	TemplateName string `json:"templateName"`
	TemplateType string `json:"templateType"`
	Content string `json:"content"`
	IsDefault int `json:"isDefault"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ContractTemplateUpdateInput 更新会员合同模板输入
type ContractTemplateUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TemplateName string `json:"templateName"`
	TemplateType string `json:"templateType"`
	Content string `json:"content"`
	IsDefault int `json:"isDefault"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ContractTemplateDetailOutput 会员合同模板详情输出
type ContractTemplateDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TemplateName string `json:"templateName"`
	TemplateType string `json:"templateType"`
	Content string `json:"content"`
	IsDefault int `json:"isDefault"`
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

// ContractTemplateListOutput 会员合同模板列表输出
type ContractTemplateListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TemplateName string `json:"templateName"`
	TemplateType string `json:"templateType"`
	Content string `json:"content"`
	IsDefault int `json:"isDefault"`
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

// ContractTemplateListInput 会员合同模板列表查询输入
type ContractTemplateListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	TemplateName string `json:"templateName"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	IsDefault *int `json:"isDefault"`
	Status *int `json:"status"`
	TemplateType *string `json:"templateType"`
}

// ContractTemplateBatchUpdateInput 批量编辑会员合同模板输入
type ContractTemplateBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	IsDefault *int `json:"isDefault"`
	Status *int `json:"status"`
}

