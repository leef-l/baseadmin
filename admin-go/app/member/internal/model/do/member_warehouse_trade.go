// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWarehouseTrade is the golang structure of table member_warehouse_trade for DAO operations like Where/Data.
type MemberWarehouseTrade struct {
	g.Meta       `orm:"table:member_warehouse_trade, do:true"`
	Id           any         // ID（Snowflake）
	TradeNo      any         // 交易编号|search:eq|keyword:on|priority:100
	GoodsId      any         // 仓库商品|ref:member_warehouse_goods.title|search:select
	ListingId    any         // 挂卖记录|ref:member_warehouse_listing.id
	SellerId     any         // 卖家|ref:member_user.nickname|search:select
	BuyerId      any         // 买家|ref:member_user.nickname|search:select
	TradePrice   any         // 成交价格（分）
	PlatformFee  any         // 平台扣除费用（分）
	SellerIncome any         // 卖家实收（分）
	TradeStatus  any         // 交易状态:1=待卖家确认,2=已确认完成,3=已取消|search:select
	ConfirmedAt  *gtime.Time // 确认时间
	Remark       any         // 备注|search:off
	Status       any         // 状态:0=关闭,1=开启|search:select
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
