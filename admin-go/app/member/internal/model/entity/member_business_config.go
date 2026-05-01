// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberBusinessConfig is the golang structure for table member_business_config.
type MemberBusinessConfig struct {
	Id         uint64      `orm:"id"          description:"配置ID（Snowflake）"`                            // 配置ID（Snowflake）
	ConfigKey  string      `orm:"config_key"  description:"配置键|search:eq"`                              // 配置键|search:eq
	Payload    string      `orm:"payload"     description:"业务配置JSON（进货时间窗/寄售时间窗/工作日/返佣比例等）|search:off"` // 业务配置JSON（进货时间窗/寄售时间窗/工作日/返佣比例等）|search:off
	Remark     string      `orm:"remark"      description:"备注|search:off"`                              // 备注|search:off
	TenantId   uint64      `orm:"tenant_id"   description:"租户"`                                         // 租户
	MerchantId uint64      `orm:"merchant_id" description:"商户"`                                         // 商户
	CreatedBy  uint64      `orm:"created_by"  description:"创建人ID"`                                      // 创建人ID
	DeptId     uint64      `orm:"dept_id"     description:"所属部门ID"`                                     // 所属部门ID
	CreatedAt  *gtime.Time `orm:"created_at"  description:"创建时间"`                                       // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"  description:"更新时间"`                                       // 更新时间
	DeletedAt  *gtime.Time `orm:"deleted_at"  description:"软删除时间，非 NULL 表示已删除"`                         // 软删除时间，非 NULL 表示已删除
}
