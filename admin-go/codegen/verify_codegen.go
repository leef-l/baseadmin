//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gbaseadmin/codegen/generator/util"
	"gbaseadmin/codegen/parser"
)

// codegen 全面端到端验证 — 覆盖铁律第6条所有场景
// 运行: cd admin-go/codegen && go run verify_codegen.go

func main() {
	cases := []struct {
		name string
		meta *parser.TableMeta
	}{
		{"demo_category (树形+排序+Switch枚举+Tooltip)", buildCategoryMeta()},
		{"demo_article (外键+所有组件+金额+密码+字典+搜索+验证)", buildArticleMeta()},
		{"demo_tag (最简表+Import+BatchEdit)", buildTagMeta()},
	}

	cwd, _ := os.Getwd()
	tplDir := filepath.Join(cwd, "templates")
	outDir := filepath.Join(cwd, "verify_output")
	os.RemoveAll(outDir)

	backendTpls := []string{"api.tpl", "model.tpl", "controller.tpl", "logic.tpl", "service.tpl", "consts.tpl"}
	frontendTpls := []string{"types.tpl", "api.tpl", "list.tpl", "form.tpl", "detail-drawer.tpl"}

	totalErrors := 0
	totalChecks := 0

	for _, tc := range cases {
		fmt.Printf("\n========== %s ==========\n", tc.name)
		for _, tpl := range backendTpls {
			totalChecks++
			if err := renderAndCheck(tplDir, outDir, "backend/"+tpl, tc.meta); err != nil {
				totalErrors++
			}
		}
		for _, tpl := range frontendTpls {
			totalChecks++
			if err := renderAndCheck(tplDir, outDir, "frontend/"+tpl, tc.meta); err != nil {
				totalErrors++
			}
		}
	}

	fmt.Printf("\n========== 结果 ==========\n")
	fmt.Printf("总检查: %d, 失败: %d\n", totalChecks, totalErrors)
	if totalErrors > 0 {
		os.Exit(1)
	}
	fmt.Println("全部通过!")
}

func renderAndCheck(tplDir, outDir, tplFile string, meta *parser.TableMeta) error {
	tplPath := filepath.Join(tplDir, tplFile)
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(util.SharedFuncMap).ParseFiles(tplPath)
	if err != nil {
		fmt.Printf("  FAIL [parse] %s: %v\n", tplFile, err)
		return err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, meta); err != nil {
		fmt.Printf("  FAIL [render] %s: %v\n", tplFile, err)
		return err
	}
	output := buf.String()
	if strings.TrimSpace(output) == "" {
		fmt.Printf("  WARN [empty] %s\n", tplFile)
	}

	// 内容检查
	errs := checkOutput(tplFile, output, meta)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("  FAIL [check] %s: %s\n", tplFile, e)
		}
		return fmt.Errorf("%d checks failed", len(errs))
	}

	// 写入文件供人工审查
	outPath := filepath.Join(outDir, meta.ModuleName, tplFile+".out")
	os.MkdirAll(filepath.Dir(outPath), 0755)
	os.WriteFile(outPath, buf.Bytes(), 0644)

	fmt.Printf("  OK   %s (%d bytes)\n", tplFile, len(output))
	return nil
}

