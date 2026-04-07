package backend

import (
	"gbaseadmin/codegen/generator/util"
	"gbaseadmin/codegen/parser"
)

var mappings = []util.TemplateMapping{
	{TplFile: "api.tpl", OutputPath: "api/{app}/v1/{module}.go"},
	{TplFile: "controller.tpl", OutputPath: "internal/controller/{module}/{module}.go"},
	{TplFile: "logic.tpl", OutputPath: "internal/logic/{module}/{module}.go"},
	{TplFile: "service.tpl", OutputPath: "internal/service/{module}.go"},
	{TplFile: "model.tpl", OutputPath: "internal/model/{module}.go"},
	{TplFile: "consts.tpl", OutputPath: "internal/consts/{module}.go"},
}

// Config 后端生成器配置
type Config struct {
	TemplateDir string              // 模板目录路径
	OutputDir   string              // 输出根目录，如 ./app/system/
	Force       bool                // 是否强制覆盖
	Cache       *util.TemplateCache // 模板缓存（可选）
}

// Generator 后端代码生成器
type Generator struct {
	config Config
}

// New 创建后端代码生成器实例
func New(cfg Config) *Generator {
	return &Generator{config: cfg}
}

// Generate 为一张表生成所有后端代码
func (g *Generator) Generate(meta *parser.TableMeta) ([]string, error) {
	return util.GenerateFiles(mappings, g.config.TemplateDir, g.config.OutputDir, meta.AppName, meta.ModuleName, g.config.Force, meta, g.config.Cache)
}

// GenerateToMemory 生成到内存（用于 dry-run diff 预览）
func (g *Generator) GenerateToMemory(meta *parser.TableMeta) (map[string][]byte, error) {
	return util.GenerateToMemory(mappings, g.config.TemplateDir, g.config.OutputDir, meta.AppName, meta.ModuleName, meta, g.config.Cache)
}
