package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Cron DTO 模型

// CronCreateInput 创建定时任务表输入
type CronCreateInput struct {
	Name string `json:"name"`
	Expression string `json:"expression"`
	Handler string `json:"handler"`
	Params string `json:"params"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	LastRunAt *gtime.Time `json:"lastRunAt"`
	LastResult string `json:"lastResult"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CronUpdateInput 更新定时任务表输入
type CronUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	Expression string `json:"expression"`
	Handler string `json:"handler"`
	Params string `json:"params"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	LastRunAt *gtime.Time `json:"lastRunAt"`
	LastResult string `json:"lastResult"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CronDetailOutput 定时任务表详情输出
type CronDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	Expression string `json:"expression"`
	Handler string `json:"handler"`
	Params string `json:"params"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	LastRunAt *gtime.Time `json:"lastRunAt"`
	LastResult string `json:"lastResult"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// CronListOutput 定时任务表列表输出
type CronListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	Expression string `json:"expression"`
	Handler string `json:"handler"`
	Params string `json:"params"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	LastRunAt *gtime.Time `json:"lastRunAt"`
	LastResult string `json:"lastResult"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// CronListInput 定时任务表列表查询输入
type CronListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	Name string `json:"name"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
	Remark string `json:"remark"`
	LastRunAtStart string `json:"lastRunAtStart"`
	LastRunAtEnd string `json:"lastRunAtEnd"`
}

// CronBatchUpdateInput 批量编辑定时任务表输入
type CronBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

// CronLogItem 定时任务执行日志
type CronLogItem struct {
	ID         int64       `json:"id"`
	CronID     int64       `json:"cronId"`
	CronName   string      `json:"cronName"`
	StartAt    *gtime.Time `json:"startAt"`
	EndAt      *gtime.Time `json:"endAt"`
	DurationMs int         `json:"durationMs"`
	Result     string      `json:"result"`
	Message    string      `json:"message"`
}

