// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWarehouseGoods is the golang structure for table member_warehouse_goods.
type MemberWarehouseGoods struct {
	Id              uint64      `orm:"id"                description:"ID（Snowflake）"`                                // ID（Snowflake）
	GoodsNo         string      `orm:"goods_no"          description:"商品编号|search:eq|keyword:on|priority:100"`       // 商品编号|search:eq|keyword:on|priority:100
	Title           string      `orm:"title"             description:"商品名称|search:like|keyword:on|priority:95"`      // 商品名称|search:like|keyword:on|priority:95
	Cover           string      `orm:"cover"             description:"商品封面"`                                         // 商品封面
	InitPrice       uint64      `orm:"init_price"        description:"初始价格（分）"`                                      // 初始价格（分）
	CurrentPrice    uint64      `orm:"current_price"     description:"当前价格（分）"`                                      // 当前价格（分）
	PriceRiseRate   uint        `orm:"price_rise_rate"   description:"每次加价比例（百分比，如10=10%）"`                          // 每次加价比例（百分比，如10=10%）
	PlatformFeeRate uint        `orm:"platform_fee_rate" description:"平台扣除比例（百分比，如5=5%）"`                            // 平台扣除比例（百分比，如5=5%）
	OwnerId         uint64      `orm:"owner_id"          description:"当前持有人|ref:member_user.nickname|search:select"` // 当前持有人|ref:member_user.nickname|search:select
	TradeCount      uint        `orm:"trade_count"       description:"流转次数"`                                         // 流转次数
	GoodsStatus     int         `orm:"goods_status"      description:"商品状态:1=持有中,2=挂卖中,3=交易中|search:select"`         // 商品状态:1=持有中,2=挂卖中,3=交易中|search:select
	Remark          string      `orm:"remark"            description:"备注|search:off"`                                // 备注|search:off
	Sort            int         `orm:"sort"              description:"排序（升序）"`                                       // 排序（升序）
	Status          int         `orm:"status"            description:"状态:0=关闭,1=开启|search:select"`                   // 状态:0=关闭,1=开启|search:select
	TenantId        uint64      `orm:"tenant_id"         description:"租户"`                                           // 租户
	MerchantId      uint64      `orm:"merchant_id"       description:"商户"`                                           // 商户
	CreatedBy       uint64      `orm:"created_by"        description:"创建人ID"`                                        // 创建人ID
	DeptId          uint64      `orm:"dept_id"           description:"所属部门ID"`                                       // 所属部门ID
	CreatedAt       *gtime.Time `orm:"created_at"        description:"创建时间"`                                         // 创建时间
	UpdatedAt       *gtime.Time `orm:"updated_at"        description:"更新时间"`                                         // 更新时间
	DeletedAt       *gtime.Time `orm:"deleted_at"        description:"软删除时间，非 NULL 表示已删除"`                           // 软删除时间，非 NULL 表示已删除
}