func checkOutput(tplFile, output string, meta *parser.TableMeta) []string {
	var errs []string
	chk := func(cond bool, msg string) {
		if !cond {
			errs = append(errs, msg)
		}
	}

	isBackend := strings.HasPrefix(tplFile, "backend/")

	if isBackend {
		// Go 代码检查
		chk(!strings.Contains(output, "{{"), "包含未渲染的模板标记 {{")
		chk(strings.Contains(output, "package "), "缺少 package 声明")

		if strings.Contains(tplFile, "logic") {
			chk(strings.Contains(output, "func init()"), "缺少 init() 函数")
			chk(strings.Contains(output, "applyListFilter"), "缺少 applyListFilter")
			chk(strings.Contains(output, "isAllowedOrderField"), "缺少 isAllowedOrderField")
			if meta.HasParentID {
				chk(strings.Contains(output, "collectChildIDs"), "树形表缺少 collectChildIDs")
				chk(strings.Contains(output, "doCollectChildIDs"), "树形表缺少 doCollectChildIDs")
				chk(strings.Contains(output, "Tree("), "树形表缺少 Tree 方法")
			}
			if meta.HasMoney {
				chk(strings.Contains(output, "LockUpdate"), "金额表缺少行锁")
				chk(strings.Contains(output, "Transaction"), "金额表缺少事务")
			}
			if meta.HasPassword {
				chk(strings.Contains(output, "bcrypt"), "密码表缺少 bcrypt")
			}
			if meta.HasImport {
				chk(strings.Contains(output, "Import("), "缺少 Import 方法")
			}
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "BatchUpdate("), "缺少 BatchUpdate 方法")
			}
			if meta.HasCreatedBy || meta.HasDeptID {
				chk(strings.Contains(output, "middleware.ApplyDataScope"), "缺少数据权限过滤")
				chk(strings.Contains(output, `"gbaseadmin/app/`+meta.AppName+`/internal/middleware"`), "缺少 middleware 导入")
			}
			// 检查外键关联字段填充
			for _, f := range meta.Fields {
				if f.RefFieldName != "" && !f.IsHidden {
					chk(strings.Contains(output, "fillRefFields"), "有外键但缺少 fillRefFields")
					break
				}
			}
		}

		if strings.Contains(tplFile, "api") {
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "BatchUpdateReq"), "缺少 BatchUpdateReq")
			}
			if meta.HasParentID {
				chk(strings.Contains(output, "TreeReq"), "树形表缺少 TreeReq")
			}
			if meta.HasImport {
				chk(strings.Contains(output, "ImportReq"), "缺少 ImportReq")
				chk(strings.Contains(output, "ImportTemplateReq"), "缺少 ImportTemplateReq")
			}
		}

		if strings.Contains(tplFile, "model") {
			chk(strings.Contains(output, "CreateInput"), "缺少 CreateInput")
			chk(strings.Contains(output, "UpdateInput"), "缺少 UpdateInput")
			chk(strings.Contains(output, "DetailOutput"), "缺少 DetailOutput")
			chk(strings.Contains(output, "ListOutput"), "缺少 ListOutput")
			chk(strings.Contains(output, "ListInput"), "缺少 ListInput")
			if meta.HasParentID {
				chk(strings.Contains(output, "TreeOutput"), "树形表缺少 TreeOutput")
				chk(strings.Contains(output, "TreeInput"), "树形表缺少 TreeInput")
			}
			// 检查外键字段使用 snowflake.JsonInt64
			for _, f := range meta.Fields {
				if (f.IsForeignKey || f.IsParentID) && !f.IsHidden {
					chk(strings.Contains(output, f.NameCamel+" snowflake.JsonInt64"), f.Name+" 应使用 snowflake.JsonInt64")
					break
				}
			}
		}

		if strings.Contains(tplFile, "controller") {
			chk(strings.Contains(output, "csv.NewWriter"), "Export 应使用 csv.Writer")
			if meta.HasImport {
				chk(strings.Contains(output, "ImportTemplate"), "缺少 ImportTemplate")
			}
		}

		if strings.Contains(tplFile, "consts") {
			for _, f := range meta.Fields {
				if f.IsEnum {
					for _, ev := range f.EnumValues {
						if ev.NameIdent != "" {
							chk(strings.Contains(output, meta.ModelName+f.NameCamel+ev.NameIdent), "缺少枚举常量 "+ev.NameIdent)
						}
					}
				}
			}
		}
	} else {
		// 前端代码检查
		if strings.Contains(tplFile, "list") {
			if meta.HasEnum {
				chk(strings.Contains(output, "Tag"), "有枚举但缺少 Tag 导入")
				chk(strings.Contains(output, "TAG_COLORS"), "有枚举但缺少 TAG_COLORS")
			} else {
				chk(!strings.Contains(output, "TAG_COLORS"), "无枚举但生成了 TAG_COLORS")
			}
			if meta.HasTooltip {
				chk(strings.Contains(output, "Tooltip"), "有 Tooltip 但缺少导入")
				chk(strings.Contains(output, "tooltipHeader"), "有 Tooltip 但缺少 tooltipHeader")
			}
			if meta.HasParentID {
				chk(strings.Contains(output, "treeConfig"), "树形表缺少 treeConfig")
				chk(!strings.Contains(output, "checkbox"), "树形表不应有 checkbox")
			} else {
				chk(strings.Contains(output, "checkbox"), "非树形表应有 checkbox")
			}
			if meta.HasImport {
				chk(strings.Contains(output, "handleImport"), "缺少 handleImport")
				chk(strings.Contains(output, "handleDownloadTemplate"), "缺少 handleDownloadTemplate")
			}
		}

		if strings.Contains(tplFile, "form") {
			if meta.HasTooltip {
				chk(strings.Contains(output, "tooltipLabel"), "有 Tooltip 但缺少 tooltipLabel")
			}
			if meta.HasParentID {
				chk(strings.Contains(output, "TreeSelect"), "树形表 form 缺少 TreeSelect")
				chk(strings.Contains(output, "treeData"), "树形表 form 缺少 treeData")
				chk(strings.Contains(output, "顶级节点"), "树形表 form 缺少顶级节点")
			}
			if meta.HasDict {
				chk(strings.Contains(output, "getDictByType"), "有字典但缺少 getDictByType")
			}
			// 检查外键 API 导入路径
			for _, f := range meta.Fields {
				if f.IsForeignKey && !f.IsHidden && f.RefTable != "" {
					expectedImport := fmt.Sprintf("#/api/%s/%s", f.RefTableApp, f.RefTable)
					chk(strings.Contains(output, expectedImport), f.Name+" 外键 API 路径应为 "+expectedImport)
				}
			}
			// 检查密码字段特殊处理
			if meta.HasPassword {
				chk(strings.Contains(output, "InputPassword"), "密码字段应使用 InputPassword")
				chk(strings.Contains(output, "不填则不修改"), "密码字段编辑时应提示不填则不修改")
			}
		}

		if strings.Contains(tplFile, "detail-drawer") {
			if meta.HasEnum {
				chk(strings.Contains(output, "Tag"), "有枚举但 detail-drawer 缺少 Tag")
			}
			// 检查图片字段
			if meta.HasImage {
				chk(strings.Contains(output, "<img"), "有图片字段但 detail-drawer 缺少 img 标签")
			}
		}

		if strings.Contains(tplFile, "types") {
			chk(strings.Contains(output, "interface "+meta.ModelName+"Item"), "缺少 Item 接口")
			chk(strings.Contains(output, "interface "+meta.ModelName+"ListParams"), "缺少 ListParams 接口")
			chk(strings.Contains(output, "interface "+meta.ModelName+"CreateParams"), "缺少 CreateParams 接口")
			chk(strings.Contains(output, "interface "+meta.ModelName+"UpdateParams"), "缺少 UpdateParams 接口")
			if meta.HasParentID {
				chk(strings.Contains(output, "children?"), "树形表 Item 缺少 children")
				chk(strings.Contains(output, "TreeParams"), "树形表缺少 TreeParams")
			}
		}

		if strings.HasSuffix(tplFile, "frontend/api.tpl") {
			chk(strings.Contains(output, "get"+meta.ModelName+"Detail"), "缺少 getDetail")
			chk(strings.Contains(output, "create"+meta.ModelName), "缺少 create")
			chk(strings.Contains(output, "update"+meta.ModelName), "缺少 update")
			chk(strings.Contains(output, "delete"+meta.ModelName), "缺少 delete")
			chk(strings.Contains(output, "export"+meta.ModelName), "缺少 export")
			if meta.HasParentID {
				chk(strings.Contains(output, "get"+meta.ModelName+"Tree"), "树形表缺少 getTree")
			} else {
				chk(strings.Contains(output, "get"+meta.ModelName+"List"), "非树形表缺少 getList")
			}
		}
	}

	return errs
}

