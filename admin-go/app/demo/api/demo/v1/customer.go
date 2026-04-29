package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Customer API

// CustomerCreateReq 创建体验客户请求
type CustomerCreateReq struct {
	g.Meta `path:"/customer/create" method:"post" tags:"体验客户" summary:"创建体验客户"`
	Avatar string `json:"avatar" v:"max-length:500" dc:"头像"`
	Name string `json:"name" v:"required|max-length:80" dc:"客户名称"`
	CustomerNo string `json:"customerNo" v:"required|max-length:50" dc:"客户编号"`
	Phone string `json:"phone" v:"phone-loose|max-length:30" dc:"联系电话"`
	Email string `json:"email" v:"email|max-length:120" dc:"邮箱"`
	Gender int `json:"gender"  dc:"性别"`
	Level int `json:"level"  dc:"等级"`
	SourceType int `json:"sourceType"  dc:"来源"`
	IsVip int `json:"isVip"  dc:"是否VIP"`
	RegisteredAt *gtime.Time `json:"registeredAt"  dc:"注册时间"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CustomerCreateRes 创建体验客户响应
type CustomerCreateRes struct {
	g.Meta `mime:"application/json"`
}

// CustomerUpdateReq 更新体验客户请求
type CustomerUpdateReq struct {
	g.Meta `path:"/customer/update" method:"put" tags:"体验客户" summary:"更新体验客户"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验客户ID"`
	Avatar string `json:"avatar" v:"max-length:500" dc:"头像"`
	Name string `json:"name" v:"max-length:80" dc:"客户名称"`
	CustomerNo string `json:"customerNo" v:"max-length:50" dc:"客户编号"`
	Phone string `json:"phone" v:"phone-loose|max-length:30" dc:"联系电话"`
	Email string `json:"email" v:"email|max-length:120" dc:"邮箱"`
	Gender int `json:"gender"  dc:"性别"`
	Level int `json:"level"  dc:"等级"`
	SourceType int `json:"sourceType"  dc:"来源"`
	IsVip int `json:"isVip"  dc:"是否VIP"`
	RegisteredAt *gtime.Time `json:"registeredAt"  dc:"注册时间"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CustomerUpdateRes 更新体验客户响应
type CustomerUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CustomerDeleteReq 删除体验客户请求
type CustomerDeleteReq struct {
	g.Meta `path:"/customer/delete" method:"delete" tags:"体验客户" summary:"删除体验客户"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验客户ID"`
}

// CustomerDeleteRes 删除体验客户响应
type CustomerDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CustomerBatchDeleteReq 批量删除体验客户请求
type CustomerBatchDeleteReq struct {
	g.Meta `path:"/customer/batch-delete" method:"delete" tags:"体验客户" summary:"批量删除体验客户"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验客户ID列表"`
}

// CustomerBatchDeleteRes 批量删除体验客户响应
type CustomerBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CustomerBatchUpdateReq 批量编辑体验客户请求
type CustomerBatchUpdateReq struct {
	g.Meta `path:"/customer/batch-update" method:"put" tags:"体验客户" summary:"批量编辑体验客户"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验客户ID列表"`
	Gender *int `json:"gender" dc:"性别"`
	Level *int `json:"level" dc:"等级"`
	SourceType *int `json:"sourceType" dc:"来源"`
	IsVip *int `json:"isVip" dc:"是否VIP"`
	Status *int `json:"status" dc:"状态"`
}

// CustomerBatchUpdateRes 批量编辑体验客户响应
type CustomerBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CustomerDetailReq 获取体验客户详情请求
type CustomerDetailReq struct {
	g.Meta `path:"/customer/detail" method:"get" tags:"体验客户" summary:"获取体验客户详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验客户ID"`
}

// CustomerDetailRes 获取体验客户详情响应
type CustomerDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.CustomerDetailOutput
}

// CustomerListReq 获取体验客户列表请求
type CustomerListReq struct {
	g.Meta    `path:"/customer/list" method:"get" tags:"体验客户" summary:"获取体验客户列表"`
	PageNum   int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	CustomerNo string `json:"customerNo" dc:"客户编号"`
	Name string `json:"name" dc:"客户名称"`
	Phone string `json:"phone" dc:"联系电话"`
	Email string `json:"email" dc:"邮箱"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Gender *int `json:"gender" dc:"性别"`
	Level *int `json:"level" dc:"等级"`
	SourceType *int `json:"sourceType" dc:"来源"`
	IsVip *int `json:"isVip" dc:"是否VIP"`
	Status *int `json:"status" dc:"状态"`
	RegisteredAtStart string `json:"registeredAtStart" dc:"注册时间开始时间"`
	RegisteredAtEnd string `json:"registeredAtEnd" dc:"注册时间结束时间"`
}

// CustomerListRes 获取体验客户列表响应
type CustomerListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.CustomerListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// CustomerExportReq 导出体验客户请求
type CustomerExportReq struct {
	g.Meta    `path:"/customer/export" method:"get" tags:"体验客户" summary:"导出体验客户"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	CustomerNo string `json:"customerNo" dc:"客户编号"`
	Name string `json:"name" dc:"客户名称"`
	Phone string `json:"phone" dc:"联系电话"`
	Email string `json:"email" dc:"邮箱"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Gender *int `json:"gender" dc:"性别"`
	Level *int `json:"level" dc:"等级"`
	SourceType *int `json:"sourceType" dc:"来源"`
	IsVip *int `json:"isVip" dc:"是否VIP"`
	Status *int `json:"status" dc:"状态"`
	RegisteredAtStart string `json:"registeredAtStart" dc:"注册时间开始时间"`
	RegisteredAtEnd string `json:"registeredAtEnd" dc:"注册时间结束时间"`
}

// CustomerExportRes 导出体验客户响应
type CustomerExportRes struct {
	g.Meta `mime:"text/csv"`
}

// CustomerImportReq 导入体验客户请求
type CustomerImportReq struct {
	g.Meta `path:"/customer/import" method:"post" mime:"multipart/form-data" tags:"体验客户" summary:"导入体验客户"`
}

// CustomerImportRes 导入体验客户响应
type CustomerImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// CustomerImportTemplateReq 下载体验客户导入模板
type CustomerImportTemplateReq struct {
	g.Meta `path:"/customer/import-template" method:"get" tags:"体验客户" summary:"下载体验客户导入模板"`
}

// CustomerImportTemplateRes 下载体验客户导入模板响应
type CustomerImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

