package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// User API

// UserCreateReq 创建会员用户请求
type UserCreateReq struct {
	g.Meta `path:"/user/create" method:"post" tags:"会员用户" summary:"创建会员用户"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"上级会员"`
	Username string `json:"username" v:"required|max-length:50" dc:"用户名（登录账号）"`
	Password string `json:"password" v:"length:6,32" dc:"密码（bcrypt加密）"`
	Nickname string `json:"nickname" v:"max-length:50" dc:"昵称"`
	Phone string `json:"phone" v:"phone-loose|max-length:20" dc:"手机号"`
	Avatar string `json:"avatar" v:"max-length:500" dc:"头像"`
	RealName string `json:"realName" v:"max-length:50" dc:"真实姓名"`
	LevelID snowflake.JsonInt64 `json:"levelID"  dc:"当前等级"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"  dc:"等级到期时间"`
	TeamCount int `json:"teamCount"  dc:"团队总人数"`
	DirectCount int `json:"directCount"  dc:"直推人数"`
	ActiveCount int `json:"activeCount"  dc:"有效用户数"`
	TeamTurnover int64 `json:"teamTurnover"  dc:"团队总营业额（分）"`
	IsActive int `json:"isActive"  dc:"是否激活"`
	IsQualified int `json:"isQualified"  dc:"仓库资格"`
	InviteCode string `json:"inviteCode" v:"max-length:32" dc:"邀请码"`
	RegisterIP string `json:"registerIP" v:"max-length:45" dc:"注册IP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"  dc:"最后登录时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// UserCreateRes 创建会员用户响应
type UserCreateRes struct {
	g.Meta `mime:"application/json"`
}

// UserUpdateReq 更新会员用户请求
type UserUpdateReq struct {
	g.Meta `path:"/user/update" method:"put" tags:"会员用户" summary:"更新会员用户"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员用户ID"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"上级会员"`
	Username string `json:"username" v:"max-length:50" dc:"用户名（登录账号）"`
	Password string `json:"password" v:"length:6,32" dc:"密码（bcrypt加密）"`
	Nickname string `json:"nickname" v:"max-length:50" dc:"昵称"`
	Phone string `json:"phone" v:"phone-loose|max-length:20" dc:"手机号"`
	Avatar string `json:"avatar" v:"max-length:500" dc:"头像"`
	RealName string `json:"realName" v:"max-length:50" dc:"真实姓名"`
	LevelID snowflake.JsonInt64 `json:"levelID"  dc:"当前等级"`
	LevelExpireAt *gtime.Time `json:"levelExpireAt"  dc:"等级到期时间"`
	TeamCount int `json:"teamCount"  dc:"团队总人数"`
	DirectCount int `json:"directCount"  dc:"直推人数"`
	ActiveCount int `json:"activeCount"  dc:"有效用户数"`
	TeamTurnover int64 `json:"teamTurnover"  dc:"团队总营业额（分）"`
	IsActive int `json:"isActive"  dc:"是否激活"`
	IsQualified int `json:"isQualified"  dc:"仓库资格"`
	InviteCode string `json:"inviteCode" v:"max-length:32" dc:"邀请码"`
	RegisterIP string `json:"registerIP" v:"max-length:45" dc:"注册IP"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"  dc:"最后登录时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// UserUpdateRes 更新会员用户响应
type UserUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// UserDeleteReq 删除会员用户请求
type UserDeleteReq struct {
	g.Meta `path:"/user/delete" method:"delete" tags:"会员用户" summary:"删除会员用户"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员用户ID"`
}

// UserDeleteRes 删除会员用户响应
type UserDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// UserBatchDeleteReq 批量删除会员用户请求
type UserBatchDeleteReq struct {
	g.Meta `path:"/user/batch-delete" method:"delete" tags:"会员用户" summary:"批量删除会员用户"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员用户ID列表"`
}

// UserBatchDeleteRes 批量删除会员用户响应
type UserBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// UserBatchUpdateReq 批量编辑会员用户请求
type UserBatchUpdateReq struct {
	g.Meta `path:"/user/batch-update" method:"put" tags:"会员用户" summary:"批量编辑会员用户"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"会员用户ID列表"`
	IsActive *int `json:"isActive" dc:"是否激活"`
	IsQualified *int `json:"isQualified" dc:"仓库资格"`
	Status *int `json:"status" dc:"状态"`
}

