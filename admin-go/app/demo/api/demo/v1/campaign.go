package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Campaign API

// CampaignCreateReq 创建体验活动请求
type CampaignCreateReq struct {
	g.Meta `path:"/campaign/create" method:"post" tags:"体验活动" summary:"创建体验活动"`
	CampaignNo string `json:"campaignNo" v:"required|max-length:50" dc:"活动编号"`
	Title string `json:"title" v:"required|max-length:120" dc:"活动标题"`
	Banner string `json:"banner" v:"max-length:500" dc:"横幅图"`
	Type int `json:"type"  dc:"活动类型"`
	Channel int `json:"channel"  dc:"投放渠道"`
	BudgetAmount int `json:"budgetAmount"  dc:"预算金额（分）"`
	LandingURL string `json:"landingURL" v:"url|max-length:500" dc:"落地页URL"`
	RuleJSON string `json:"ruleJSON" v:"max-length:65535" dc:"规则JSON"`
	IntroContent string `json:"introContent" v:"max-length:65535" dc:"活动介绍"`
	StartAt *gtime.Time `json:"startAt"  dc:"开始时间"`
	EndAt *gtime.Time `json:"endAt"  dc:"结束时间"`
	IsPublic int `json:"isPublic"  dc:"是否公开"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CampaignCreateRes 创建体验活动响应
type CampaignCreateRes struct {
	g.Meta `mime:"application/json"`
}

// CampaignUpdateReq 更新体验活动请求
type CampaignUpdateReq struct {
	g.Meta `path:"/campaign/update" method:"put" tags:"体验活动" summary:"更新体验活动"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验活动ID"`
	CampaignNo string `json:"campaignNo" v:"max-length:50" dc:"活动编号"`
	Title string `json:"title" v:"max-length:120" dc:"活动标题"`
	Banner string `json:"banner" v:"max-length:500" dc:"横幅图"`
	Type int `json:"type"  dc:"活动类型"`
	Channel int `json:"channel"  dc:"投放渠道"`
	BudgetAmount int `json:"budgetAmount"  dc:"预算金额（分）"`
	LandingURL string `json:"landingURL" v:"url|max-length:500" dc:"落地页URL"`
	RuleJSON string `json:"ruleJSON" v:"max-length:65535" dc:"规则JSON"`
	IntroContent string `json:"introContent" v:"max-length:65535" dc:"活动介绍"`
	StartAt *gtime.Time `json:"startAt"  dc:"开始时间"`
	EndAt *gtime.Time `json:"endAt"  dc:"结束时间"`
	IsPublic int `json:"isPublic"  dc:"是否公开"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CampaignUpdateRes 更新体验活动响应
type CampaignUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CampaignDeleteReq 删除体验活动请求
type CampaignDeleteReq struct {
	g.Meta `path:"/campaign/delete" method:"delete" tags:"体验活动" summary:"删除体验活动"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验活动ID"`
}

// CampaignDeleteRes 删除体验活动响应
type CampaignDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CampaignBatchDeleteReq 批量删除体验活动请求
type CampaignBatchDeleteReq struct {
	g.Meta `path:"/campaign/batch-delete" method:"delete" tags:"体验活动" summary:"批量删除体验活动"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验活动ID列表"`
}

// CampaignBatchDeleteRes 批量删除体验活动响应
type CampaignBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CampaignBatchUpdateReq 批量编辑体验活动请求
type CampaignBatchUpdateReq struct {
	g.Meta `path:"/campaign/batch-update" method:"put" tags:"体验活动" summary:"批量编辑体验活动"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验活动ID列表"`
	Type *int `json:"type" dc:"活动类型"`
	Channel *int `json:"channel" dc:"投放渠道"`
	IsPublic *int `json:"isPublic" dc:"是否公开"`
	Status *int `json:"status" dc:"状态"`
}

// CampaignBatchUpdateRes 批量编辑体验活动响应
type CampaignBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CampaignDetailReq 获取体验活动详情请求
type CampaignDetailReq struct {
	g.Meta `path:"/campaign/detail" method:"get" tags:"体验活动" summary:"获取体验活动详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验活动ID"`
}

// CampaignDetailRes 获取体验活动详情响应
type CampaignDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.CampaignDetailOutput
}

// CampaignListReq 获取体验活动列表请求
type CampaignListReq struct {
	g.Meta    `path:"/campaign/list" method:"get" tags:"体验活动" summary:"获取体验活动列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	CampaignNo string `json:"campaignNo" dc:"活动编号"`
	Title string `json:"title" dc:"活动标题"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Type *int `json:"type" dc:"活动类型"`
	Channel *int `json:"channel" dc:"投放渠道"`
	IsPublic *int `json:"isPublic" dc:"是否公开"`
	Status *int `json:"status" dc:"状态"`
	StartAtStart string `json:"startAtStart" dc:"开始时间开始时间"`
	StartAtEnd string `json:"startAtEnd" dc:"开始时间结束时间"`
	EndAtStart string `json:"endAtStart" dc:"结束时间开始时间"`
	EndAtEnd string `json:"endAtEnd" dc:"结束时间结束时间"`
}

// CampaignListRes 获取体验活动列表响应
type CampaignListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.CampaignListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// CampaignExportReq 导出体验活动请求
type CampaignExportReq struct {
	g.Meta    `path:"/campaign/export" method:"get" tags:"体验活动" summary:"导出体验活动"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	CampaignNo string `json:"campaignNo" dc:"活动编号"`
	Title string `json:"title" dc:"活动标题"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Type *int `json:"type" dc:"活动类型"`
	Channel *int `json:"channel" dc:"投放渠道"`
	IsPublic *int `json:"isPublic" dc:"是否公开"`
	Status *int `json:"status" dc:"状态"`
	StartAtStart string `json:"startAtStart" dc:"开始时间开始时间"`
	StartAtEnd string `json:"startAtEnd" dc:"开始时间结束时间"`
	EndAtStart string `json:"endAtStart" dc:"结束时间开始时间"`
	EndAtEnd string `json:"endAtEnd" dc:"结束时间结束时间"`
}

// CampaignExportRes 导出体验活动响应
type CampaignExportRes struct {
	g.Meta `mime:"text/csv"`
}

// CampaignImportReq 导入体验活动请求
type CampaignImportReq struct {
	g.Meta `path:"/campaign/import" method:"post" mime:"multipart/form-data" tags:"体验活动" summary:"导入体验活动"`
}

// CampaignImportRes 导入体验活动响应
type CampaignImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// CampaignImportTemplateReq 下载体验活动导入模板
type CampaignImportTemplateReq struct {
	g.Meta `path:"/campaign/import-template" method:"get" tags:"体验活动" summary:"下载体验活动导入模板"`
}

// CampaignImportTemplateRes 下载体验活动导入模板响应
type CampaignImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

