package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// ShopOrder DTO 模型

// ShopOrderCreateInput 创建商城订单输入
type ShopOrderCreateInput struct {
	OrderNo string `json:"orderNo"`
	UserID snowflake.JsonInt64 `json:"userID"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	GoodsTitle string `json:"goodsTitle"`
	GoodsCover string `json:"goodsCover"`
	Quantity int `json:"quantity"`
	TotalPrice int64 `json:"totalPrice"`
	PayWallet int `json:"payWallet"`
	OrderStatus int `json:"orderStatus"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ShopOrderUpdateInput 更新商城订单输入
type ShopOrderUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	OrderNo string `json:"orderNo"`
	UserID snowflake.JsonInt64 `json:"userID"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	GoodsTitle string `json:"goodsTitle"`
	GoodsCover string `json:"goodsCover"`
	Quantity int `json:"quantity"`
	TotalPrice int64 `json:"totalPrice"`
	PayWallet int `json:"payWallet"`
	OrderStatus int `json:"orderStatus"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ShopOrderDetailOutput 商城订单详情输出
type ShopOrderDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	OrderNo string `json:"orderNo"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	ShopGoodsTitle string `json:"shopGoodsTitle"`
	GoodsTitle string `json:"goodsTitle"`
	GoodsCover string `json:"goodsCover"`
	Quantity int `json:"quantity"`
	TotalPrice int64 `json:"totalPrice"`
	PayWallet int `json:"payWallet"`
	OrderStatus int `json:"orderStatus"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ShopOrderListOutput 商城订单列表输出
type ShopOrderListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	OrderNo string `json:"orderNo"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	ShopGoodsTitle string `json:"shopGoodsTitle"`
	GoodsTitle string `json:"goodsTitle"`
	GoodsCover string `json:"goodsCover"`
	Quantity int `json:"quantity"`
	TotalPrice int64 `json:"totalPrice"`
	PayWallet int `json:"payWallet"`
	OrderStatus int `json:"orderStatus"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ShopOrderListInput 商城订单列表查询输入
type ShopOrderListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	OrderNo string `json:"orderNo"`
	UserID *snowflake.JsonInt64 `json:"userID"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	PayWallet *int `json:"payWallet"`
	OrderStatus *int `json:"orderStatus"`
	Status *int `json:"status"`
	GoodsTitle string `json:"goodsTitle"`
}

// ShopOrderBatchUpdateInput 批量编辑商城订单输入
type ShopOrderBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	PayWallet *int `json:"payWallet"`
	OrderStatus *int `json:"orderStatus"`
	Status *int `json:"status"`
}

