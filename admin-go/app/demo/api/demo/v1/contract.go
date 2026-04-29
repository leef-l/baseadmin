package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Contract API

// ContractCreateReq 创建体验合同请求
type ContractCreateReq struct {
	g.Meta `path:"/contract/create" method:"post" tags:"体验合同" summary:"创建体验合同"`
	ContractNo string `json:"contractNo" v:"required|max-length:50" dc:"合同编号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	OrderID snowflake.JsonInt64 `json:"orderID"  dc:"订单"`
	Title string `json:"title" v:"required|max-length:120" dc:"合同标题"`
	ContractFile string `json:"contractFile" v:"max-length:500" dc:"合同文件"`
	SignImage string `json:"signImage" v:"max-length:500" dc:"签章图片"`
	ContractAmount int `json:"contractAmount"  dc:"合同金额（分）"`
	SignPassword string `json:"signPassword" v:"length:6,32" dc:"签署密码"`
	SignedAt *gtime.Time `json:"signedAt"  dc:"签署时间"`
	ExpiresAt *gtime.Time `json:"expiresAt"  dc:"到期时间"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ContractCreateRes 创建体验合同响应
type ContractCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ContractUpdateReq 更新体验合同请求
type ContractUpdateReq struct {
	g.Meta `path:"/contract/update" method:"put" tags:"体验合同" summary:"更新体验合同"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验合同ID"`
	ContractNo string `json:"contractNo" v:"max-length:50" dc:"合同编号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	OrderID snowflake.JsonInt64 `json:"orderID"  dc:"订单"`
	Title string `json:"title" v:"max-length:120" dc:"合同标题"`
	ContractFile string `json:"contractFile" v:"max-length:500" dc:"合同文件"`
	SignImage string `json:"signImage" v:"max-length:500" dc:"签章图片"`
	ContractAmount int `json:"contractAmount"  dc:"合同金额（分）"`
	SignPassword string `json:"signPassword" v:"length:6,32" dc:"签署密码"`
	SignedAt *gtime.Time `json:"signedAt"  dc:"签署时间"`
	ExpiresAt *gtime.Time `json:"expiresAt"  dc:"到期时间"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ContractUpdateRes 更新体验合同响应
type ContractUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ContractDeleteReq 删除体验合同请求
type ContractDeleteReq struct {
	g.Meta `path:"/contract/delete" method:"delete" tags:"体验合同" summary:"删除体验合同"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验合同ID"`
}

// ContractDeleteRes 删除体验合同响应
type ContractDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ContractBatchDeleteReq 批量删除体验合同请求
type ContractBatchDeleteReq struct {
	g.Meta `path:"/contract/batch-delete" method:"delete" tags:"体验合同" summary:"批量删除体验合同"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验合同ID列表"`
}

// ContractBatchDeleteRes 批量删除体验合同响应
type ContractBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ContractBatchUpdateReq 批量编辑体验合同请求
type ContractBatchUpdateReq struct {
	g.Meta `path:"/contract/batch-update" method:"put" tags:"体验合同" summary:"批量编辑体验合同"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验合同ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// ContractBatchUpdateRes 批量编辑体验合同响应
type ContractBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ContractDetailReq 获取体验合同详情请求
type ContractDetailReq struct {
	g.Meta `path:"/contract/detail" method:"get" tags:"体验合同" summary:"获取体验合同详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验合同ID"`
}

// ContractDetailRes 获取体验合同详情响应
type ContractDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ContractDetailOutput
}

// ContractListReq 获取体验合同列表请求
type ContractListReq struct {
	g.Meta    `path:"/contract/list" method:"get" tags:"体验合同" summary:"获取体验合同列表"`
	PageNum   int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	ContractNo string `json:"contractNo" dc:"合同编号"`
	Title string `json:"title" dc:"合同标题"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	OrderID *snowflake.JsonInt64 `json:"orderID" dc:"订单"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	SignedAtStart string `json:"signedAtStart" dc:"签署时间开始时间"`
	SignedAtEnd string `json:"signedAtEnd" dc:"签署时间结束时间"`
	ExpiresAtStart string `json:"expiresAtStart" dc:"到期时间开始时间"`
	ExpiresAtEnd string `json:"expiresAtEnd" dc:"到期时间结束时间"`
}

// ContractListRes 获取体验合同列表响应
type ContractListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ContractListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// ContractExportReq 导出体验合同请求
type ContractExportReq struct {
	g.Meta    `path:"/contract/export" method:"get" tags:"体验合同" summary:"导出体验合同"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	ContractNo string `json:"contractNo" dc:"合同编号"`
	Title string `json:"title" dc:"合同标题"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	OrderID *snowflake.JsonInt64 `json:"orderID" dc:"订单"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	SignedAtStart string `json:"signedAtStart" dc:"签署时间开始时间"`
	SignedAtEnd string `json:"signedAtEnd" dc:"签署时间结束时间"`
	ExpiresAtStart string `json:"expiresAtStart" dc:"到期时间开始时间"`
	ExpiresAtEnd string `json:"expiresAtEnd" dc:"到期时间结束时间"`
}

// ContractExportRes 导出体验合同响应
type ContractExportRes struct {
	g.Meta `mime:"text/csv"`
}

// ContractImportReq 导入体验合同请求
type ContractImportReq struct {
	g.Meta `path:"/contract/import" method:"post" mime:"multipart/form-data" tags:"体验合同" summary:"导入体验合同"`
}

// ContractImportRes 导入体验合同响应
type ContractImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// ContractImportTemplateReq 下载体验合同导入模板
type ContractImportTemplateReq struct {
	g.Meta `path:"/contract/import-template" method:"get" tags:"体验合同" summary:"下载体验合同导入模板"`
}

// ContractImportTemplateRes 下载体验合同导入模板响应
type ContractImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

