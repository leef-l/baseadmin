// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberShopGoods is the golang structure for table member_shop_goods.
type MemberShopGoods struct {
	Id            uint64      `orm:"id"             description:"ID（Snowflake）"`                                    // ID（Snowflake）
	CategoryId    uint64      `orm:"category_id"    description:"商品分类|ref:member_shop_category.name|search:select"` // 商品分类|ref:member_shop_category.name|search:select
	Title         string      `orm:"title"          description:"商品名称|search:like|keyword:on|priority:100"`         // 商品名称|search:like|keyword:on|priority:100
	Cover         string      `orm:"cover"          description:"封面图"`                                              // 封面图
	Images        string      `orm:"images"         description:"商品图片（JSON数组）|search:off"`                          // 商品图片（JSON数组）|search:off
	Price         uint64      `orm:"price"          description:"售价（分，优惠券余额支付）"`                                    // 售价（分，优惠券余额支付）
	OriginalPrice uint64      `orm:"original_price" description:"原价（分）"`                                            // 原价（分）
	Stock         uint        `orm:"stock"          description:"库存"`                                               // 库存
	Sales         uint        `orm:"sales"          description:"销量"`                                               // 销量
	Content       string      `orm:"content"        description:"商品详情|search:off"`                                  // 商品详情|search:off
	Sort          int         `orm:"sort"           description:"排序（升序）"`                                           // 排序（升序）
	IsRecommend   int         `orm:"is_recommend"   description:"是否推荐:0=否,1=是|search:select"`                       // 是否推荐:0=否,1=是|search:select
	Status        int         `orm:"status"         description:"状态:0=下架,1=上架|search:select"`                       // 状态:0=下架,1=上架|search:select
	TenantId      uint64      `orm:"tenant_id"      description:"租户"`                                               // 租户
	MerchantId    uint64      `orm:"merchant_id"    description:"商户"`                                               // 商户
	CreatedBy     uint64      `orm:"created_by"     description:"创建人ID"`                                            // 创建人ID
	DeptId        uint64      `orm:"dept_id"        description:"所属部门ID"`                                           // 所属部门ID
	CreatedAt     *gtime.Time `orm:"created_at"     description:"创建时间"`                                             // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     description:"更新时间"`                                             // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at"     description:"软删除时间，非 NULL 表示已删除"`                               // 软删除时间，非 NULL 表示已删除
}
