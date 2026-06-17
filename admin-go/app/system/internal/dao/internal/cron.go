// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type CronDao struct {
	table    string
	group    string
	columns  CronColumns
	handlers []gdb.ModelHandler
}

type CronColumns struct {
	Id         string
	Name       string
	Expression string
	Handler    string
	Params     string
	Remark     string
	Status     string
	LastRunAt  string
	LastResult string
	CreatedBy  string
	DeptId     string
	TenantId   string
	MerchantId string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}

var cronColumns = CronColumns{
	Id:         "id",
	Name:       "name",
	Expression: "expression",
	Handler:    "handler",
	Params:     "params",
	Remark:     "remark",
	Status:     "status",
	LastRunAt:  "last_run_at",
	LastResult: "last_result",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	TenantId:   "tenant_id",
	MerchantId: "merchant_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

func NewCronDao(handlers ...gdb.ModelHandler) *CronDao {
	return &CronDao{
		group:    "default",
		table:    "system_cron",
		columns:  cronColumns,
		handlers: handlers,
	}
}

func (dao *CronDao) DB() gdb.DB              { return g.DB(dao.group) }
func (dao *CronDao) Table() string            { return dao.table }
func (dao *CronDao) Columns() CronColumns     { return dao.columns }
func (dao *CronDao) Group() string            { return dao.group }
func (dao *CronDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}
func (dao *CronDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
