package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// TenantCreateReq 创建租户请求
type TenantCreateReq struct {
	g.Meta        `path:"/tenant/create" method:"post" tags:"租户表" summary:"创建租户"`
	Name          string      `json:"name" v:"required#租户名称不能为空" dc:"租户名称"`
	Code          string      `json:"code" v:"required#租户编码不能为空" dc:"租户编码"`
	ContactName   string      `json:"contactName" dc:"联系人"`
	ContactPhone  string      `json:"contactPhone" dc:"联系电话"`
	Domain        string      `json:"domain" dc:"租户域名"`
	ExpireAt      *gtime.Time `json:"expireAt" dc:"到期时间"`
	Status        int         `json:"status" dc:"状态"`
	Remark        string      `json:"remark" dc:"备注"`
	CreateAdmin   int         `json:"createAdmin" dc:"是否同步创建管理员:0=否,1=是"`
	AdminUsername string      `json:"adminUsername" dc:"管理员用户名"`
	AdminPassword string      `json:"adminPassword" dc:"管理员密码"`
	AdminNickname string      `json:"adminNickname" dc:"管理员昵称"`
	AdminEmail    string      `json:"adminEmail" dc:"管理员邮箱"`
}

type TenantCreateRes struct {
	g.Meta `mime:"application/json"`
}

// TenantUpdateReq 更新租户请求
type TenantUpdateReq struct {
	g.Meta       `path:"/tenant/update" method:"put" tags:"租户表" summary:"更新租户"`
	ID           snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"租户ID"`
	Name         string              `json:"name" v:"required#租户名称不能为空" dc:"租户名称"`
	Code         string              `json:"code" v:"required#租户编码不能为空" dc:"租户编码"`
	ContactName  string              `json:"contactName" dc:"联系人"`
	ContactPhone string              `json:"contactPhone" dc:"联系电话"`
	Domain       string              `json:"domain" dc:"租户域名"`
	ExpireAt     *gtime.Time         `json:"expireAt" dc:"到期时间"`
	Status       int                 `json:"status" dc:"状态"`
	Remark       string              `json:"remark" dc:"备注"`
}

type TenantUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// TenantDeleteReq 删除租户请求
type TenantDeleteReq struct {
	g.Meta `path:"/tenant/delete" method:"delete" tags:"租户表" summary:"删除租户"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"租户ID"`
}

type TenantDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// TenantBatchDeleteReq 批量删除租户请求
type TenantBatchDeleteReq struct {
	g.Meta `path:"/tenant/batch-delete" method:"delete" tags:"租户表" summary:"批量删除租户"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"租户ID列表"`
}

type TenantBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// TenantDetailReq 获取租户详情请求
type TenantDetailReq struct {
	g.Meta `path:"/tenant/detail" method:"get" tags:"租户表" summary:"获取租户详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"租户ID"`
}

type TenantDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.TenantDetailOutput
}

// TenantListReq 获取租户列表请求
type TenantListReq struct {
	g.Meta   `path:"/tenant/list" method:"get" tags:"租户表" summary:"获取租户列表"`
	PageNum  int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize int    `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword  string `json:"keyword" dc:"关键词"`
	Code     string `json:"code" dc:"租户编码"`
	Status   *int   `json:"status" dc:"状态"`
}

type TenantListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.TenantListOutput `json:"list" dc:"列表数据"`
	Total  int                       `json:"total" dc:"总数"`
}
