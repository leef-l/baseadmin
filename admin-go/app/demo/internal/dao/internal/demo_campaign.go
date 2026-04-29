// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoCampaignDao is the data access object for the table demo_campaign.
type DemoCampaignDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DemoCampaignColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DemoCampaignColumns defines and stores column names for the table demo_campaign.
type DemoCampaignColumns struct {
	Id           string // 活动ID（Snowflake）
	CampaignNo   string // 活动编号|search:eq|priority:100
	Title        string // 活动标题|search:like|keyword:on|priority:95
	Banner       string // 横幅图
	Type         string // 活动类型:1=免费,2=付费,3=公开,4=私密
	Channel      string // 投放渠道:1=官网,2=小程序,3=短信,4=线下
	BudgetAmount string // 预算金额（分）
	LandingUrl   string // 落地页URL
	RuleJson     string // 规则JSON
	IntroContent string // 活动介绍
	StartAt      string // 开始时间
	EndAt        string // 结束时间
	IsPublic     string // 是否公开:0=否,1=是
	Status       string // 状态:0=草稿,1=已发布,2=已下架
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// demoCampaignColumns holds the columns for the table demo_campaign.
var demoCampaignColumns = DemoCampaignColumns{
	Id:           "id",
	CampaignNo:   "campaign_no",
	Title:        "title",
	Banner:       "banner",
	Type:         "type",
	Channel:      "channel",
	BudgetAmount: "budget_amount",
	LandingUrl:   "landing_url",
	RuleJson:     "rule_json",
	IntroContent: "intro_content",
	StartAt:      "start_at",
	EndAt:        "end_at",
	IsPublic:     "is_public",
	Status:       "status",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewDemoCampaignDao creates and returns a new DAO object for table data access.
func NewDemoCampaignDao(handlers ...gdb.ModelHandler) *DemoCampaignDao {
	return &DemoCampaignDao{
		group:    "default",
		table:    "demo_campaign",
		columns:  demoCampaignColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoCampaignDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoCampaignDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoCampaignDao) Columns() DemoCampaignColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoCampaignDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoCampaignDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoCampaignDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
