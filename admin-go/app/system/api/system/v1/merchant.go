package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// MerchantCreateReq 创建商户请求
type MerchantCreateReq struct {
	g.Meta        `path:"/merchant/create" method:"post" tags:"商户表" summary:"创建商户"`
	TenantID      snowflake.JsonInt64 `json:"tenantId" dc:"租户ID"`
	Name          string              `json:"name" v:"required#商户名称不能为空" dc:"商户名称"`
	Code          string              `json:"code" v:"required#商户编码不能为空" dc:"商户编码"`
	ContactName   string              `json:"contactName" dc:"联系人"`
	ContactPhone  string              `json:"contactPhone" dc:"联系电话"`
	Address       string              `json:"address" dc:"商户地址"`
	Status        int                 `json:"status" dc:"状态"`
	Remark        string              `json:"remark" dc:"备注"`
	CreateAdmin   int                 `json:"createAdmin" dc:"是否同步创建管理员:0=否,1=是"`
	AdminUsername string              `json:"adminUsername" dc:"管理员用户名"`
	AdminPassword string              `json:"adminPassword" dc:"管理员密码"`
	AdminNickname string              `json:"adminNickname" dc:"管理员昵称"`
	AdminEmail    string              `json:"adminEmail" dc:"管理员邮箱"`
}

type MerchantCreateRes struct {
	g.Meta `mime:"application/json"`
}

// MerchantUpdateReq 更新商户请求
type MerchantUpdateReq struct {
	g.Meta       `path:"/merchant/update" method:"put" tags:"商户表" summary:"更新商户"`
	ID           snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商户ID"`
	TenantID     snowflake.JsonInt64 `json:"tenantId" dc:"租户ID"`
	Name         string              `json:"name" v:"required#商户名称不能为空" dc:"商户名称"`
	Code         string              `json:"code" v:"required#商户编码不能为空" dc:"商户编码"`
	ContactName  string              `json:"contactName" dc:"联系人"`
	ContactPhone string              `json:"contactPhone" dc:"联系电话"`
	Address      string              `json:"address" dc:"商户地址"`
	Status       int                 `json:"status" dc:"状态"`
	Remark       string              `json:"remark" dc:"备注"`
}

type MerchantUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// MerchantDeleteReq 删除商户请求
type MerchantDeleteReq struct {
	g.Meta `path:"/merchant/delete" method:"delete" tags:"商户表" summary:"删除商户"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商户ID"`
}

type MerchantDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// MerchantBatchDeleteReq 批量删除商户请求
type MerchantBatchDeleteReq struct {
	g.Meta `path:"/merchant/batch-delete" method:"delete" tags:"商户表" summary:"批量删除商户"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"商户ID列表"`
}

type MerchantBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// MerchantDetailReq 获取商户详情请求
type MerchantDetailReq struct {
	g.Meta `path:"/merchant/detail" method:"get" tags:"商户表" summary:"获取商户详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商户ID"`
}

type MerchantDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.MerchantDetailOutput
}

// MerchantListReq 获取商户列表请求
type MerchantListReq struct {
	g.Meta   `path:"/merchant/list" method:"get" tags:"商户表" summary:"获取商户列表"`
	PageNum  int                 `json:"pageNum" d:"1" dc:"页码"`
	PageSize int                 `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword  string              `json:"keyword" dc:"关键词"`
	TenantID snowflake.JsonInt64 `json:"tenantId" dc:"租户ID"`
	Code     string              `json:"code" dc:"商户编码"`
	Status   *int                `json:"status" dc:"状态"`
}

type MerchantListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.MerchantListOutput `json:"list" dc:"列表数据"`
	Total  int                         `json:"total" dc:"总数"`
}
