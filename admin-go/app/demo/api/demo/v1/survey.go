package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Survey API

// SurveyCreateReq 创建体验问卷请求
type SurveyCreateReq struct {
	g.Meta `path:"/survey/create" method:"post" tags:"体验问卷" summary:"创建体验问卷"`
	SurveyNo string `json:"surveyNo" v:"required|max-length:50" dc:"问卷编号"`
	Title string `json:"title" v:"required|max-length:120" dc:"问卷标题"`
	Poster string `json:"poster" v:"max-length:500" dc:"海报"`
	QuestionJSON string `json:"questionJSON" v:"max-length:65535" dc:"问题JSON"`
	IntroContent string `json:"introContent" v:"max-length:65535" dc:"问卷介绍"`
	PublishAt *gtime.Time `json:"publishAt"  dc:"发布时间"`
	ExpireAt *gtime.Time `json:"expireAt"  dc:"过期时间"`
	IsAnonymous int `json:"isAnonymous"  dc:"是否匿名"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// SurveyCreateRes 创建体验问卷响应
type SurveyCreateRes struct {
	g.Meta `mime:"application/json"`
}

// SurveyUpdateReq 更新体验问卷请求
type SurveyUpdateReq struct {
	g.Meta `path:"/survey/update" method:"put" tags:"体验问卷" summary:"更新体验问卷"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验问卷ID"`
	SurveyNo string `json:"surveyNo" v:"max-length:50" dc:"问卷编号"`
	Title string `json:"title" v:"max-length:120" dc:"问卷标题"`
	Poster string `json:"poster" v:"max-length:500" dc:"海报"`
	QuestionJSON string `json:"questionJSON" v:"max-length:65535" dc:"问题JSON"`
	IntroContent string `json:"introContent" v:"max-length:65535" dc:"问卷介绍"`
	PublishAt *gtime.Time `json:"publishAt"  dc:"发布时间"`
	ExpireAt *gtime.Time `json:"expireAt"  dc:"过期时间"`
	IsAnonymous int `json:"isAnonymous"  dc:"是否匿名"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// SurveyUpdateRes 更新体验问卷响应
type SurveyUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// SurveyDeleteReq 删除体验问卷请求
type SurveyDeleteReq struct {
	g.Meta `path:"/survey/delete" method:"delete" tags:"体验问卷" summary:"删除体验问卷"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验问卷ID"`
}

// SurveyDeleteRes 删除体验问卷响应
type SurveyDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// SurveyBatchDeleteReq 批量删除体验问卷请求
type SurveyBatchDeleteReq struct {
	g.Meta `path:"/survey/batch-delete" method:"delete" tags:"体验问卷" summary:"批量删除体验问卷"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验问卷ID列表"`
}

// SurveyBatchDeleteRes 批量删除体验问卷响应
type SurveyBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// SurveyBatchUpdateReq 批量编辑体验问卷请求
type SurveyBatchUpdateReq struct {
	g.Meta `path:"/survey/batch-update" method:"put" tags:"体验问卷" summary:"批量编辑体验问卷"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验问卷ID列表"`
	IsAnonymous *int `json:"isAnonymous" dc:"是否匿名"`
	Status *int `json:"status" dc:"状态"`
}

// SurveyBatchUpdateRes 批量编辑体验问卷响应
type SurveyBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// SurveyDetailReq 获取体验问卷详情请求
type SurveyDetailReq struct {
	g.Meta `path:"/survey/detail" method:"get" tags:"体验问卷" summary:"获取体验问卷详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验问卷ID"`
}

// SurveyDetailRes 获取体验问卷详情响应
type SurveyDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.SurveyDetailOutput
}

// SurveyListReq 获取体验问卷列表请求
type SurveyListReq struct {
	g.Meta    `path:"/survey/list" method:"get" tags:"体验问卷" summary:"获取体验问卷列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	SurveyNo string `json:"surveyNo" dc:"问卷编号"`
	Title string `json:"title" dc:"问卷标题"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsAnonymous *int `json:"isAnonymous" dc:"是否匿名"`
	Status *int `json:"status" dc:"状态"`
	PublishAtStart string `json:"publishAtStart" dc:"发布时间开始时间"`
	PublishAtEnd string `json:"publishAtEnd" dc:"发布时间结束时间"`
	ExpireAtStart string `json:"expireAtStart" dc:"过期时间开始时间"`
	ExpireAtEnd string `json:"expireAtEnd" dc:"过期时间结束时间"`
}

// SurveyListRes 获取体验问卷列表响应
type SurveyListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.SurveyListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// SurveyExportReq 导出体验问卷请求
type SurveyExportReq struct {
	g.Meta    `path:"/survey/export" method:"get" tags:"体验问卷" summary:"导出体验问卷"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	SurveyNo string `json:"surveyNo" dc:"问卷编号"`
	Title string `json:"title" dc:"问卷标题"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsAnonymous *int `json:"isAnonymous" dc:"是否匿名"`
	Status *int `json:"status" dc:"状态"`
	PublishAtStart string `json:"publishAtStart" dc:"发布时间开始时间"`
	PublishAtEnd string `json:"publishAtEnd" dc:"发布时间结束时间"`
	ExpireAtStart string `json:"expireAtStart" dc:"过期时间开始时间"`
	ExpireAtEnd string `json:"expireAtEnd" dc:"过期时间结束时间"`
}

// SurveyExportRes 导出体验问卷响应
type SurveyExportRes struct {
	g.Meta `mime:"text/csv"`
}

// SurveyImportReq 导入体验问卷请求
type SurveyImportReq struct {
	g.Meta `path:"/survey/import" method:"post" mime:"multipart/form-data" tags:"体验问卷" summary:"导入体验问卷"`
}

// SurveyImportRes 导入体验问卷响应
type SurveyImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// SurveyImportTemplateReq 下载体验问卷导入模板
type SurveyImportTemplateReq struct {
	g.Meta `path:"/survey/import-template" method:"get" tags:"体验问卷" summary:"下载体验问卷导入模板"`
}

// SurveyImportTemplateRes 下载体验问卷导入模板响应
type SurveyImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

