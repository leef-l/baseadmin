package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Level API

// LevelCreateReq 创建会员等级配置请求
type LevelCreateReq struct {
	g.Meta `path:"/level/create" method:"post" tags:"会员等级配置" summary:"创建会员等级配置"`
	Name string `json:"name" v:"required|max-length:50" dc:"等级名称"`
	LevelNo int `json:"levelNo"  dc:"等级编号（越大越高）"`
	Icon string `json:"icon" v:"max-length:500" dc:"等级图标"`
	DurationDays int `json:"durationDays"  dc:"有效天数（0=永久）"`
	NeedActiveCount int `json:"needActiveCount"  dc:"升级所需有效用户数"`
	NeedTeamTurnover int64 `json:"needTeamTurnover"  dc:"升级所需团队营业额（分）"`
	IsTop int `json:"isTop"  dc:"是否最高等级"`
	AutoDeploy int `json:"autoDeploy"  dc:"到达后自动部署站点"`
	Remark string `json:"remark" v:"max-length:500" dc:"等级说明"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// LevelCreateRes 创建会员等级配置响应
type LevelCreateRes struct {
	g.Meta `mime:"application/json"`
}

// LevelUpdateReq 更新会员等级配置请求
type LevelUpdateReq struct {
	g.Meta `path:"/level/update" method:"put" tags:"会员等级配置" summary:"更新会员等级配置"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员等级配置ID"`
	Name string `json:"name" v:"max-length:50" dc:"等级名称"`
	LevelNo int `json:"levelNo"  dc:"等级编号（越大越高）"`
	Icon string `json:"icon" v:"max-length:500" dc:"等级图标"`
	DurationDays int `json:"durationDays"  dc:"有效天数（0=永久）"`
	NeedActiveCount int `json:"needActiveCount"  dc:"升级所需有效用户数"`
	NeedTeamTurnover int64 `json:"needTeamTurnover"  dc:"升级所需团队营业额（分）"`
	IsTop int `json:"isTop"  dc:"是否最高等级"`
	AutoDeploy int `json:"autoDeploy"  dc:"到达后自动部署站点"`
	Remark string `json:"remark" v:"max-length:500" dc:"等级说明"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// LevelUpdateRes 更新会员等级配置响应
type LevelUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// LevelDeleteReq 删除会员等级配置请求
type LevelDeleteReq struct {
	g.Meta `path:"/level/delete" method:"delete" tags:"会员等级配置" summary:"删除会员等级配置"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员等级配置ID"`
}

// LevelDeleteRes 删除会员等级配置响应
type LevelDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// LevelBatchDeleteReq 批量删除会员等级配置请求
type LevelBatchDeleteReq struct {
	g.Meta `path:"/level/batch-delete" method:"delete" tags:"会员等级配置" summary:"批量删除会员等级配置"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员等级配置ID列表"`
}

// LevelBatchDeleteRes 批量删除会员等级配置响应
type LevelBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// LevelBatchUpdateReq 批量编辑会员等级配置请求
type LevelBatchUpdateReq struct {
	g.Meta `path:"/level/batch-update" method:"put" tags:"会员等级配置" summary:"批量编辑会员等级配置"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员等级配置ID列表"`
	IsTop *int `json:"isTop" dc:"是否最高等级"`
	AutoDeploy *int `json:"autoDeploy" dc:"到达后自动部署站点"`
	Status *int `json:"status" dc:"状态"`
}

// LevelBatchUpdateRes 批量编辑会员等级配置响应
type LevelBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// LevelDetailReq 获取会员等级配置详情请求
type LevelDetailReq struct {
	g.Meta `path:"/level/detail" method:"get" tags:"会员等级配置" summary:"获取会员等级配置详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员等级配置ID"`
}

// LevelDetailRes 获取会员等级配置详情响应
type LevelDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.LevelDetailOutput
}

// LevelListReq 获取会员等级配置列表请求
type LevelListReq struct {
	g.Meta    `path:"/level/list" method:"get" tags:"会员等级配置" summary:"获取会员等级配置列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"等级名称"`
	LevelNo string `json:"levelNo" dc:"等级编号（越大越高）"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsTop *int `json:"isTop" dc:"是否最高等级"`
	AutoDeploy *int `json:"autoDeploy" dc:"到达后自动部署站点"`
	Status *int `json:"status" dc:"状态"`
}

// LevelListRes 获取会员等级配置列表响应
type LevelListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.LevelListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// LevelExportReq 导出会员等级配置请求
type LevelExportReq struct {
	g.Meta    `path:"/level/export" method:"get" tags:"会员等级配置" summary:"导出会员等级配置"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"等级名称"`
	LevelNo string `json:"levelNo" dc:"等级编号（越大越高）"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsTop *int `json:"isTop" dc:"是否最高等级"`
	AutoDeploy *int `json:"autoDeploy" dc:"到达后自动部署站点"`
	Status *int `json:"status" dc:"状态"`
}

// LevelExportRes 导出会员等级配置响应
type LevelExportRes struct {
	g.Meta `mime:"text/csv"`
}

// LevelImportReq 导入会员等级配置请求
type LevelImportReq struct {
	g.Meta `path:"/level/import" method:"post" mime:"multipart/form-data" tags:"会员等级配置" summary:"导入会员等级配置"`
}

// LevelImportRes 导入会员等级配置响应
type LevelImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// LevelImportTemplateReq 下载会员等级配置导入模板
type LevelImportTemplateReq struct {
	g.Meta `path:"/level/import-template" method:"get" tags:"会员等级配置" summary:"下载会员等级配置导入模板"`
}

// LevelImportTemplateRes 下载会员等级配置导入模板响应
type LevelImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

