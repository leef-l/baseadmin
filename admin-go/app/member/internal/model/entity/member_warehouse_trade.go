// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWarehouseTrade is the golang structure for table member_warehouse_trade.
type MemberWarehouseTrade struct {
	Id           uint64      `orm:"id"            description:"ID（Snowflake）"`                                       // ID（Snowflake）
	TradeNo      string      `orm:"trade_no"      description:"交易编号|search:eq|keyword:on|priority:100"`              // 交易编号|search:eq|keyword:on|priority:100
	GoodsId      uint64      `orm:"goods_id"      description:"仓库商品|ref:member_warehouse_goods.title|search:select"` // 仓库商品|ref:member_warehouse_goods.title|search:select
	ListingId    uint64      `orm:"listing_id"    description:"挂卖记录|ref:member_warehouse_listing.id"`                // 挂卖记录|ref:member_warehouse_listing.id
	SellerId     uint64      `orm:"seller_id"     description:"卖家|ref:member_user.nickname|search:select"`           // 卖家|ref:member_user.nickname|search:select
	BuyerId      uint64      `orm:"buyer_id"      description:"买家|ref:member_user.nickname|search:select"`           // 买家|ref:member_user.nickname|search:select
	TradePrice   uint64      `orm:"trade_price"   description:"成交价格（分）"`                                             // 成交价格（分）
	PlatformFee  uint64      `orm:"platform_fee"  description:"平台扣除费用（分）"`                                           // 平台扣除费用（分）
	SellerIncome uint64      `orm:"seller_income" description:"卖家实收（分）"`                                             // 卖家实收（分）
	TradeStatus  int         `orm:"trade_status"  description:"交易状态:1=待卖家确认,2=已确认完成,3=已取消|search:select"`            // 交易状态:1=待卖家确认,2=已确认完成,3=已取消|search:select
	ConfirmedAt  *gtime.Time `orm:"confirmed_at"  description:"确认时间"`                                                // 确认时间
	Remark       string      `orm:"remark"        description:"备注|search:off"`                                       // 备注|search:off
	Status       int         `orm:"status"        description:"状态:0=关闭,1=开启|search:select"`                          // 状态:0=关闭,1=开启|search:select
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                                                  // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                                                  // 商户
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                                               // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                                              // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                                                // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                                                // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`                                  // 软删除时间，非 NULL 表示已删除
}
