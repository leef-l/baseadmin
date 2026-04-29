package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Appointment API

// AppointmentCreateReq 创建体验预约请求
type AppointmentCreateReq struct {
	g.Meta `path:"/appointment/create" method:"post" tags:"体验预约" summary:"创建体验预约"`
	AppointmentNo string `json:"appointmentNo" v:"required|max-length:50" dc:"预约编号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	Subject string `json:"subject" v:"required|max-length:120" dc:"预约主题"`
	AppointmentAt *gtime.Time `json:"appointmentAt"  dc:"预约时间"`
	ContactPhone string `json:"contactPhone" v:"phone-loose|max-length:30" dc:"联系电话"`
	Address string `json:"address" v:"max-length:255" dc:"预约地址"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// AppointmentCreateRes 创建体验预约响应
type AppointmentCreateRes struct {
	g.Meta `mime:"application/json"`
}

// AppointmentUpdateReq 更新体验预约请求
type AppointmentUpdateReq struct {
	g.Meta `path:"/appointment/update" method:"put" tags:"体验预约" summary:"更新体验预约"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验预约ID"`
	AppointmentNo string `json:"appointmentNo" v:"max-length:50" dc:"预约编号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	Subject string `json:"subject" v:"max-length:120" dc:"预约主题"`
	AppointmentAt *gtime.Time `json:"appointmentAt"  dc:"预约时间"`
	ContactPhone string `json:"contactPhone" v:"phone-loose|max-length:30" dc:"联系电话"`
	Address string `json:"address" v:"max-length:255" dc:"预约地址"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// AppointmentUpdateRes 更新体验预约响应
type AppointmentUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// AppointmentDeleteReq 删除体验预约请求
type AppointmentDeleteReq struct {
	g.Meta `path:"/appointment/delete" method:"delete" tags:"体验预约" summary:"删除体验预约"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验预约ID"`
}

// AppointmentDeleteRes 删除体验预约响应
type AppointmentDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// AppointmentBatchDeleteReq 批量删除体验预约请求
type AppointmentBatchDeleteReq struct {
	g.Meta `path:"/appointment/batch-delete" method:"delete" tags:"体验预约" summary:"批量删除体验预约"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验预约ID列表"`
}

// AppointmentBatchDeleteRes 批量删除体验预约响应
type AppointmentBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// AppointmentBatchUpdateReq 批量编辑体验预约请求
type AppointmentBatchUpdateReq struct {
	g.Meta `path:"/appointment/batch-update" method:"put" tags:"体验预约" summary:"批量编辑体验预约"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验预约ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// AppointmentBatchUpdateRes 批量编辑体验预约响应
type AppointmentBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// AppointmentDetailReq 获取体验预约详情请求
type AppointmentDetailReq struct {
	g.Meta `path:"/appointment/detail" method:"get" tags:"体验预约" summary:"获取体验预约详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验预约ID"`
}

// AppointmentDetailRes 获取体验预约详情响应
type AppointmentDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.AppointmentDetailOutput
}

// AppointmentListReq 获取体验预约列表请求
type AppointmentListReq struct {
	g.Meta    `path:"/appointment/list" method:"get" tags:"体验预约" summary:"获取体验预约列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	AppointmentNo string `json:"appointmentNo" dc:"预约编号"`
	Subject string `json:"subject" dc:"预约主题"`
	ContactPhone string `json:"contactPhone" dc:"联系电话"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	AppointmentAtStart string `json:"appointmentAtStart" dc:"预约时间开始时间"`
	AppointmentAtEnd string `json:"appointmentAtEnd" dc:"预约时间结束时间"`
}

// AppointmentListRes 获取体验预约列表响应
type AppointmentListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.AppointmentListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// AppointmentExportReq 导出体验预约请求
type AppointmentExportReq struct {
	g.Meta    `path:"/appointment/export" method:"get" tags:"体验预约" summary:"导出体验预约"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	AppointmentNo string `json:"appointmentNo" dc:"预约编号"`
	Subject string `json:"subject" dc:"预约主题"`
	ContactPhone string `json:"contactPhone" dc:"联系电话"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	AppointmentAtStart string `json:"appointmentAtStart" dc:"预约时间开始时间"`
	AppointmentAtEnd string `json:"appointmentAtEnd" dc:"预约时间结束时间"`
}

// AppointmentExportRes 导出体验预约响应
type AppointmentExportRes struct {
	g.Meta `mime:"text/csv"`
}

// AppointmentImportReq 导入体验预约请求
type AppointmentImportReq struct {
	g.Meta `path:"/appointment/import" method:"post" mime:"multipart/form-data" tags:"体验预约" summary:"导入体验预约"`
}

// AppointmentImportRes 导入体验预约响应
type AppointmentImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// AppointmentImportTemplateReq 下载体验预约导入模板
type AppointmentImportTemplateReq struct {
	g.Meta `path:"/appointment/import-template" method:"get" tags:"体验预约" summary:"下载体验预约导入模板"`
}

// AppointmentImportTemplateRes 下载体验预约导入模板响应
type AppointmentImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

