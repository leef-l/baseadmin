package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// 登录
type AuthLoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"认证" summary:"用户登录"`
	Username string `json:"username" v:"required#用户名不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}

type AuthLoginRes struct {
	Token    string              `json:"token"`
	UserID   snowflake.JsonInt64 `json:"userId"`
	Username string              `json:"username"`
	Nickname string              `json:"nickname"`
	Avatar   string              `json:"avatar"`
}

// 票据登录
type AuthTicketLoginReq struct {
	g.Meta `path:"/auth/ticket-login" method:"post" tags:"认证" summary:"票据登录"`
	Ticket string `json:"ticket" v:"required#票据不能为空"`
}

type AuthTicketLoginRes = AuthLoginRes

// 获取当前用户信息
type AuthInfoReq struct {
	g.Meta `path:"/auth/info" method:"get" tags:"认证" summary:"获取当前用户信息"`
}

type AuthInfoRes struct {
	UserID   snowflake.JsonInt64 `json:"userId"`
	Username string              `json:"username"`
	Nickname string              `json:"nickname"`
	Email    string              `json:"email"`
	Avatar   string              `json:"avatar"`
	DeptID   snowflake.JsonInt64 `json:"deptId"`
	Status   int                 `json:"status"`
	Roles    []string            `json:"roles"`
	Perms    []string            `json:"perms"`
}

// 修改密码
type AuthChangePasswordReq struct {
	g.Meta      `path:"/auth/change-password" method:"post" tags:"认证" summary:"修改密码"`
	OldPassword string `json:"oldPassword" v:"required#旧密码不能为空"`
	NewPassword string `json:"newPassword" v:"required|length:8,64#新密码不能为空|密码长度8-64位"`
}

type AuthChangePasswordRes struct{}

// 获取当前用户菜单（动态路由）
type AuthMenusReq struct {
	g.Meta `path:"/auth/menus" method:"get" tags:"认证" summary:"获取当前用户菜单树"`
}

type AuthMenusRes struct {
	Menus interface{} `json:"menus"`
}
