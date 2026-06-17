package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemPlan is the golang structure of table system_plan for DAO operations like Where/Data.
type SystemPlan struct {
	g.Meta      `orm:"table:system_plan, do:true"`
	Id          any
	Name        any
	Code        any
	Description any
	UserLimit   any
	StorageMb   any
	Features    any
	Price       any
	Sort        any
	Status      any
	CreatedBy   any
	DeptId      any
	TenantId    any
	MerchantId  any
	CreatedAt   *gtime.Time
	UpdatedAt   *gtime.Time
	DeletedAt   *gtime.Time
}
