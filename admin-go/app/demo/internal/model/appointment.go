package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Appointment DTO 模型

// AppointmentCreateInput 创建体验预约输入
type AppointmentCreateInput struct {
	AppointmentNo string `json:"appointmentNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	Subject string `json:"subject"`
	AppointmentAt *gtime.Time `json:"appointmentAt"`
	ContactPhone string `json:"contactPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// AppointmentUpdateInput 更新体验预约输入
type AppointmentUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	AppointmentNo string `json:"appointmentNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	Subject string `json:"subject"`
	AppointmentAt *gtime.Time `json:"appointmentAt"`
	ContactPhone string `json:"contactPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// AppointmentDetailOutput 体验预约详情输出
type AppointmentDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	AppointmentNo string `json:"appointmentNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	Subject string `json:"subject"`
	AppointmentAt *gtime.Time `json:"appointmentAt"`
	ContactPhone string `json:"contactPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// AppointmentListOutput 体验预约列表输出
type AppointmentListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	AppointmentNo string `json:"appointmentNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	Subject string `json:"subject"`
	AppointmentAt *gtime.Time `json:"appointmentAt"`
	ContactPhone string `json:"contactPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// AppointmentListInput 体验预约列表查询输入
type AppointmentListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	AppointmentNo string `json:"appointmentNo"`
	Subject string `json:"subject"`
	ContactPhone string `json:"contactPhone"`
	CustomerID *snowflake.JsonInt64 `json:"customerID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
	AppointmentAtStart string `json:"appointmentAtStart"`
	AppointmentAtEnd string `json:"appointmentAtEnd"`
}

// AppointmentBatchUpdateInput 批量编辑体验预约输入
type AppointmentBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