// UserBatchUpdateRes 批量编辑会员用户响应
type UserBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// UserDetailReq 获取会员用户详情请求
type UserDetailReq struct {
	g.Meta `path:"/user/detail" method:"get" tags:"会员用户" summary:"获取会员用户详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"会员用户ID"`
}

// UserDetailRes 获取会员用户详情响应
type UserDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.UserDetailOutput
}

// UserListReq 获取会员用户列表请求
type UserListReq struct {
	g.Meta    `path:"/user/list" method:"get" tags:"会员用户" summary:"获取会员用户列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Username string `json:"username" dc:"用户名（登录账号）"`
	InviteCode string `json:"inviteCode" dc:"邀请码"`
	Nickname string `json:"nickname" dc:"昵称"`
	RealName string `json:"realName" dc:"真实姓名"`
	Phone string `json:"phone" dc:"手机号"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"上级会员"`
	LevelID *snowflake.JsonInt64 `json:"levelID" dc:"当前等级"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsActive *int `json:"isActive" dc:"是否激活"`
	IsQualified *int `json:"isQualified" dc:"仓库资格"`
	Status *int `json:"status" dc:"状态"`
	LevelExpireAtStart string `json:"levelExpireAtStart" dc:"等级到期时间开始时间"`
	LevelExpireAtEnd string `json:"levelExpireAtEnd" dc:"等级到期时间结束时间"`
	LastLoginAtStart string `json:"lastLoginAtStart" dc:"最后登录时间开始时间"`
	LastLoginAtEnd string `json:"lastLoginAtEnd" dc:"最后登录时间结束时间"`
}

// UserListRes 获取会员用户列表响应
type UserListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.UserListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// UserExportReq 导出会员用户请求
type UserExportReq struct {
	g.Meta    `path:"/user/export" method:"get" tags:"会员用户" summary:"导出会员用户"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Username string `json:"username" dc:"用户名（登录账号）"`
	InviteCode string `json:"inviteCode" dc:"邀请码"`
	Nickname string `json:"nickname" dc:"昵称"`
	RealName string `json:"realName" dc:"真实姓名"`
	Phone string `json:"phone" dc:"手机号"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"上级会员"`
	LevelID *snowflake.JsonInt64 `json:"levelID" dc:"当前等级"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsActive *int `json:"isActive" dc:"是否激活"`
	IsQualified *int `json:"isQualified" dc:"仓库资格"`
	Status *int `json:"status" dc:"状态"`
	LevelExpireAtStart string `json:"levelExpireAtStart" dc:"等级到期时间开始时间"`
	LevelExpireAtEnd string `json:"levelExpireAtEnd" dc:"等级到期时间结束时间"`
	LastLoginAtStart string `json:"lastLoginAtStart" dc:"最后登录时间开始时间"`
	LastLoginAtEnd string `json:"lastLoginAtEnd" dc:"最后登录时间结束时间"`
}

// UserExportRes 导出会员用户响应
type UserExportRes struct {
	g.Meta `mime:"text/csv"`
}

// UserTreeReq 获取会员用户树形结构请求
type UserTreeReq struct {
	g.Meta    `path:"/user/tree" method:"get" tags:"会员用户" summary:"获取会员用户树形结构"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Username string `json:"username" dc:"用户名（登录账号）"`
	InviteCode string `json:"inviteCode" dc:"邀请码"`
	Nickname string `json:"nickname" dc:"昵称"`
	RealName string `json:"realName" dc:"真实姓名"`
	Phone string `json:"phone" dc:"手机号"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"上级会员"`
	LevelID *snowflake.JsonInt64 `json:"levelID" dc:"当前等级"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsActive *int `json:"isActive" dc:"是否激活"`
	IsQualified *int `json:"isQualified" dc:"仓库资格"`
	Status *int `json:"status" dc:"状态"`
	LevelExpireAtStart string `json:"levelExpireAtStart" dc:"等级到期时间开始时间"`
	LevelExpireAtEnd string `json:"levelExpireAtEnd" dc:"等级到期时间结束时间"`
	LastLoginAtStart string `json:"lastLoginAtStart" dc:"最后登录时间开始时间"`
	LastLoginAtEnd string `json:"lastLoginAtEnd" dc:"最后登录时间结束时间"`
}

// UserTreeRes 获取会员用户树形结构响应
type UserTreeRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.UserTreeOutput `json:"list" dc:"树形数据"`
}