// PLACEHOLDER_META_BUILDERS

// --- 表1: demo_category（树形+排序+Switch枚举+Tooltip）---
func buildCategoryMeta() *parser.TableMeta {
	m := &parser.TableMeta{
		TableName: "demo_category", AppName: "demo", AppNameCamel: "Demo",
		ModelName: "Category", DaoName: "DemoCategory", ModuleName: "category", PackageName: "category",
		Comment: "分类", HasParentID: true, HasStatus: true, HasSort: true,
		HasTooltip: true, HasSearchable: true, HasCreatedBy: true, HasDeptID: true,
		HasEnum: true, HasBatchEdit: true, HasImport: false, HasTreeSelect: true,
		ParentDisplayField: "name",
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		f("parent_id", "ParentID", "ParentId", "parentId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsParentID = true; f.Component = "TreeSelectSingle"
			f.RefTable = "category"; f.RefTableDB = "demo_category"; f.RefTableApp = "demo"
			f.RefTableCamel = "Category"; f.RefTableLower = "category"
			f.RefDisplayField = "name"; f.RefDisplayCamel = "Name"; f.RefDisplayLower = "name"
			f.RefFieldName = "CategoryName"; f.RefFieldJSON = "categoryName"; f.RefIsTree = true
		}),
		f("name", "Name", "Name", "name", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true; f.IsSearchable = true; f.MaxLength = 50; f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:50"}; f.UpdateValidationRules = []string{"max-length:50"}
		}),
		f("icon", "Icon", "Icon", "icon", "string", "string", func(f *parser.FieldMeta) { f.Component = "IconPicker" }),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.TooltipText = "升序"; f.Component = "InputNumber" }),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true; f.Component = "Switch"; f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "禁用", NameIdent: "Disabled"}, {Value: "1", Label: "启用", NameIdent: "Enabled"}}
		}),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return m
}

