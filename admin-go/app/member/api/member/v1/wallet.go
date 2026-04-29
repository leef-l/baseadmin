package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Wallet API

// WalletCreateReq 创建会员钱包请求
type WalletCreateReq struct {
	g.Meta `path:"/wallet/create" method:"post" tags:"会员钱包" summary:"创建会员钱包"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	WalletType int `json:"walletType"  dc:"钱包类型"`
	Balance int64 `json:"balance"  dc:"当前余额（分）"`
	TotalIncome int64 `json:"totalIncome"  dc:"累计收入（分）"`
	TotalExpense int64 `json:"totalExpense"  dc:"累计支出（分）"`
	FrozenAmount int64 `json:"frozenAmount"  dc:"冻结金额（分）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WalletCreateRes 创建会员钱包响应
type WalletCreateRes struct {
	g.Meta `mime:"application/json"`
}

// WalletUpdateReq 更新会员钱包请求
type WalletUpdateReq struct {
	g.Meta `path:"/wallet/update" method:"put" tags:"会员钱包" summary:"更新会员钱包"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员钱包ID"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	WalletType int `json:"walletType"  dc:"钱包类型"`
	Balance int64 `json:"balance"  dc:"当前余额（分）"`
	TotalIncome int64 `json:"totalIncome"  dc:"累计收入（分）"`
	TotalExpense int64 `json:"totalExpense"  dc:"累计支出（分）"`
	FrozenAmount int64 `json:"frozenAmount"  dc:"冻结金额（分）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WalletUpdateRes 更新会员钱包响应
type WalletUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WalletDeleteReq 删除会员钱包请求
type WalletDeleteReq struct {
	g.Meta `path:"/wallet/delete" method:"delete" tags:"会员钱包" summary:"删除会员钱包"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员钱包ID"`
}

// WalletDeleteRes 删除会员钱包响应
type WalletDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WalletBatchDeleteReq 批量删除会员钱包请求
type WalletBatchDeleteReq struct {
	g.Meta `path:"/wallet/batch-delete" method:"delete" tags:"会员钱包" summary:"批量删除会员钱包"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员钱包ID列表"`
}

// WalletBatchDeleteRes 批量删除会员钱包响应
type WalletBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WalletBatchUpdateReq 批量编辑会员钱包请求
type WalletBatchUpdateReq struct {
	g.Meta `path:"/wallet/batch-update" method:"put" tags:"会员钱包" summary:"批量编辑会员钱包"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员钱包ID列表"`
	WalletType *int `json:"walletType" dc:"钱包类型"`
	Status *int `json:"status" dc:"状态"`
}

// WalletBatchUpdateRes 批量编辑会员钱包响应
type WalletBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WalletDetailReq 获取会员钱包详情请求
type WalletDetailReq struct {
	g.Meta `path:"/wallet/detail" method:"get" tags:"会员钱包" summary:"获取会员钱包详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员钱包ID"`
}

// WalletDetailRes 获取会员钱包详情响应
type WalletDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.WalletDetailOutput
}

// WalletListReq 获取会员钱包列表请求
type WalletListReq struct {
	g.Meta    `path:"/wallet/list" method:"get" tags:"会员钱包" summary:"获取会员钱包列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	WalletType *int `json:"walletType" dc:"钱包类型"`
	Status *int `json:"status" dc:"状态"`
}

// WalletListRes 获取会员钱包列表响应
type WalletListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.WalletListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// WalletExportReq 导出会员钱包请求
type WalletExportReq struct {
	g.Meta    `path:"/wallet/export" method:"get" tags:"会员钱包" summary:"导出会员钱包"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	WalletType *int `json:"walletType" dc:"钱包类型"`
	Status *int `json:"status" dc:"状态"`
}

// WalletExportRes 导出会员钱包响应
type WalletExportRes struct {
	g.Meta `mime:"text/csv"`
}

// WalletImportReq 导入会员钱包请求
type WalletImportReq struct {
	g.Meta `path:"/wallet/import" method:"post" mime:"multipart/form-data" tags:"会员钱包" summary:"导入会员钱包"`
}

// WalletImportRes 导入会员钱包响应
type WalletImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// WalletImportTemplateReq 下载会员钱包导入模板
type WalletImportTemplateReq struct {
	g.Meta `path:"/wallet/import-template" method:"get" tags:"会员钱包" summary:"下载会员钱包导入模板"`
}

// WalletImportTemplateRes 下载会员钱包导入模板响应
type WalletImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

