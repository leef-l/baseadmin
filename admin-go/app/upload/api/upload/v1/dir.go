package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Dir API

// DirCreateReq 创建文件目录请求
type DirCreateReq struct {
	g.Meta   `path:"/dir/create" method:"post" tags:"文件目录" summary:"创建文件目录"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"上级目录"`
	Name     string              `json:"name" v:"required#目录名称不能为空" dc:"目录名称"`
	Path     string              `json:"path" v:"required#目录路径不能为空" dc:"目录路径"`
	Sort     int                 `json:"sort"  dc:"排序"`
	Status   int                 `json:"status"  dc:"状态"`
}

// DirCreateRes 创建文件目录响应
type DirCreateRes struct {
	g.Meta `mime:"application/json"`
}

// DirUpdateReq 更新文件目录请求
type DirUpdateReq struct {
	g.Meta   `path:"/dir/update" method:"put" tags:"文件目录" summary:"更新文件目录"`
	ID       snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件目录ID"`
	ParentID snowflake.JsonInt64 `json:"parentID" dc:"上级目录"`
	Name     string              `json:"name" dc:"目录名称"`
	Path     string              `json:"path" dc:"目录路径"`
	Sort     int                 `json:"sort" dc:"排序"`
	Status   int                 `json:"status" dc:"状态"`
}

// DirUpdateRes 更新文件目录响应
type DirUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// DirDeleteReq 删除文件目录请求
type DirDeleteReq struct {
	g.Meta `path:"/dir/delete" method:"delete" tags:"文件目录" summary:"删除文件目录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件目录ID"`
}

// DirDeleteRes 删除文件目录响应
type DirDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DirBatchDeleteReq 批量删除文件目录请求
type DirBatchDeleteReq struct {
	g.Meta `path:"/dir/batch-delete" method:"delete" tags:"文件目录" summary:"批量删除文件目录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"文件目录ID列表"`
}

// DirBatchDeleteRes 批量删除文件目录响应
type DirBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DirDetailReq 获取文件目录详情请求
type DirDetailReq struct {
	g.Meta `path:"/dir/detail" method:"get" tags:"文件目录" summary:"获取文件目录详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件目录ID"`
}

// DirDetailRes 获取文件目录详情响应
type DirDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.DirDetailOutput
}

// DirListReq 获取文件目录列表请求
type DirListReq struct {
	g.Meta   `path:"/dir/list" method:"get" tags:"文件目录" summary:"获取文件目录列表"`
	PageNum  int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize int    `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword  string `json:"keyword" dc:"关键词"`
	Status   *int   `json:"status" dc:"状态"`
}

// DirListRes 获取文件目录列表响应
type DirListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.DirListOutput `json:"list" dc:"列表数据"`
	Total  int                    `json:"total" dc:"总数"`
}

// DirTreeReq 获取文件目录树形结构请求
type DirTreeReq struct {
	g.Meta  `path:"/dir/tree" method:"get" tags:"文件目录" summary:"获取文件目录树形结构"`
	Keyword string `json:"keyword" dc:"关键词"`
	Status  *int   `json:"status" dc:"状态"`
}

// DirTreeRes 获取文件目录树形结构响应
type DirTreeRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.DirTreeOutput `json:"list" dc:"树形数据"`
}
