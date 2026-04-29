// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Domain is the golang structure of table system_domain for DAO operations like Where/Data.
type Domain struct {
	g.Meta       `orm:"table:system_domain, do:true"`
	Id           any         // 域名ID（Snowflake）
	Domain       any         // 绑定域名
	OwnerType    any         // 主体类型:1=租户,2=商户
	TenantId     any         // 租户
	MerchantId   any         // 商户
	AppCode      any         // 应用编码：admin/upload/shop
	VerifyToken  any         // 域名校验令牌
	VerifyStatus any         // 校验状态:0=未校验,1=已校验
	SslStatus    any         // SSL状态:0=未配置,1=已配置
	NginxStatus  any         // Nginx配置状态:0=未应用,1=已应用
	Status       any         // 状态:0=关闭,1=开启
	Remark       any         // 备注
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