// --- 表2: demo_article（复杂表：外键+所有组件+金额+密码+字典+搜索+验证规则）---
func buildArticleMeta() *parser.TableMeta {
	m := &parser.TableMeta{
		TableName: "demo_article", AppName: "demo", AppNameCamel: "Demo",
		ModelName: "Article", DaoName: "DemoArticle", ModuleName: "article", PackageName: "article",
		Comment: "文章", HasStatus: true, HasSort: true, HasPassword: true,
		HasTooltip: true, HasRichText: true, HasMoney: true, HasSearchable: true,
		HasCreatedBy: true, HasDeptID: true, HasEnum: true, HasImage: true,
		HasBatchEdit: true, HasImport: true, HasForeignKey: true, HasDict: true,
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		// 同应用树形外键
		f("category_id", "CategoryID", "CategoryId", "categoryId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsForeignKey = true; f.Component = "Select"; f.IsRequired = true
			f.ValidationRules = []string{"required"}
			f.RefTable = "category"; f.RefTableDB = "demo_category"; f.RefTableApp = "demo"
			f.RefTableCamel = "Category"; f.RefTableLower = "category"
			f.RefDisplayField = "name"; f.RefDisplayCamel = "Name"; f.RefDisplayLower = "name"
			f.RefFieldName = "CategoryName"; f.RefFieldJSON = "categoryName"; f.RefIsTree = true
		}),
		// 跨应用普通外键
		f("user_id", "UserID", "UserId", "userId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsForeignKey = true; f.Component = "Select"; f.IsRequired = true
			f.ValidationRules = []string{"required"}
			f.RefTable = "users"; f.RefTableDB = "system_users"; f.RefTableApp = "system"
			f.RefTableCamel = "Users"; f.RefTableLower = "users"
			f.RefDisplayField = "username"; f.RefDisplayCamel = "Username"; f.RefDisplayLower = "username"
			f.RefFieldName = "UsersUsername"; f.RefFieldJSON = "usersUsername"; f.RefIsTree = false
		}),
		f("title", "Title", "Title", "title", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true; f.IsSearchable = true; f.MaxLength = 200; f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:200"}; f.UpdateValidationRules = []string{"max-length:200"}
		}),
		f("order_no", "OrderNo", "OrderNo", "orderNo", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true; f.IsSearchable = true; f.IsExactSearch = true; f.MaxLength = 50; f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:50"}; f.UpdateValidationRules = []string{"max-length:50"}
		}),
		f("cover", "Cover", "Cover", "cover", "string", "string", func(f *parser.FieldMeta) { f.Component = "ImageUpload" }),
		f("attachment_file", "AttachmentFile", "AttachmentFile", "attachmentFile", "string", "string", func(f *parser.FieldMeta) { f.Component = "FileUpload" }),
		f("body_content", "BodyContent", "BodyContent", "bodyContent", "string", "string", func(f *parser.FieldMeta) { f.DBType = "text"; f.Component = "RichText" }),
		f("extra_json", "ExtraJSON", "ExtraJson", "extraJSON", "string", "string", func(f *parser.FieldMeta) { f.DBType = "text"; f.Component = "JsonEditor" }),
		f("link_url", "LinkURL", "LinkUrl", "linkURL", "string", "string", func(f *parser.FieldMeta) {
			f.Component = "InputUrl"; f.FrontendRules = "url"; f.MaxLength = 500
			f.ValidationRules = []string{"url", "max-length:500"}; f.UpdateValidationRules = []string{"url", "max-length:500"}
		}),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true; f.Component = "Radio"; f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "草稿", NameIdent: "Draft"}, {Value: "1", Label: "已发布", NameIdent: "Published"}, {Value: "2", Label: "已下架", NameIdent: "Offline"}}
		}),
		f("type", "Type", "Type", "type", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true; f.Component = "Select"; f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "1", Label: "普通", NameIdent: "Regular"}, {Value: "2", Label: "置顶", NameIdent: "Pinned"}, {Value: "3", Label: "推荐", NameIdent: "Recommended"}, {Value: "4", Label: "热门", NameIdent: "Hot"}}
		}),
		f("is_top", "IsTop", "IsTop", "isTop", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true; f.Component = "Switch"; f.DefaultValue = "0"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "否", NameIdent: "No"}, {Value: "1", Label: "是", NameIdent: "Yes"}}
		}),
		f("price", "Price", "Price", "price", "int", "number", func(f *parser.FieldMeta) { f.IsMoney = true; f.Component = "InputNumber" }),
		f("pay_password", "PayPassword", "PayPassword", "payPassword", "string", "string", func(f *parser.FieldMeta) {
			f.IsPassword = true; f.Component = "Password"
			f.ValidationRules = []string{"length:6,32"}; f.UpdateValidationRules = []string{"length:6,32"}
		}),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.TooltipText = "升序"; f.Component = "InputNumber" }),
		f("icon", "Icon", "Icon", "icon", "string", "string", func(f *parser.FieldMeta) { f.Component = "IconPicker" }),
		f("email", "Email", "Email", "email", "string", "string", func(f *parser.FieldMeta) {
			f.Component = "Input"; f.FrontendRules = "email"; f.MaxLength = 100
			f.ValidationRules = []string{"email", "max-length:100"}; f.UpdateValidationRules = []string{"email", "max-length:100"}
		}),
		f("phone", "Phone", "Phone", "phone", "string", "string", func(f *parser.FieldMeta) {
			f.IsSearchable = true; f.Component = "Input"; f.FrontendRules = "phone"; f.MaxLength = 20
			f.ValidationRules = []string{"phone-loose", "max-length:20"}; f.UpdateValidationRules = []string{"phone-loose", "max-length:20"}
		}),
		f("remark", "Remark", "Remark", "remark", "string", "string", func(f *parser.FieldMeta) { f.DBType = "text"; f.IsSearchable = true; f.Component = "Textarea" }),
		f("level", "Level", "Level", "level", "string", "string", func(f *parser.FieldMeta) { f.DictType = "article_level"; f.Component = "Select" }),
		f("extra_field", "ExtraField", "ExtraField", "extraField", "string", "string", func(f *parser.FieldMeta) { f.Component = "Input" }),
		f("publish_at", "PublishAt", "PublishAt", "publishAt", "*gtime.Time", "string", func(f *parser.FieldMeta) { f.IsTimeField = true; f.Component = "DateTimePicker" }),
		f("expire_at", "ExpireAt", "ExpireAt", "expireAt", "*gtime.Time", "string", func(f *parser.FieldMeta) { f.IsTimeField = true; f.Component = "DateTimePicker" }),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return m
}

