// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWarehouseListing is the golang structure for table member_warehouse_listing.
type MemberWarehouseListing struct {
	Id            uint64      `orm:"id"             description:"ID（Snowflake）"`                                       // ID（Snowflake）
	GoodsId       uint64      `orm:"goods_id"       description:"仓库商品|ref:member_warehouse_goods.title|search:select"` // 仓库商品|ref:member_warehouse_goods.title|search:select
	SellerId      uint64      `orm:"seller_id"      description:"卖家|ref:member_user.nickname|search:select"`           // 卖家|ref:member_user.nickname|search:select
	ListingPrice  uint64      `orm:"listing_price"  description:"挂卖价格（分，自动加价后）"`                                       // 挂卖价格（分，自动加价后）
	ListingStatus int         `orm:"listing_status" description:"挂卖状态:1=挂卖中,2=已售出,3=已取消|search:select"`                // 挂卖状态:1=挂卖中,2=已售出,3=已取消|search:select
	ListedAt      *gtime.Time `orm:"listed_at"      description:"挂卖时间"`                                                // 挂卖时间
	SoldAt        *gtime.Time `orm:"sold_at"        description:"售出时间"`                                                // 售出时间
	Remark        string      `orm:"remark"         description:"备注|search:off"`                                       // 备注|search:off
	Status        int         `orm:"status"         description:"状态:0=关闭,1=开启|search:select"`                          // 状态:0=关闭,1=开启|search:select
	TenantId      uint64      `orm:"tenant_id"      description:"租户"`                                                  // 租户
	MerchantId    uint64      `orm:"merchant_id"    description:"商户"`                                                  // 商户
	CreatedBy     uint64      `orm:"created_by"     description:"创建人ID"`                                               // 创建人ID
	DeptId        uint64      `orm:"dept_id"        description:"所属部门ID"`                                              // 所属部门ID
	CreatedAt     *gtime.Time `orm:"created_at"     description:"创建时间"`                                                // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     description:"更新时间"`                                                // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at"     description:"软删除时间，非 NULL 表示已删除"`                                  // 软删除时间，非 NULL 表示已删除
}
