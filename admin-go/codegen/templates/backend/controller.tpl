package {{.PackageName}}

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/{{.AppName}}/api/{{.AppName}}/v1"
	"gbaseadmin/app/{{.AppName}}/internal/model"
	"gbaseadmin/app/{{.AppName}}/internal/service"
)

func csvSafe{{.ModelName}}(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var {{.ModelName}} = c{{.ModelName}}{}

type c{{.ModelName}} struct{}

// Create 创建{{.Comment}}
func (c *c{{.ModelName}}) Create(ctx context.Context, req *v1.{{.ModelName}}CreateReq) (res *v1.{{.ModelName}}CreateRes, err error) {
	err = service.{{.ModelName}}().Create(ctx, &model.{{.ModelName}}CreateInput{
{{- range .Fields}}
{{- if and (not .IsID) (not .IsHidden)}}
		{{.NameCamel}}: req.{{.NameCamel}},
{{- end}}
{{- end}}
	})
	return
}

// Update 更新{{.Comment}}
func (c *c{{.ModelName}}) Update(ctx context.Context, req *v1.{{.ModelName}}UpdateReq) (res *v1.{{.ModelName}}UpdateRes, err error) {
	err = service.{{.ModelName}}().Update(ctx, &model.{{.ModelName}}UpdateInput{
		ID: req.ID,
{{- range .Fields}}
{{- if and (not .IsID) (not .IsHidden)}}
		{{.NameCamel}}: req.{{.NameCamel}},
{{- end}}
{{- end}}
	})
	return
}

// Delete 删除{{.Comment}}
func (c *c{{.ModelName}}) Delete(ctx context.Context, req *v1.{{.ModelName}}DeleteReq) (res *v1.{{.ModelName}}DeleteRes, err error) {
	err = service.{{.ModelName}}().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除{{.Comment}}
func (c *c{{.ModelName}}) BatchDelete(ctx context.Context, req *v1.{{.ModelName}}BatchDeleteReq) (res *v1.{{.ModelName}}BatchDeleteRes, err error) {
	err = service.{{.ModelName}}().BatchDelete(ctx, req.IDs)
	return
}
{{- if .HasBatchEdit}}

// BatchUpdate 批量编辑{{.Comment}}
func (c *c{{.ModelName}}) BatchUpdate(ctx context.Context, req *v1.{{.ModelName}}BatchUpdateReq) (res *v1.{{.ModelName}}BatchUpdateRes, err error) {
	err = service.{{.ModelName}}().BatchUpdate(ctx, &model.{{.ModelName}}BatchUpdateInput{
		IDs: req.IDs,
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (.IsEnum)}}
		{{.NameCamel}}: req.{{.NameCamel}},
{{- end}}
{{- end}}
	})
	return
}
{{- end}}

// Detail 获取{{.Comment}}详情
func (c *c{{.ModelName}}) Detail(ctx context.Context, req *v1.{{.ModelName}}DetailReq) (res *v1.{{.ModelName}}DetailRes, err error) {
	res = &v1.{{.ModelName}}DetailRes{}
	res.{{.ModelName}}DetailOutput, err = service.{{.ModelName}}().Detail(ctx, req.ID)
	return
}

// List 获取{{.Comment}}列表
func (c *c{{.ModelName}}) List(ctx context.Context, req *v1.{{.ModelName}}ListReq) (res *v1.{{.ModelName}}ListRes, err error) {
	res = &v1.{{.ModelName}}ListRes{}
	res.List, res.Total, err = service.{{.ModelName}}().List(ctx, &model.{{.ModelName}}ListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
{{- if .HasKeywordSearch}}
		Keyword: req.Keyword,
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
		{{.NameCamel}}Start: req.{{.NameCamel}}Start,
		{{.NameCamel}}End: req.{{.NameCamel}}End,
{{- else}}
		{{.NameCamel}}: req.{{.NameCamel}},
{{- end}}
{{- end}}
	})
	return
}
// Export 导出{{.Comment}}
func (c *c{{.ModelName}}) Export(ctx context.Context, req *v1.{{.ModelName}}ExportReq) (res *v1.{{.ModelName}}ExportRes, err error) {
	list, err := service.{{.ModelName}}().Export(ctx, &model.{{.ModelName}}ListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
{{- if .HasKeywordSearch}}
		Keyword: req.Keyword,
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
		{{.NameCamel}}Start: req.{{.NameCamel}}Start,
		{{.NameCamel}}End: req.{{.NameCamel}}End,
{{- else}}
		{{.NameCamel}}: req.{{.NameCamel}},
{{- end}}
{{- end}}
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="{{.ModuleName}}.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头
	_ = w.Write([]string{ {{- $first := true}}{{- range .Fields}}{{- if and (not .IsHidden) (not .IsID) (not .IsPassword)}}{{if not $first}}, {{end}}"{{.ShortLabel}}"{{$first = false}}{{- end}}{{- end}}, "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsPassword)}}
{{- if .RefFieldJSON}}
			csvSafe{{$.ModelName}}(item.{{.RefFieldName}}),
{{- else if eq .GoType "*gtime.Time"}}
			func() string { if item.{{.NameCamel}} != nil { return item.{{.NameCamel}}.String() }; return "" }(),
{{- else if .IsMoney}}
			fmt.Sprintf("%.2f", float64(item.{{.NameCamel}})/100),
{{- else if eq .GoType "string"}}
			csvSafe{{$.ModelName}}(item.{{.NameCamel}}),
{{- else}}
			fmt.Sprintf("%v", item.{{.NameCamel}}),
{{- end}}
{{- end}}
{{- end}}
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}
{{- if .HasImport}}

// Import 导入{{.Comment}}
func (c *c{{.ModelName}}) Import(ctx context.Context, req *v1.{{.ModelName}}ImportReq) (res *v1.{{.ModelName}}ImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.{{.ModelName}}().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.{{.ModelName}}ImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载{{.Comment}}导入模板
func (c *c{{.ModelName}}) ImportTemplate(ctx context.Context, req *v1.{{.ModelName}}ImportTemplateReq) (res *v1.{{.ModelName}}ImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="{{.ModuleName}}_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{ {{- $first := true}}{{- range .Fields}}{{- if and (not .IsHidden) (not .IsID) (not .IsPassword) (not .IsTimeField) (or (not $.HasTenantScope) (and (ne .Name "tenant_id") (ne .Name "merchant_id")))}}{{if not $first}}, {{end}}"{{.ShortLabel}}"{{$first = false}}{{- end}}{{- end}}})
	w.Flush()
	return
}
{{- end}}
{{- if .HasParentID}}

// Tree 获取{{.Comment}}树形结构
func (c *c{{.ModelName}}) Tree(ctx context.Context, req *v1.{{.ModelName}}TreeReq) (res *v1.{{.ModelName}}TreeRes, err error) {
	res = &v1.{{.ModelName}}TreeRes{}
	res.List, err = service.{{.ModelName}}().Tree(ctx, &model.{{.ModelName}}TreeInput{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
{{- if .HasKeywordSearch}}
		Keyword: req.Keyword,
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
		{{.NameCamel}}Start: req.{{.NameCamel}}Start,
		{{.NameCamel}}End: req.{{.NameCamel}}End,
{{- else}}
		{{.NameCamel}}: req.{{.NameCamel}},
{{- end}}
{{- end}}
	})
	return
}
{{end}}
