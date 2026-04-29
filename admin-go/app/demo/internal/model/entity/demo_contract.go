// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoContract is the golang structure for table demo_contract.
type DemoContract struct {
	Id             uint64      `orm:"id"              description:"合同ID（Snowflake）"`                         // 合同ID（Snowflake）
	ContractNo     string      `orm:"contract_no"     description:"合同编号|search:eq|priority:100"`             // 合同编号|search:eq|priority:100
	CustomerId     uint64      `orm:"customer_id"     description:"客户"`                                      // 客户
	OrderId        uint64      `orm:"order_id"        description:"订单"`                                      // 订单
	Title          string      `orm:"title"           description:"合同标题|search:like|keyword:on|priority:95"` // 合同标题|search:like|keyword:on|priority:95
	ContractFile   string      `orm:"contract_file"   description:"合同文件"`                                    // 合同文件
	SignImage      string      `orm:"sign_image"      description:"签章图片"`                                    // 签章图片
	ContractAmount int         `orm:"contract_amount" description:"合同金额（分）"`                                 // 合同金额（分）
	SignPassword   string      `orm:"sign_password"   description:"签署密码"`                                    // 签署密码
	SignedAt       *gtime.Time `orm:"signed_at"       description:"签署时间"`                                    // 签署时间
	ExpiresAt      *gtime.Time `orm:"expires_at"      description:"到期时间"`                                    // 到期时间
	Status         int         `orm:"status"          description:"状态:0=待审核,1=已通过,2=已拒绝,3=已取消"`              // 状态:0=待审核,1=已通过,2=已拒绝,3=已取消
	TenantId       uint64      `orm:"tenant_id"       description:"租户"`                                      // 租户
	MerchantId     uint64      `orm:"merchant_id"     description:"商户"`                                      // 商户
	CreatedBy      uint64      `orm:"created_by"      description:"创建人ID"`                                   // 创建人ID
	DeptId         uint64      `orm:"dept_id"         description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt      *gtime.Time `orm:"created_at"      description:"创建时间"`                                    // 创建时间
	UpdatedAt      *gtime.Time `orm:"updated_at"      description:"更新时间"`                                    // 更新时间
	DeletedAt      *gtime.Time `orm:"deleted_at"      description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
