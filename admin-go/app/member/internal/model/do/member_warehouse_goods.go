// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWarehouseGoods is the golang structure of table member_warehouse_goods for DAO operations like Where/Data.
type MemberWarehouseGoods struct {
	g.Meta          `orm:"table:member_warehouse_goods, do:true"`
	Id              any         // ID（Snowflake）
	GoodsNo         any         // 商品编号|search:eq|keyword:on|priority:100
	Title           any         // 商品名称|search:like|keyword:on|priority:95
	Cover           any         // 商品封面
	InitPrice       any         // 初始价格（分）
	CurrentPrice    any         // 当前价格（分）
	PriceRiseRate   any         // 每次加价比例（百分比，如10=10%）
	PlatformFeeRate any         // 平台扣除比例（百分比，如5=5%）
	OwnerId         any         // 当前持有人|ref:member_user.nickname|search:select
	TradeCount      any         // 流转次数
	GoodsStatus     any         // 商品状态:1=持有中,2=挂卖中,3=交易中|search:select
	Remark          any         // 备注|search:off
	Sort            any         // 排序（升序）
	Status          any         // 状态:0=关闭,1=开启|search:select
	TenantId        any         // 租户
	MerchantId      any         // 商户
	CreatedBy       any         // 创建人ID
	DeptId          any         // 所属部门ID
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
	DeletedAt       *gtime.Time // 软删除时间，非 NULL 表示已删除
}
