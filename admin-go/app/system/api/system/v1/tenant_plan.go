package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// TenantPlan API

// TenantPlanCreateReq 创建租户套餐订阅表请求
type TenantPlanCreateReq struct {
	g.Meta `path:"/tenant_plan/create" method:"post" tags:"租户套餐订阅表" summary:"创建租户套餐订阅表"`
	TenantID snowflake.JsonInt64 `json:"tenantID" v:"required" dc:"租户"`
	PlanID snowflake.JsonInt64 `json:"planID" v:"required" dc:"套餐ID"`
	StartAt *gtime.Time `json:"startAt"  dc:"生效时间"`
	ExpireAt *gtime.Time `json:"expireAt"  dc:"到期时间"`
	Status int `json:"status"  dc:"状态"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// TenantPlanCreateRes 创建租户套餐订阅表响应
type TenantPlanCreateRes struct {
	g.Meta `mime:"application/json"`
}

// TenantPlanUpdateReq 更新租户套餐订阅表请求
type TenantPlanUpdateReq struct {
	g.Meta `path:"/tenant_plan/update" method:"put" tags:"租户套餐订阅表" summary:"更新租户套餐订阅表"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"租户套餐订阅表ID"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	PlanID snowflake.JsonInt64 `json:"planID"  dc:"套餐ID"`
	StartAt *gtime.Time `json:"startAt"  dc:"生效时间"`
	ExpireAt *gtime.Time `json:"expireAt"  dc:"到期时间"`
	Status int `json:"status"  dc:"状态"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// TenantPlanUpdateRes 更新租户套餐订阅表响应
type TenantPlanUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// TenantPlanDeleteReq 删除租户套餐订阅表请求
type TenantPlanDeleteReq struct {
	g.Meta `path:"/tenant_plan/delete" method:"delete" tags:"租户套餐订阅表" summary:"删除租户套餐订阅表"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"租户套餐订阅表ID"`
}

// TenantPlanDeleteRes 删除租户套餐订阅表响应
type TenantPlanDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// TenantPlanBatchDeleteReq 批量删除租户套餐订阅表请求
type TenantPlanBatchDeleteReq struct {
	g.Meta `path:"/tenant_plan/batch-delete" method:"delete" tags:"租户套餐订阅表" summary:"批量删除租户套餐订阅表"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"租户套餐订阅表ID列表"`
}

// TenantPlanBatchDeleteRes 批量删除租户套餐订阅表响应
type TenantPlanBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// TenantPlanBatchUpdateReq 批量编辑租户套餐订阅表请求
type TenantPlanBatchUpdateReq struct {
	g.Meta `path:"/tenant_plan/batch-update" method:"put" tags:"租户套餐订阅表" summary:"批量编辑租户套餐订阅表"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"租户套餐订阅表ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// TenantPlanBatchUpdateRes 批量编辑租户套餐订阅表响应
type TenantPlanBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// TenantPlanDetailReq 获取租户套餐订阅表详情请求
type TenantPlanDetailReq struct {
	g.Meta `path:"/tenant_plan/detail" method:"get" tags:"租户套餐订阅表" summary:"获取租户套餐订阅表详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"租户套餐订阅表ID"`
}

// TenantPlanDetailRes 获取租户套餐订阅表详情响应
type TenantPlanDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.TenantPlanDetailOutput
}

// TenantPlanListReq 获取租户套餐订阅表列表请求
type TenantPlanListReq struct {
	g.Meta    `path:"/tenant_plan/list" method:"get" tags:"租户套餐订阅表" summary:"获取租户套餐订阅表列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	PlanID *snowflake.JsonInt64 `json:"planID" dc:"套餐ID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	Remark string `json:"remark" dc:"备注"`
	StartAtStart string `json:"startAtStart" dc:"生效时间开始时间"`
	StartAtEnd string `json:"startAtEnd" dc:"生效时间结束时间"`
	ExpireAtStart string `json:"expireAtStart" dc:"到期时间开始时间"`
	ExpireAtEnd string `json:"expireAtEnd" dc:"到期时间结束时间"`
}

// TenantPlanListRes 获取租户套餐订阅表列表响应
type TenantPlanListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.TenantPlanListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// TenantPlanExportReq 导出租户套餐订阅表请求
type TenantPlanExportReq struct {
	g.Meta    `path:"/tenant_plan/export" method:"get" tags:"租户套餐订阅表" summary:"导出租户套餐订阅表"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	PlanID *snowflake.JsonInt64 `json:"planID" dc:"套餐ID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	Remark string `json:"remark" dc:"备注"`
	StartAtStart string `json:"startAtStart" dc:"生效时间开始时间"`
	StartAtEnd string `json:"startAtEnd" dc:"生效时间结束时间"`
	ExpireAtStart string `json:"expireAtStart" dc:"到期时间开始时间"`
	ExpireAtEnd string `json:"expireAtEnd" dc:"到期时间结束时间"`
}

// TenantPlanExportRes 导出租户套餐订阅表响应
type TenantPlanExportRes struct {
	g.Meta `mime:"text/csv"`
}

// TenantPlanImportReq 导入租户套餐订阅表请求
type TenantPlanImportReq struct {
	g.Meta `path:"/tenant_plan/import" method:"post" mime:"multipart/form-data" tags:"租户套餐订阅表" summary:"导入租户套餐订阅表"`
}

// TenantPlanImportRes 导入租户套餐订阅表响应
type TenantPlanImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// TenantPlanImportTemplateReq 下载租户套餐订阅表导入模板
type TenantPlanImportTemplateReq struct {
	g.Meta `path:"/tenant_plan/import-template" method:"get" tags:"租户套餐订阅表" summary:"下载租户套餐订阅表导入模板"`
}

// TenantPlanImportTemplateRes 下载租户套餐订阅表导入模板响应
type TenantPlanImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

