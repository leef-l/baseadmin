// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoSurveyDao is the data access object for the table demo_survey.
type DemoSurveyDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DemoSurveyColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DemoSurveyColumns defines and stores column names for the table demo_survey.
type DemoSurveyColumns struct {
	Id           string // 问卷ID（Snowflake）
	SurveyNo     string // 问卷编号|search:eq|priority:100
	Title        string // 问卷标题|search:like|keyword:on|priority:95
	Poster       string // 海报
	QuestionJson string // 问题JSON
	IntroContent string // 问卷介绍
	PublishAt    string // 发布时间
	ExpireAt     string // 过期时间
	IsAnonymous  string // 是否匿名:0=否,1=是
	Status       string // 状态:0=草稿,1=已发布,2=已下架
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// demoSurveyColumns holds the columns for the table demo_survey.
var demoSurveyColumns = DemoSurveyColumns{
	Id:           "id",
	SurveyNo:     "survey_no",
	Title:        "title",
	Poster:       "poster",
	QuestionJson: "question_json",
	IntroContent: "intro_content",
	PublishAt:    "publish_at",
	ExpireAt:     "expire_at",
	IsAnonymous:  "is_anonymous",
	Status:       "status",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewDemoSurveyDao creates and returns a new DAO object for table data access.
func NewDemoSurveyDao(handlers ...gdb.ModelHandler) *DemoSurveyDao {
	return &DemoSurveyDao{
		group:    "default",
		table:    "demo_survey",
		columns:  demoSurveyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoSurveyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoSurveyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoSurveyDao) Columns() DemoSurveyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoSurveyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoSurveyDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DemoSurveyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
