// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberShopOrder is the golang structure for table member_shop_order.
type MemberShopOrder struct {
	Id          uint64      `orm:"id"           description:"ID（Snowflake）"`                                // ID（Snowflake）
	OrderNo     string      `orm:"order_no"     description:"订单号|search:eq|keyword:on|priority:100"`        // 订单号|search:eq|keyword:on|priority:100
	UserId      uint64      `orm:"user_id"      description:"购买会员|ref:member_user.nickname|search:select"`  // 购买会员|ref:member_user.nickname|search:select
	GoodsId     uint64      `orm:"goods_id"     description:"商品|ref:member_shop_goods.title|search:select"` // 商品|ref:member_shop_goods.title|search:select
	GoodsTitle  string      `orm:"goods_title"  description:"商品名称（快照）"`                                     // 商品名称（快照）
	GoodsCover  string      `orm:"goods_cover"  description:"商品封面（快照）"`                                     // 商品封面（快照）
	Quantity    uint        `orm:"quantity"     description:"购买数量"`                                         // 购买数量
	TotalPrice  uint64      `orm:"total_price"  description:"订单总价（分）"`                                      // 订单总价（分）
	PayWallet   int         `orm:"pay_wallet"   description:"支付钱包:1=优惠券余额"`                                 // 支付钱包:1=优惠券余额
	OrderStatus int         `orm:"order_status" description:"订单状态:1=已完成,2=已取消|search:select"`               // 订单状态:1=已完成,2=已取消|search:select
	Remark      string      `orm:"remark"       description:"订单备注|search:off"`                              // 订单备注|search:off
	Status      int         `orm:"status"       description:"状态:0=关闭,1=开启|search:select"`                   // 状态:0=关闭,1=开启|search:select
	TenantId    uint64      `orm:"tenant_id"    description:"租户"`                                           // 租户
	MerchantId  uint64      `orm:"merchant_id"  description:"商户"`                                           // 商户
	CreatedBy   uint64      `orm:"created_by"   description:"创建人ID"`                                        // 创建人ID
	DeptId      uint64      `orm:"dept_id"      description:"所属部门ID"`                                       // 所属部门ID
	CreatedAt   *gtime.Time `orm:"created_at"   description:"创建时间"`                                         // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"   description:"更新时间"`                                         // 更新时间
	DeletedAt   *gtime.Time `orm:"deleted_at"   description:"软删除时间，非 NULL 表示已删除"`                           // 软删除时间，非 NULL 表示已删除
}
