// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWarehouseListing is the golang structure of table member_warehouse_listing for DAO operations like Where/Data.
type MemberWarehouseListing struct {
	g.Meta        `orm:"table:member_warehouse_listing, do:true"`
	Id            any         // ID（Snowflake）
	GoodsId       any         // 仓库商品|ref:member_warehouse_goods.title|search:select
	SellerId      any         // 卖家|ref:member_user.nickname|search:select
	ListingPrice  any         // 挂卖价格（分，自动加价后）
	ListingStatus any         // 挂卖状态:1=挂卖中,2=已售出,3=已取消|search:select
	ListedAt      *gtime.Time // 挂卖时间
	SoldAt        *gtime.Time // 售出时间
	Remark        any         // 备注|search:off
	Status        any         // 状态:0=关闭,1=开启|search:select
	TenantId      any         // 租户
	MerchantId    any         // 商户
	CreatedBy     any         // 创建人ID
	DeptId        any         // 所属部门ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间，非 NULL 表示已删除
}
