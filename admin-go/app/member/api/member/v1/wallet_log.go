package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// WalletLog API

// WalletLogCreateReq 创建钱包流水记录请求
type WalletLogCreateReq struct {
	g.Meta `path:"/wallet_log/create" method:"post" tags:"钱包流水记录" summary:"创建钱包流水记录"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	WalletType int `json:"walletType"  dc:"钱包类型"`
	ChangeType int `json:"changeType"  dc:"变动类型"`
	ChangeAmount int64 `json:"changeAmount"  dc:"变动金额（分，正增负减）"`
	BeforeBalance int64 `json:"beforeBalance"  dc:"变动前余额（分）"`
	AfterBalance int64 `json:"afterBalance"  dc:"变动后余额（分）"`
	RelatedOrderNo string `json:"relatedOrderNo" v:"max-length:64" dc:"关联单号"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注说明"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WalletLogCreateRes 创建钱包流水记录响应
type WalletLogCreateRes struct {
	g.Meta `mime:"application/json"`
}

// WalletLogUpdateReq 更新钱包流水记录请求
type WalletLogUpdateReq struct {
	g.Meta `path:"/wallet_log/update" method:"put" tags:"钱包流水记录" summary:"更新钱包流水记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"钱包流水记录ID"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	WalletType int `json:"walletType"  dc:"钱包类型"`
	ChangeType int `json:"changeType"  dc:"变动类型"`
	ChangeAmount int64 `json:"changeAmount"  dc:"变动金额（分，正增负减）"`
	BeforeBalance int64 `json:"beforeBalance"  dc:"变动前余额（分）"`
	AfterBalance int64 `json:"afterBalance"  dc:"变动后余额（分）"`
	RelatedOrderNo string `json:"relatedOrderNo" v:"max-length:64" dc:"关联单号"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注说明"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WalletLogUpdateRes 更新钱包流水记录响应
type WalletLogUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WalletLogDeleteReq 删除钱包流水记录请求
type WalletLogDeleteReq struct {
	g.Meta `path:"/wallet_log/delete" method:"delete" tags:"钱包流水记录" summary:"删除钱包流水记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"钱包流水记录ID"`
}

// WalletLogDeleteRes 删除钱包流水记录响应
type WalletLogDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WalletLogBatchDeleteReq 批量删除钱包流水记录请求
type WalletLogBatchDeleteReq struct {
	g.Meta `path:"/wallet_log/batch-delete" method:"delete" tags:"钱包流水记录" summary:"批量删除钱包流水记录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"钱包流水记录ID列表"`
}

// WalletLogBatchDeleteRes 批量删除钱包流水记录响应
type WalletLogBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WalletLogDetailReq 获取钱包流水记录详情请求
type WalletLogDetailReq struct {
	g.Meta `path:"/wallet_log/detail" method:"get" tags:"钱包流水记录" summary:"获取钱包流水记录详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"钱包流水记录ID"`
}

// WalletLogDetailRes 获取钱包流水记录详情响应
type WalletLogDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.WalletLogDetailOutput
}

// WalletLogListReq 获取钱包流水记录列表请求
type WalletLogListReq struct {
	g.Meta    `path:"/wallet_log/list" method:"get" tags:"钱包流水记录" summary:"获取钱包流水记录列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	RelatedOrderNo string `json:"relatedOrderNo" dc:"关联单号"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	WalletType *int `json:"walletType" dc:"钱包类型"`
	ChangeType *int `json:"changeType" dc:"变动类型"`
}

// WalletLogListRes 获取钱包流水记录列表响应
type WalletLogListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.WalletLogListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// WalletLogExportReq 导出钱包流水记录请求
type WalletLogExportReq struct {
	g.Meta    `path:"/wallet_log/export" method:"get" tags:"钱包流水记录" summary:"导出钱包流水记录"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	RelatedOrderNo string `json:"relatedOrderNo" dc:"关联单号"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	WalletType *int `json:"walletType" dc:"钱包类型"`
	ChangeType *int `json:"changeType" dc:"变动类型"`
}

// WalletLogExportRes 导出钱包流水记录响应
type WalletLogExportRes struct {
	g.Meta `mime:"text/csv"`
}

// WalletLogImportReq 导入钱包流水记录请求
type WalletLogImportReq struct {
	g.Meta `path:"/wallet_log/import" method:"post" mime:"multipart/form-data" tags:"钱包流水记录" summary:"导入钱包流水记录"`
}

// WalletLogImportRes 导入钱包流水记录响应
type WalletLogImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// WalletLogImportTemplateReq 下载钱包流水记录导入模板
type WalletLogImportTemplateReq struct {
	g.Meta `path:"/wallet_log/import-template" method:"get" tags:"钱包流水记录" summary:"下载钱包流水记录导入模板"`
}

// WalletLogImportTemplateRes 下载钱包流水记录导入模板响应
type WalletLogImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

