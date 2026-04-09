package frontend

import (
	"gbaseadmin/codegen/generator/util"
	"gbaseadmin/codegen/parser"
)

var mappings = []util.TemplateMapping{
	{TplFile: "types.tpl", OutputPath: "api/{app}/{module}/types.ts"},
	{TplFile: "api.tpl", OutputPath: "api/{app}/{module}/index.ts"},
	{TplFile: "list.tpl", OutputPath: "views/{app}/{module}/index.vue"},
	{TplFile: "form.tpl", OutputPath: "views/{app}/{module}/modules/form.vue"},
	{TplFile: "detail-drawer.tpl", OutputPath: "views/{app}/{module}/modules/detail-drawer.vue"},
}

// Config 前端生成器配置
type Config struct {
	TemplateDir string              // 模板目录路径
	OutputDir   string              // 输出根目录
	Force       bool                // 是否强制覆盖
	Cache       *util.TemplateCache // 模板缓存（可选）
}

// Generator 前端代码生成器
type Generator struct {
	config Config
}

// New 创建前端代码生成器实例
func New(cfg Config) *Generator {
	return &Generator{config: cfg}
}

// Plan 为一张表规划所有前端输出文件，不直接落盘。
func (g *Generator) Plan(meta *parser.TableMeta) ([]util.PlannedFile, error) {
	return util.PlanFiles(mappings, g.config.TemplateDir, g.config.OutputDir, meta.AppName, meta.ModuleName, g.config.Force, meta, g.config.Cache)
}

// Generate 为一张表生成所有前端代码
func (g *Generator) Generate(meta *parser.TableMeta) ([]string, error) {
	plans, err := g.Plan(meta)
	if err != nil {
		return nil, err
	}
	return util.CommitPlannedFiles(plans)
}
