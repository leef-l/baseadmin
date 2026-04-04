package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Dir DTO 模型

// DirCreateInput 创建文件目录输入
type DirCreateInput struct {
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Name string `json:"name"`
	Path string `json:"path"`
	Sort int `json:"sort"`
	Status int `json:"status"`
}

// DirUpdateInput 更新文件目录输入
type DirUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Name string `json:"name"`
	Path string `json:"path"`
	Sort int `json:"sort"`
	Status int `json:"status"`
}

// DirDetailOutput 文件目录详情输出
type DirDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	DirName string `json:"dirName"`
	Name string `json:"name"`
	Path string `json:"path"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// DirListOutput 文件目录列表输出
type DirListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	DirName string `json:"dirName"`
	Name string `json:"name"`
	Path string `json:"path"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// DirListInput 文件目录列表查询输入
type DirListInput struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	Status int `json:"status"`
}

// DirTreeOutput 文件目录树形输出
type DirTreeOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	DirName string `json:"dirName"`
	Name string `json:"name"`
	Path string `json:"path"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	Children []*DirTreeOutput `json:"children"`
}

