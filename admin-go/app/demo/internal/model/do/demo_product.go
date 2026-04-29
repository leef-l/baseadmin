// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoProduct is the golang structure of table demo_product for DAO operations like Where/Data.
type DemoProduct struct {
	g.Meta        `orm:"table:demo_product, do:true"`
	Id            any         // 商品ID（Snowflake）
	CategoryId    any         // 商品分类
	SkuNo         any         // SKU编号|search:eq|priority:100
	Name          any         // 商品名称|search:like|keyword:on|priority:95
	Cover         any         // 封面
	ManualFile    any         // 说明书文件
	DetailContent any         // 详情内容
	SpecJson      any         // 规格JSON
	WebsiteUrl    any         // 官网URL
	Type          any         // 类型:1=普通,2=置顶,3=推荐,4=热门
	IsRecommend   any         // 是否推荐:0=否,1=是
	SalePrice     any         // 销售价（分）
	StockNum      any         // 库存数量
	WeightNum     any         // 重量（克）
	Sort          any         // 排序（升序）
	Icon          any         // 图标
	Status        any         // 状态:0=草稿,1=上架,2=下架
	TenantId      any         // 租户
	MerchantId    any         // 商户
	CreatedBy     any         // 创建人ID
	DeptId        any         // 所属部门ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间，非 NULL 表示已删除
}
