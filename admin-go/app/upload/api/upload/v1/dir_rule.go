package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// DirRule API

// DirRuleCreateReq 创建文件目录规则请求
type DirRuleCreateReq struct {
	g.Meta   `path:"/dir_rule/create" method:"post" tags:"文件目录规则" summary:"创建文件目录规则"`
	DirID    snowflake.JsonInt64 `json:"dirID" v:"required#目录ID不能为空" dc:"目录ID"`
	Category int                 `json:"category"  dc:"类别"`
	SavePath string              `json:"savePath"  dc:"保存目录"`
	Status   int                 `json:"status"  dc:"状态"`
}

// DirRuleCreateRes 创建文件目录规则响应
type DirRuleCreateRes struct {
	g.Meta `mime:"application/json"`
}

// DirRuleUpdateReq 更新文件目录规则请求
type DirRuleUpdateReq struct {
	g.Meta   `path:"/dir_rule/update" method:"put" tags:"文件目录规则" summary:"更新文件目录规则"`
	ID       snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件目录规则ID"`
	DirID    snowflake.JsonInt64 `json:"dirID" dc:"目录ID"`
	Category int                 `json:"category" dc:"类别"`
	SavePath string              `json:"savePath" dc:"保存目录"`
	Status   int                 `json:"status" dc:"状态"`
}

// DirRuleUpdateRes 更新文件目录规则响应
type DirRuleUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// DirRuleDeleteReq 删除文件目录规则请求
type DirRuleDeleteReq struct {
	g.Meta `path:"/dir_rule/delete" method:"delete" tags:"文件目录规则" summary:"删除文件目录规则"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件目录规则ID"`
}

// DirRuleDeleteRes 删除文件目录规则响应
type DirRuleDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DirRuleBatchDeleteReq 批量删除文件目录规则请求
type DirRuleBatchDeleteReq struct {
	g.Meta `path:"/dir_rule/batch-delete" method:"delete" tags:"文件目录规则" summary:"批量删除文件目录规则"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"文件目录规则ID列表"`
}

// DirRuleBatchDeleteRes 批量删除文件目录规则响应
type DirRuleBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DirRuleDetailReq 获取文件目录规则详情请求
type DirRuleDetailReq struct {
	g.Meta `path:"/dir_rule/detail" method:"get" tags:"文件目录规则" summary:"获取文件目录规则详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件目录规则ID"`
}

// DirRuleDetailRes 获取文件目录规则详情响应
type DirRuleDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.DirRuleDetailOutput
}

// DirRuleListReq 获取文件目录规则列表请求
type DirRuleListReq struct {
	g.Meta   `path:"/dir_rule/list" method:"get" tags:"文件目录规则" summary:"获取文件目录规则列表"`
	PageNum  int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize int    `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword  string `json:"keyword" dc:"关键词"`
	Category int    `json:"category" dc:"类别"`
	Status   *int   `json:"status" dc:"状态"`
}

// DirRuleListRes 获取文件目录规则列表响应
type DirRuleListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.DirRuleListOutput `json:"list" dc:"列表数据"`
	Total  int                        `json:"total" dc:"总数"`
}
