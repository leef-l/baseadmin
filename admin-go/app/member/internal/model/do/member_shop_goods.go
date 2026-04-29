// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberShopGoods is the golang structure of table member_shop_goods for DAO operations like Where/Data.
type MemberShopGoods struct {
	g.Meta        `orm:"table:member_shop_goods, do:true"`
	Id            any         // ID（Snowflake）
	CategoryId    any         // 商品分类|ref:member_shop_category.name|search:select
	Title         any         // 商品名称|search:like|keyword:on|priority:100
	Cover         any         // 封面图
	Images        any         // 商品图片（JSON数组）|search:off
	Price         any         // 售价（分，优惠券余额支付）
	OriginalPrice any         // 原价（分）
	Stock         any         // 库存
	Sales         any         // 销量
	Content       any         // 商品详情|search:off
	Sort          any         // 排序（升序）
	IsRecommend   any         // 是否推荐:0=否,1=是|search:select
	Status        any         // 状态:0=下架,1=上架|search:select
	TenantId      any         // 租户
	MerchantId    any         // 商户
	CreatedBy     any         // 创建人ID
	DeptId        any         // 所属部门ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间，非 NULL 表示已删除
}
