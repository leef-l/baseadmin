// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoProduct is the golang structure for table demo_product.
type DemoProduct struct {
	Id            uint64      `orm:"id"             description:"商品ID（Snowflake）"`                         // 商品ID（Snowflake）
	CategoryId    uint64      `orm:"category_id"    description:"商品分类"`                                    // 商品分类
	SkuNo         string      `orm:"sku_no"         description:"SKU编号|search:eq|priority:100"`            // SKU编号|search:eq|priority:100
	Name          string      `orm:"name"           description:"商品名称|search:like|keyword:on|priority:95"` // 商品名称|search:like|keyword:on|priority:95
	Cover         string      `orm:"cover"          description:"封面"`                                      // 封面
	ManualFile    string      `orm:"manual_file"    description:"说明书文件"`                                   // 说明书文件
	DetailContent string      `orm:"detail_content" description:"详情内容"`                                    // 详情内容
	SpecJson      string      `orm:"spec_json"      description:"规格JSON"`                                  // 规格JSON
	WebsiteUrl    string      `orm:"website_url"    description:"官网URL"`                                   // 官网URL
	Type          int         `orm:"type"           description:"类型:1=普通,2=置顶,3=推荐,4=热门"`                  // 类型:1=普通,2=置顶,3=推荐,4=热门
	IsRecommend   int         `orm:"is_recommend"   description:"是否推荐:0=否,1=是"`                            // 是否推荐:0=否,1=是
	SalePrice     int         `orm:"sale_price"     description:"销售价（分）"`                                  // 销售价（分）
	StockNum      int         `orm:"stock_num"      description:"库存数量"`                                    // 库存数量
	WeightNum     int         `orm:"weight_num"     description:"重量（克）"`                                   // 重量（克）
	Sort          int         `orm:"sort"           description:"排序（升序）"`                                  // 排序（升序）
	Icon          string      `orm:"icon"           description:"图标"`                                      // 图标
	Status        int         `orm:"status"         description:"状态:0=草稿,1=上架,2=下架"`                       // 状态:0=草稿,1=上架,2=下架
	TenantId      uint64      `orm:"tenant_id"      description:"租户"`                                      // 租户
	MerchantId    uint64      `orm:"merchant_id"    description:"商户"`                                      // 商户
	CreatedBy     uint64      `orm:"created_by"     description:"创建人ID"`                                   // 创建人ID
	DeptId        uint64      `orm:"dept_id"        description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt     *gtime.Time `orm:"created_at"     description:"创建时间"`                                    // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     description:"更新时间"`                                    // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at"     description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
