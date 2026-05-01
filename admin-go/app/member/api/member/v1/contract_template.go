package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// ContractTemplate API

// ContractTemplateCreateReq 创建会员合同模板请求
type ContractTemplateCreateReq struct {
	g.Meta `path:"/contract_template/create" method:"post" tags:"会员合同模板" summary:"创建会员合同模板"`
	TemplateName string `json:"templateName" v:"required|max-length:64" dc:"模板名称"`
	TemplateType string `json:"templateType" v:"max-length:32" dc:"模板类型"`
	Content string `json:"content" v:"required|max-length:16777215" dc:"模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）"`
	IsDefault int `json:"isDefault"  dc:"是否默认模板"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ContractTemplateCreateRes 创建会员合同模板响应
type ContractTemplateCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ContractTemplateUpdateReq 更新会员合同模板请求
type ContractTemplateUpdateReq struct {
	g.Meta `path:"/contract_template/update" method:"put" tags:"会员合同模板" summary:"更新会员合同模板"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员合同模板ID"`
	TemplateName string `json:"templateName" v:"max-length:64" dc:"模板名称"`
	TemplateType string `json:"templateType" v:"max-length:32" dc:"模板类型"`
	Content string `json:"content" v:"max-length:16777215" dc:"模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）"`
	IsDefault int `json:"isDefault"  dc:"是否默认模板"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ContractTemplateUpdateRes 更新会员合同模板响应
type ContractTemplateUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ContractTemplateDeleteReq 删除会员合同模板请求
type ContractTemplateDeleteReq struct {
	g.Meta `path:"/contract_template/delete" method:"delete" tags:"会员合同模板" summary:"删除会员合同模板"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员合同模板ID"`
}

// ContractTemplateDeleteRes 删除会员合同模板响应
type ContractTemplateDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ContractTemplateBatchDeleteReq 批量删除会员合同模板请求
type ContractTemplateBatchDeleteReq struct {
	g.Meta `path:"/contract_template/batch-delete" method:"delete" tags:"会员合同模板" summary:"批量删除会员合同模板"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员合同模板ID列表"`
}

// ContractTemplateBatchDeleteRes 批量删除会员合同模板响应
type ContractTemplateBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ContractTemplateBatchUpdateReq 批量编辑会员合同模板请求
type ContractTemplateBatchUpdateReq struct {
	g.Meta `path:"/contract_template/batch-update" method:"put" tags:"会员合同模板" summary:"批量编辑会员合同模板"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员合同模板ID列表"`
	IsDefault *int `json:"isDefault" dc:"是否默认模板"`
	Status *int `json:"status" dc:"状态"`
}

// ContractTemplateBatchUpdateRes 批量编辑会员合同模板响应
type ContractTemplateBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ContractTemplateDetailReq 获取会员合同模板详情请求
type ContractTemplateDetailReq struct {
	g.Meta `path:"/contract_template/detail" method:"get" tags:"会员合同模板" summary:"获取会员合同模板详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员合同模板ID"`
}

// ContractTemplateDetailRes 获取会员合同模板详情响应
type ContractTemplateDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ContractTemplateDetailOutput
}

// ContractTemplateListReq 获取会员合同模板列表请求
type ContractTemplateListReq struct {
	g.Meta    `path:"/contract_template/list" method:"get" tags:"会员合同模板" summary:"获取会员合同模板列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	TemplateName string `json:"templateName" dc:"模板名称"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsDefault *int `json:"isDefault" dc:"是否默认模板"`
	Status *int `json:"status" dc:"状态"`
	TemplateType *string `json:"templateType" dc:"模板类型"`
}

// ContractTemplateListRes 获取会员合同模板列表响应
type ContractTemplateListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ContractTemplateListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// ContractTemplateExportReq 导出会员合同模板请求
type ContractTemplateExportReq struct {
	g.Meta    `path:"/contract_template/export" method:"get" tags:"会员合同模板" summary:"导出会员合同模板"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	TemplateName string `json:"templateName" dc:"模板名称"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsDefault *int `json:"isDefault" dc:"是否默认模板"`
	Status *int `json:"status" dc:"状态"`
	TemplateType *string `json:"templateType" dc:"模板类型"`
}

// ContractTemplateExportRes 导出会员合同模板响应
type ContractTemplateExportRes struct {
	g.Meta `mime:"text/csv"`
}

// ContractTemplateImportReq 导入会员合同模板请求
type ContractTemplateImportReq struct {
	g.Meta `path:"/contract_template/import" method:"post" mime:"multipart/form-data" tags:"会员合同模板" summary:"导入会员合同模板"`
}

// ContractTemplateImportRes 导入会员合同模板响应
type ContractTemplateImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// ContractTemplateImportTemplateReq 下载会员合同模板导入模板
type ContractTemplateImportTemplateReq struct {
	g.Meta `path:"/contract_template/import-template" method:"get" tags:"会员合同模板" summary:"下载会员合同模板导入模板"`
}

// ContractTemplateImportTemplateRes 下载会员合同模板导入模板响应
type ContractTemplateImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

