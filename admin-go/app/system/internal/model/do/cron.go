package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type SystemCron struct {
	g.Meta     `orm:"table:system_cron, do:true"`
	Id         any
	Name       any
	Expression any
	Handler    any
	Params     any
	Remark     any
	Status     any
	LastRunAt  *gtime.Time
	LastResult any
	CreatedBy  any
	DeptId     any
	TenantId   any
	MerchantId any
	CreatedAt  *gtime.Time
	UpdatedAt  *gtime.Time
	DeletedAt  *gtime.Time
}
