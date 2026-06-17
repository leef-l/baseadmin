package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type SystemTenantPlan struct {
	g.Meta     `orm:"table:system_tenant_plan, do:true"`
	Id         any
	TenantId   any
	PlanId     any
	StartAt    *gtime.Time
	ExpireAt   *gtime.Time
	Status     any
	Remark     any
	CreatedBy  any
	DeptId     any
	MerchantId any
	CreatedAt  *gtime.Time
	UpdatedAt  *gtime.Time
	DeletedAt  *gtime.Time
}
