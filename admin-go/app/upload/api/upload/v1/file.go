package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// File API

// FileCreateReq 创建文件记录请求
type FileCreateReq struct {
	g.Meta  `path:"/file/create" method:"post" tags:"文件记录" summary:"创建文件记录"`
	DirID   snowflake.JsonInt64 `json:"dirID"  dc:"所属目录"`
	Name    string              `json:"name" v:"required#文件名称不能为空" dc:"文件名称"`
	URL     string              `json:"url" v:"required#文件地址不能为空" dc:"文件地址"`
	Ext     string              `json:"ext"  dc:"文件扩展名"`
	Size    int64               `json:"size"  dc:"文件大小"`
	Mime    string              `json:"mime"  dc:"MIME类型"`
	Storage int                 `json:"storage"  dc:"存储类型"`
	IsImage int                 `json:"isImage"  dc:"是否图片"`
}

// FileCreateRes 创建文件记录响应
type FileCreateRes struct {
	g.Meta `mime:"application/json"`
}

// FileUpdateReq 更新文件记录请求
type FileUpdateReq struct {
	g.Meta  `path:"/file/update" method:"put" tags:"文件记录" summary:"更新文件记录"`
	ID      snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件记录ID"`
	DirID   snowflake.JsonInt64 `json:"dirID" dc:"所属目录"`
	Name    string              `json:"name" dc:"文件名称"`
	URL     string              `json:"url" dc:"文件地址"`
	Ext     string              `json:"ext" dc:"文件扩展名"`
	Size    int64               `json:"size" dc:"文件大小"`
	Mime    string              `json:"mime" dc:"MIME类型"`
	Storage int                 `json:"storage" dc:"存储类型"`
	IsImage int                 `json:"isImage" dc:"是否图片"`
}

// FileUpdateRes 更新文件记录响应
type FileUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// FileDeleteReq 删除文件记录请求
type FileDeleteReq struct {
	g.Meta `path:"/file/delete" method:"delete" tags:"文件记录" summary:"删除文件记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件记录ID"`
}

// FileDeleteRes 删除文件记录响应
type FileDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// FileBatchDeleteReq 批量删除文件记录请求
type FileBatchDeleteReq struct {
	g.Meta `path:"/file/batch-delete" method:"delete" tags:"文件记录" summary:"批量删除文件记录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"文件记录ID列表"`
}

// FileBatchDeleteRes 批量删除文件记录响应
type FileBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// FileDetailReq 获取文件记录详情请求
type FileDetailReq struct {
	g.Meta `path:"/file/detail" method:"get" tags:"文件记录" summary:"获取文件记录详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"文件记录ID"`
}

// FileDetailRes 获取文件记录详情响应
type FileDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.FileDetailOutput
}

// FileListReq 获取文件记录列表请求
type FileListReq struct {
	g.Meta   `path:"/file/list" method:"get" tags:"文件记录" summary:"获取文件记录列表"`
	PageNum  int                 `json:"pageNum" d:"1" dc:"页码"`
	PageSize int                 `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword  string              `json:"keyword" dc:"关键词"`
	DirID    snowflake.JsonInt64 `json:"dirID" dc:"所属目录"`
	Name     string              `json:"name" dc:"文件名称"`
	Storage  int                 `json:"storage" dc:"存储类型"`
	IsImage  *int                `json:"isImage" dc:"是否图片"`
}

// FileListRes 获取文件记录列表响应
type FileListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.FileListOutput `json:"list" dc:"列表数据"`
	Total  int                     `json:"total" dc:"总数"`
}
