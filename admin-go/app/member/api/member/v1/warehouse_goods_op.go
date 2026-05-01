package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// WarehouseGoodsAssignReq 后台把仓库商品分配给指定会员。
//
// 业务约束：
//   - 商品必须存在且未删除
//   - 商品当前状态必须是 1=持有中（挂卖中 / 交易中不允许换持有人）
//   - 目标会员必须存在、status=1
//
// 不做钱包变动；只改 owner_id。后续会员可在 H5 一键挂卖。
type WarehouseGoodsAssignReq struct {
	g.Meta  `path:"/warehouse_goods/assign" method:"post" tags:"仓库商品" summary:"分配仓库商品给指定会员"`
	GoodsID snowflake.JsonInt64 `json:"goodsId" v:"required#商品 ID 不能为空"`
	OwnerID snowflake.JsonInt64 `json:"ownerId" v:"required#目标会员 ID 不能为空"`
	Remark  string              `json:"remark" v:"max-length:500" dc:"分配备注"`
}

// WarehouseGoodsAssignRes 响应。
type WarehouseGoodsAssignRes struct {
	g.Meta `mime:"application/json"`
}