// --- 表3: demo_tag（最简表）---
func buildTagMeta() *parser.TableMeta {
	m := &parser.TableMeta{
		TableName: "demo_tag", AppName: "demo", AppNameCamel: "Demo",
		ModelName: "Tag", DaoName: "DemoTag", ModuleName: "tag", PackageName: "tag",
		Comment: "标签", HasStatus: true, HasSort: true, HasSearchable: true,
		HasCreatedBy: true, HasDeptID: true, HasEnum: true, HasBatchEdit: true, HasImport: true,
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		f("name", "Name", "Name", "name", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true; f.IsSearchable = true; f.MaxLength = 50; f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:50"}; f.UpdateValidationRules = []string{"max-length:50"}
		}),
		f("color", "Color", "Color", "color", "string", "string", func(f *parser.FieldMeta) { f.Component = "Input"; f.MaxLength = 20 }),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.Component = "InputNumber" }),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true; f.Component = "Switch"; f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "禁用", NameIdent: "Disabled"}, {Value: "1", Label: "启用", NameIdent: "Enabled"}}
		}),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return m
}

// --- 辅助函数 ---
func f(name, camel, dao, lower, goType, tsType string, customize func(*parser.FieldMeta)) parser.FieldMeta {
	fm := parser.FieldMeta{
		Name: name, NameCamel: camel, NameDao: dao, NameLower: lower,
		GoType: goType, TSType: tsType, Label: camel, ShortLabel: camel,
	}
	if customize != nil {
		customize(&fm)
	}
	return fm
}

func hiddenField(name, camel, dao, lower, goType, tsType string) parser.FieldMeta {
	return f(name, camel, dao, lower, goType, tsType, func(f *parser.FieldMeta) { f.IsHidden = true })
}

func hiddenTimeField(name, camel string) parser.FieldMeta {
	dao := strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), " ", "")
	lower := strings.ToLower(camel[:1]) + camel[1:]
	return f(name, camel, dao, lower, "*gtime.Time", "string", func(f *parser.FieldMeta) { f.IsHidden = true; f.IsTimeField = true })
}
