// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoContract is the golang structure of table demo_contract for DAO operations like Where/Data.
type DemoContract struct {
	g.Meta         `orm:"table:demo_contract, do:true"`
	Id             any         // 合同ID（Snowflake）
	ContractNo     any         // 合同编号|search:eq|priority:100
	CustomerId     any         // 客户
	OrderId        any         // 订单
	Title          any         // 合同标题|search:like|keyword:on|priority:95
	ContractFile   any         // 合同文件
	SignImage      any         // 签章图片
	ContractAmount any         // 合同金额（分）
	SignPassword   any         // 签署密码
	SignedAt       *gtime.Time // 签署时间
	ExpiresAt      *gtime.Time // 到期时间
	Status         any         // 状态:0=待审核,1=已通过,2=已拒绝,3=已取消
	TenantId       any         // 租户
	MerchantId     any         // 商户
	CreatedBy      any         // 创建人ID
	DeptId         any         // 所属部门ID
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	DeletedAt      *gtime.Time // 软删除时间，非 NULL 表示已删除
}
