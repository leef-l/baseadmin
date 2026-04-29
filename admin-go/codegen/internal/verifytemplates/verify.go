package verifytemplates

import (
	"bytes"
	"fmt"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gbaseadmin/codegen/generator/util"
	"gbaseadmin/codegen/parser"
)

type TemplateCase struct {
	Key  string
	Name string
	Meta *parser.TableMeta
}

type RenderedTemplate struct {
	CaseKey      string
	TemplateFile string
	Output       string
}

type Summary struct {
	TotalChecks int
	TotalErrors int
}

func Cases() []TemplateCase {
	return []TemplateCase{
		{Key: "demo_category", Name: "demo_category (树形+排序+Switch枚举+Tooltip)", Meta: buildCategoryMeta()},
		{Key: "demo_article", Name: "demo_article (外键+所有组件+金额+密码+字典+搜索+验证)", Meta: buildArticleMeta()},
		{Key: "demo_tag", Name: "demo_tag (最简表+Import+BatchEdit)", Meta: buildTagMeta()},
		{Key: "demo_user_review", Name: "demo_user_review (多段模块名 user_review+跨应用外键)", Meta: buildUserReviewMeta()},
	}
}

func backendTemplateFiles() []string {
	return []string{"backend/api.tpl", "backend/model.tpl", "backend/controller.tpl", "backend/logic.tpl", "backend/service.tpl", "backend/consts.tpl"}
}

func frontendTemplateFiles() []string {
	return []string{"frontend/types.tpl", "frontend/api.tpl", "frontend/list.tpl", "frontend/form.tpl", "frontend/detail-drawer.tpl"}
}

func AllTemplateFiles() []string {
	files := append([]string{}, backendTemplateFiles()...)
	return append(files, frontendTemplateFiles()...)
}

func RunCLI(rootDir string) (Summary, error) {
	outDir := filepath.Join(rootDir, "verify_output")
	if err := os.RemoveAll(outDir); err != nil {
		return Summary{}, fmt.Errorf("清理输出目录失败: %w", err)
	}

	summary := Summary{}
	for _, tc := range Cases() {
		fmt.Printf("\n========== %s ==========\n", tc.Name)
		for _, tplFile := range AllTemplateFiles() {
			summary.TotalChecks++
			rendered, err := RenderTemplateCase(filepath.Join(rootDir, "templates"), tplFile, tc)
			if err != nil {
				summary.TotalErrors++
				continue
			}
			if err := writeRenderedTemplate(outDir, rendered); err != nil {
				fmt.Printf("  FAIL [write] %s: %v\n", tplFile, err)
				summary.TotalErrors++
				continue
			}
			fmt.Printf("  OK   %s (%d bytes)\n", tplFile, len(rendered.Output))
		}
	}
	return summary, nil
}

func RenderTemplateCase(tplDir, tplFile string, tc TemplateCase) (*RenderedTemplate, error) {
	output, err := renderTemplateOutput(tplDir, tplFile, tc.Meta)
	if err != nil {
		return nil, err
	}
	return &RenderedTemplate{
		CaseKey:      tc.Key,
		TemplateFile: tplFile,
		Output:       output,
	}, nil
}

func renderTemplateOutput(tplDir, tplFile string, meta *parser.TableMeta) (string, error) {
	tplPath := filepath.Join(tplDir, tplFile)
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(util.SharedFuncMap).ParseFiles(tplPath)
	if err != nil {
		return "", fmt.Errorf("parse %s failed: %w", tplFile, err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, meta); err != nil {
		return "", fmt.Errorf("render %s failed: %w", tplFile, err)
	}
	output := buf.String()
	errs := checkOutput(tplFile, output, meta)
	if len(errs) > 0 {
		return "", fmt.Errorf("%s checks failed: %s", tplFile, strings.Join(errs, "; "))
	}
	if strings.HasPrefix(tplFile, "backend/") {
		if err := checkGoSyntax(tplFile, output); err != nil {
			return "", err
		}
	}
	return output, nil
}

func writeRenderedTemplate(outDir string, rendered *RenderedTemplate) error {
	if rendered == nil {
		return fmt.Errorf("rendered template is nil")
	}
	outPath := filepath.Join(outDir, rendered.CaseKey, rendered.TemplateFile+".out")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(rendered.Output), 0o644); err != nil {
		return fmt.Errorf("写入输出失败: %w", err)
	}
	return nil
}

func checkGoSyntax(tplFile, output string) error {
	fset := token.NewFileSet()
	if _, err := goparser.ParseFile(fset, tplFile, output, goparser.AllErrors); err != nil {
		return fmt.Errorf("生成代码存在语法错误: %w", err)
	}
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

	for _, f := range meta.Fields {
		if f.IsHidden || f.IsID {
			continue
		}
		if !parser.IsSupportedComponent(f.Component) {
			errs = append(errs, fmt.Sprintf("字段 %s 使用了未登记的组件 %s；必须先在 vue-vben-admin/apps/web-antd/src/adapter/component/index.ts 适配后再生成", f.Name, f.Component))
		}
	}

	if isBackend {
		// Go 代码检查
		chk(!strings.Contains(output, "{{"), "包含未渲染的模板标记 {{")
		chk(strings.Contains(output, "package "), "缺少 package 声明")

		if strings.Contains(tplFile, "logic") {
			chk(strings.Contains(output, "func init()"), "缺少 init() 函数")
			chk(strings.Contains(output, "applyListFilter"), "缺少 applyListFilter")
			chk(strings.Contains(output, "isAllowedOrderField"), "缺少 isAllowedOrderField")
			chk(strings.Contains(output, "applyListOrder"), "缺少 applyListOrder")
			chk(strings.Contains(output, "in = &model."+meta.ModelName+"ListInput{}"), "List/Export 缺少 nil 入参保护")
			if meta.HasParentID {
				chk(strings.Contains(output, "collectChildIDs"), "树形表缺少 collectChildIDs")
				chk(strings.Contains(output, "doCollectChildIDs"), "树形表缺少 doCollectChildIDs")
				chk(strings.Contains(output, "collectDeleteIDs"), "树形表缺少 collectDeleteIDs")
				chk(strings.Contains(output, "Tree("), "树形表缺少 Tree 方法")
				chk(strings.Contains(output, "in = &model."+meta.ModelName+"TreeInput{}"), "树形表 Tree 缺少 nil 入参保护")
			}
			chk(strings.Contains(output, "BatchDelete("), "缺少 BatchDelete 方法")
			if meta.HasMoney {
				chk(strings.Contains(output, "LockUpdate"), "金额表缺少行锁")
				chk(strings.Contains(output, "Transaction"), "金额表缺少事务")
			}
			if meta.HasPassword {
				chk(strings.Contains(output, "bcrypt"), "密码表缺少 bcrypt")
			}
			if meta.HasImport {
				chk(strings.Contains(output, "Import("), "缺少 Import 方法")
				chk(strings.Contains(output, "io.EOF"), "导入逻辑缺少 io.EOF 结束处理")
				if meta.HasCreatedBy {
					chk(strings.Contains(output, "Columns().CreatedBy"), "导入逻辑缺少 created_by 注入")
				}
				if meta.HasDeptID {
					chk(strings.Contains(output, "Columns().DeptId"), "导入逻辑缺少 dept_id 注入")
				}
			} else {
				chk(!strings.Contains(output, "Import("), "未启用导入时不应生成 Import 方法")
			}
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "BatchUpdate("), "缺少 BatchUpdate 方法")
			} else {
				chk(!strings.Contains(output, "BatchUpdate("), "未启用批量编辑时不应生成 BatchUpdate 方法")
			}
			if meta.HasCreatedBy || meta.HasDeptID {
				chk(strings.Contains(output, "middleware.ApplyDataScope"), "缺少数据权限过滤")
				chk(strings.Contains(output, `"gbaseadmin/app/`+meta.AppName+`/internal/middleware"`), "缺少 middleware 导入")
			}
			if meta.HasTenantScope {
				chk(strings.Contains(output, "ApplyTenantScopeToWrite"), "缺少租户写入守卫 ApplyTenantScopeToWrite")
				chk(strings.Contains(output, "EnsureTenantMerchantAccessible"), "缺少租户/商户可访问性校验")
				chk(strings.Contains(output, "ApplyTenantScopeToModel"), "列表查询缺少租户范围过滤 ApplyTenantScopeToModel")
				chk(strings.Contains(output, "EnsureTenantScopedRow"), "行级操作缺少租户行守卫")
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
			chk(strings.Contains(output, "BatchDeleteReq"), "缺少 BatchDeleteReq")
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "BatchUpdateReq"), "缺少 BatchUpdateReq")
			} else {
				chk(!strings.Contains(output, "BatchUpdateReq"), "未启用批量编辑时不应生成 BatchUpdateReq")
			}
			if meta.HasParentID {
				chk(strings.Contains(output, "TreeReq"), "树形表缺少 TreeReq")
			} else {
				chk(!strings.Contains(output, "TreeReq"), "非树形表不应生成 TreeReq")
			}
			if meta.HasImport {
				chk(strings.Contains(output, "ImportReq"), "缺少 ImportReq")
				chk(strings.Contains(output, "ImportTemplateReq"), "缺少 ImportTemplateReq")
			} else {
				chk(!strings.Contains(output, "ImportReq"), "未启用导入时不应生成 ImportReq")
				chk(!strings.Contains(output, "ImportTemplateReq"), "未启用导入时不应生成 ImportTemplateReq")
			}
			if meta.HasKeywordSearch {
				chk(strings.Contains(output, `json:"keyword"`), "关键词搜索缺少 keyword 请求字段")
			} else {
				chk(!strings.Contains(output, `json:"keyword"`), "未启用关键词搜索时不应生成 keyword 请求字段")
			}
			if strings.Contains(tplFile, "backend/api") {
				chk(strings.Contains(output, `type `+meta.ModelName+`ExportReq struct`), "缺少 ExportReq")
				chk(strings.Contains(output, `json:"orderBy"`), "ExportReq 缺少 orderBy")
				chk(strings.Contains(output, `json:"orderDir"`), "ExportReq 缺少 orderDir")
			}
			for _, f := range meta.Fields {
				if f.IsHidden || f.IsID || len(f.ValidationRules) == 0 {
					continue
				}
				for _, rule := range f.ValidationRules {
					chk(strings.Contains(output, rule), f.Name+" 缺少后端校验规则 "+rule)
				}
			}
		}

		if strings.Contains(tplFile, "model") {
			chk(strings.Contains(output, "CreateInput"), "缺少 CreateInput")
			chk(strings.Contains(output, "UpdateInput"), "缺少 UpdateInput")
			chk(strings.Contains(output, "DetailOutput"), "缺少 DetailOutput")
			chk(strings.Contains(output, "ListOutput"), "缺少 ListOutput")
			chk(strings.Contains(output, "ListInput"), "缺少 ListInput")
			if meta.HasKeywordSearch {
				chk(strings.Contains(output, `json:"keyword"`), "关键词搜索缺少 keyword 输入字段")
			} else {
				chk(!strings.Contains(output, `json:"keyword"`), "未启用关键词搜索时不应生成 keyword 输入字段")
			}
			if meta.HasParentID {
				chk(strings.Contains(output, "TreeOutput"), "树形表缺少 TreeOutput")
				chk(strings.Contains(output, "TreeInput"), "树形表缺少 TreeInput")
			} else {
				chk(!strings.Contains(output, "TreeOutput"), "非树形表不应生成 TreeOutput")
				chk(!strings.Contains(output, "TreeInput"), "非树形表不应生成 TreeInput")
			}
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "BatchUpdateInput"), "缺少 BatchUpdateInput")
			} else {
				chk(!strings.Contains(output, "BatchUpdateInput"), "未启用批量编辑时不应生成 BatchUpdateInput")
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
			chk(strings.Contains(output, "BatchDelete("), "controller 缺少 BatchDelete")
			if meta.HasImport {
				chk(strings.Contains(output, "ImportTemplate"), "缺少 ImportTemplate")
			} else {
				chk(!strings.Contains(output, "ImportTemplate"), "未启用导入时不应生成 ImportTemplate")
			}
		}

		if strings.Contains(tplFile, "service") {
			chk(strings.Contains(output, "BatchDelete("), "service 缺少 BatchDelete")
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "BatchUpdate("), "缺少 BatchUpdate service 方法")
			} else {
				chk(!strings.Contains(output, "BatchUpdate("), "未启用批量编辑时不应生成 BatchUpdate service 方法")
			}
			if meta.HasImport {
				chk(strings.Contains(output, "Import("), "缺少 Import service 方法")
			} else {
				chk(!strings.Contains(output, "Import("), "未启用导入时不应生成 Import service 方法")
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
		chk(!strings.Contains(output, "rules: ["), "前端模板不应生成数组 rules，必须服从当前 vben 表单规则")
		chk(!strings.Contains(output, "document.createElement"), "前端模板不应生成裸 DOM createElement 交互")
		if strings.Contains(tplFile, "list") {
			chk(strings.Contains(output, "downloadFileFromBlob"), "list 应使用 vben 下载工具")
			chk(strings.Contains(output, "const sortableFieldMap"), "list 缺少排序字段映射")
			chk(strings.Contains(output, "resolveSortField"), "list 缺少排序字段转换 helper")
			chk(strings.Contains(output, "getSortColumns"), "list 导出应读取当前表格排序状态")
			chk(strings.Contains(output, "useAccess"), "list 缺少 useAccess 权限能力")
			chk(strings.Contains(output, "const canBatchDelete = hasAccessByCodes(['"+meta.AppName+":"+meta.ModuleName+":batch-delete'])"), "list 缺少 batch-delete 权限判断")
			chk(strings.Contains(output, "checkboxConfig: canBatchDelete ? { highlight: true } : undefined"), "list 复选框高亮必须跟 batch-delete 权限绑定")
			chk(strings.Contains(output, "...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : [])"), "list 复选框列必须跟 batch-delete 权限绑定")
			if meta.HasKeywordSearch {
				chk(strings.Contains(output, "fieldName: 'keyword'"), "关键词搜索缺少 keyword 控件")
			}
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
			}
			if meta.HasTenantScope {
				chk(strings.Contains(output, "usePlatformSuperAdmin"), "租户表列表缺少平台超级管理员判断")
				chk(strings.Contains(output, "isPlatformSuperAdmin.value ? ["), "tenant_id/merchant_id 列表与筛选应按平台超级管理员条件渲染")
			}
			chk(strings.Contains(output, "checkbox"), "列表应有 checkbox")
			chk(strings.Contains(output, "handleBatchDelete"), "列表缺少 handleBatchDelete")
			chk(strings.Contains(output, ":batch-delete"), "列表缺少 batch-delete 权限按钮")
			if meta.HasImport {
				chk(strings.Contains(output, "handleImportTrigger"), "缺少 handleImportTrigger")
				chk(strings.Contains(output, "handleImportChange"), "缺少 handleImportChange")
				chk(strings.Contains(output, "handleDownloadTemplate"), "缺少 handleDownloadTemplate")
				chk(strings.Contains(output, "ref<HTMLInputElement | null>(null)"), "导入应使用 Vue ref 挂载 input")
				chk(strings.Contains(output, "accept=\".csv\""), "导入应限制为 CSV 格式")
			} else {
				chk(!strings.Contains(output, "handleImportTrigger"), "未启用导入时不应生成 handleImportTrigger")
				chk(!strings.Contains(output, "handleImportChange"), "未启用导入时不应生成 handleImportChange")
				chk(!strings.Contains(output, "handleDownloadTemplate"), "未启用导入时不应生成模板下载")
			}
			if !meta.HasParentID {
				if meta.HasSort {
					chk(strings.Contains(output, "defaultSort: { field: 'sort', order: 'asc' }"), "含 sort 字段的列表默认排序应为 sort asc")
				} else {
					chk(strings.Contains(output, "defaultSort: { field: 'createdAt', order: 'desc' }"), "无 sort 字段的列表默认排序应为 createdAt desc")
				}
			} else {
				chk(!strings.Contains(output, "sortConfig"), "树形表不应生成远程排序配置 sortConfig")
			}
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "handleBatchUpdateStatus"), "缺少 handleBatchUpdateStatus")
				chk(strings.Contains(output, ":batch-update"), "缺少 batch-update 权限按钮")
			} else {
				chk(!strings.Contains(output, "handleBatchUpdateStatus"), "未启用批量编辑时不应生成 handleBatchUpdateStatus")
				chk(!strings.Contains(output, ":batch-update"), "未启用批量编辑时不应出现 batch-update 权限按钮")
			}
			for _, f := range meta.SearchFields {
				chk(strings.Contains(output, "fieldName: '"+f.SearchFormField+"'"), "缺少搜索字段 "+f.SearchFormField)
				if f.SearchRange {
					chk(strings.Contains(output, f.NameLower+"Start"), f.Name+" 区间搜索缺少 start 参数映射")
					chk(strings.Contains(output, f.NameLower+"End"), f.Name+" 区间搜索缺少 end 参数映射")
				}
				if f.DictType != "" {
					chk(strings.Contains(output, "getDictByType"), f.Name+" 字典搜索缺少 getDictByType")
				}
				if f.IsForeignKey && f.RefTableCamel != "" {
					chk(strings.Contains(output, "get"+f.RefTableCamel), f.Name+" 外键搜索缺少选项加载")
				}
			}
		}

		if strings.Contains(tplFile, "form") {
			chk(strings.Contains(output, "const openToken = ref(0);"), "form 缺少 openToken 防串写保护")
			chk(strings.Contains(output, "if (!isOpen)"), "form 应先处理关闭分支")
			chk(strings.Contains(output, "formApi.resetForm();"), "form 打开时应先 resetForm")
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
			chk(strings.Contains(output, "const openToken = ref(0);"), "detail-drawer 缺少 openToken 防串写保护")
			chk(strings.Contains(output, "if (!isOpen)"), "detail-drawer 应先处理关闭分支")
			if meta.HasEnum {
				chk(strings.Contains(output, "Tag"), "有枚举但 detail-drawer 缺少 Tag")
			}
			// 检查图片字段
			if meta.HasImage {
				chk(strings.Contains(output, "<img"), "有图片字段但 detail-drawer 缺少 img 标签")
			}
			chk(strings.Contains(output, "displayValue("), "detail-drawer 缺少 displayValue 兜底")
			for _, f := range meta.Fields {
				if f.Component == "InputUrl" && !f.IsHidden && !f.IsID && !f.IsPassword {
					chk(strings.Contains(output, ":href=\"detail."+f.NameLower+"\""), f.Name+" 详情应渲染为可点击链接")
					chk(strings.Contains(output, `rel="noreferrer noopener"`), f.Name+" 外链应补 rel 安全属性")
					break
				}
			}
		}

		if strings.Contains(tplFile, "types") {
			chk(strings.Contains(output, "interface "+meta.ModelName+"Item"), "缺少 Item 接口")
			chk(strings.Contains(output, "interface "+meta.ModelName+"ListParams"), "缺少 ListParams 接口")
			chk(strings.Contains(output, "interface "+meta.ModelName+"CreateParams"), "缺少 CreateParams 接口")
			chk(strings.Contains(output, "interface "+meta.ModelName+"UpdateParams"), "缺少 UpdateParams 接口")
			if meta.HasKeywordSearch {
				chk(strings.Contains(output, "keyword?: string;"), "关键词搜索缺少 keyword 参数")
			} else {
				chk(!strings.Contains(output, "keyword?: string;"), "未启用关键词搜索时不应生成 keyword 参数")
			}
			for _, f := range meta.SearchFields {
				if f.SearchRange {
					chk(strings.Contains(output, f.NameLower+"Start?: string;"), f.Name+" 缺少区间 start 参数定义")
					chk(strings.Contains(output, f.NameLower+"End?: string;"), f.Name+" 缺少区间 end 参数定义")
					continue
				}
				chk(strings.Contains(output, f.NameLower+"?: "), f.Name+" 缺少列表查询参数定义")
			}
			if meta.HasParentID {
				chk(strings.Contains(output, "children?"), "树形表 Item 缺少 children")
				chk(strings.Contains(output, "TreeParams"), "树形表缺少 TreeParams")
			} else {
				chk(!strings.Contains(output, "TreeParams"), "非树形表不应生成 TreeParams")
			}
		}

		if strings.HasSuffix(tplFile, "frontend/api.tpl") {
			chk(strings.Contains(output, "get"+meta.ModelName+"Detail"), "缺少 getDetail")
			chk(strings.Contains(output, "create"+meta.ModelName), "缺少 create")
			chk(strings.Contains(output, "update"+meta.ModelName), "缺少 update")
			chk(strings.Contains(output, "delete"+meta.ModelName), "缺少 delete")
			chk(strings.Contains(output, "export"+meta.ModelName), "缺少 export")
			chk(strings.Contains(output, "Partial<"+meta.ModelName+"ListParams>"), "export API 应复用 ListParams 类型")
			if meta.HasParentID {
				chk(strings.Contains(output, "get"+meta.ModelName+"Tree"), "树形表缺少 getTree")
			} else {
				chk(strings.Contains(output, "get"+meta.ModelName+"List"), "非树形表缺少 getList")
			}
			chk(strings.Contains(output, "batchDelete"+meta.ModelName), "缺少 batchDelete API")
			if meta.HasImport {
				chk(strings.Contains(output, "import"+meta.ModelName), "缺少 import API")
				chk(strings.Contains(output, "downloadImportTemplate"+meta.ModelName), "缺少模板下载 API")
			} else {
				chk(!strings.Contains(output, "import"+meta.ModelName), "未启用导入时不应生成 import API")
				chk(!strings.Contains(output, "downloadImportTemplate"+meta.ModelName), "未启用导入时不应生成模板下载 API")
			}
			if meta.HasBatchEdit {
				chk(strings.Contains(output, "batchUpdate"+meta.ModelName), "缺少 batchUpdate API")
			} else {
				chk(!strings.Contains(output, "batchUpdate"+meta.ModelName), "未启用批量编辑时不应生成 batchUpdate API")
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
		Comment: "分类",
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		f("parent_id", "ParentID", "ParentId", "parentId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsParentID = true
			f.Component = "TreeSelectSingle"
			f.RefTable = "category"
			f.RefTableDB = "demo_category"
			f.RefTableApp = "demo"
			f.RefTableCamel = "Category"
			f.RefTableLower = "category"
			f.RefDisplayField = "name"
			f.RefDisplayCamel = "Name"
			f.RefDisplayLower = "name"
			f.RefFieldName = "CategoryName"
			f.RefFieldJSON = "categoryName"
			f.RefIsTree = true
			f.RefHasDeletedAt = true
		}),
		f("name", "Name", "Name", "name", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true
			f.IsSearchable = true
			f.MaxLength = 50
			f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:50"}
			f.UpdateValidationRules = []string{"max-length:50"}
		}),
		f("icon", "Icon", "Icon", "icon", "string", "string", func(f *parser.FieldMeta) { f.Component = "IconPicker" }),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.TooltipText = "升序"; f.Component = "InputNumber" }),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Switch"
			f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "禁用", NameIdent: "Disabled"}, {Value: "1", Label: "启用", NameIdent: "Enabled"}}
		}),
		scopeTenantField(),
		scopeMerchantField(),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return finalizeVerifyMeta(m)
}

// --- 表2: demo_article（复杂表：外键+所有组件+金额+密码+字典+搜索+验证规则）---
func buildArticleMeta() *parser.TableMeta {
	m := &parser.TableMeta{
		TableName: "demo_article", AppName: "demo", AppNameCamel: "Demo",
		ModelName: "Article", DaoName: "DemoArticle", ModuleName: "article", PackageName: "article",
		Comment: "文章",
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		// 同应用树形外键
		f("category_id", "CategoryID", "CategoryId", "categoryId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsForeignKey = true
			f.Component = "Select"
			f.IsRequired = true
			f.ValidationRules = []string{"required"}
			f.RefTable = "category"
			f.RefTableDB = "demo_category"
			f.RefTableApp = "demo"
			f.RefTableCamel = "Category"
			f.RefTableLower = "category"
			f.RefDisplayField = "name"
			f.RefDisplayCamel = "Name"
			f.RefDisplayLower = "name"
			f.RefFieldName = "CategoryName"
			f.RefFieldJSON = "categoryName"
			f.RefIsTree = true
			f.RefHasDeletedAt = true
		}),
		// 跨应用普通外键
		f("user_id", "UserID", "UserId", "userId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsForeignKey = true
			f.Component = "Select"
			f.IsRequired = true
			f.ValidationRules = []string{"required"}
			f.RefTable = "users"
			f.RefTableDB = "system_users"
			f.RefTableApp = "system"
			f.RefTableCamel = "Users"
			f.RefTableLower = "users"
			f.RefDisplayField = "username"
			f.RefDisplayCamel = "Username"
			f.RefDisplayLower = "username"
			f.RefFieldName = "UsersUsername"
			f.RefFieldJSON = "usersUsername"
			f.RefIsTree = false
			f.RefHasDeletedAt = true
		}),
		f("title", "Title", "Title", "title", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true
			f.IsSearchable = true
			f.MaxLength = 200
			f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:200"}
			f.UpdateValidationRules = []string{"max-length:200"}
		}),
		f("order_no", "OrderNo", "OrderNo", "orderNo", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true
			f.IsSearchable = true
			f.IsExactSearch = true
			f.MaxLength = 50
			f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:50"}
			f.UpdateValidationRules = []string{"max-length:50"}
		}),
		f("cover", "Cover", "Cover", "cover", "string", "string", func(f *parser.FieldMeta) { f.Component = "ImageUpload" }),
		f("attachment_file", "AttachmentFile", "AttachmentFile", "attachmentFile", "string", "string", func(f *parser.FieldMeta) { f.Component = "FileUpload" }),
		f("body_content", "BodyContent", "BodyContent", "bodyContent", "string", "string", func(f *parser.FieldMeta) { f.DBType = "text"; f.Component = "RichText" }),
		f("extra_json", "ExtraJSON", "ExtraJson", "extraJSON", "string", "string", func(f *parser.FieldMeta) { f.DBType = "text"; f.Component = "JsonEditor" }),
		f("link_url", "LinkURL", "LinkUrl", "linkURL", "string", "string", func(f *parser.FieldMeta) {
			f.Component = "InputUrl"
			f.FrontendRules = "url"
			f.MaxLength = 500
			f.ValidationRules = []string{"url", "max-length:500"}
			f.UpdateValidationRules = []string{"url", "max-length:500"}
		}),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Radio"
			f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "草稿", NameIdent: "Draft"}, {Value: "1", Label: "已发布", NameIdent: "Published"}, {Value: "2", Label: "已下架", NameIdent: "Offline"}}
		}),
		f("type", "Type", "Type", "type", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Select"
			f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "1", Label: "普通", NameIdent: "Regular"}, {Value: "2", Label: "置顶", NameIdent: "Pinned"}, {Value: "3", Label: "推荐", NameIdent: "Recommended"}, {Value: "4", Label: "热门", NameIdent: "Hot"}}
		}),
		f("is_top", "IsTop", "IsTop", "isTop", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Switch"
			f.DefaultValue = "0"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "否", NameIdent: "No"}, {Value: "1", Label: "是", NameIdent: "Yes"}}
		}),
		f("price", "Price", "Price", "price", "int", "number", func(f *parser.FieldMeta) { f.IsMoney = true; f.Component = "InputNumber" }),
		f("pay_password", "PayPassword", "PayPassword", "payPassword", "string", "string", func(f *parser.FieldMeta) {
			f.IsPassword = true
			f.Component = "Password"
			f.ValidationRules = []string{"length:6,32"}
			f.UpdateValidationRules = []string{"length:6,32"}
		}),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.TooltipText = "升序"; f.Component = "InputNumber" }),
		f("icon", "Icon", "Icon", "icon", "string", "string", func(f *parser.FieldMeta) { f.Component = "IconPicker" }),
		f("email", "Email", "Email", "email", "string", "string", func(f *parser.FieldMeta) {
			f.Component = "Input"
			f.FrontendRules = "email"
			f.MaxLength = 100
			f.ValidationRules = []string{"email", "max-length:100"}
			f.UpdateValidationRules = []string{"email", "max-length:100"}
		}),
		f("phone", "Phone", "Phone", "phone", "string", "string", func(f *parser.FieldMeta) {
			f.IsSearchable = true
			f.Component = "Input"
			f.FrontendRules = "phone"
			f.MaxLength = 20
			f.ValidationRules = []string{"phone-loose", "max-length:20"}
			f.UpdateValidationRules = []string{"phone-loose", "max-length:20"}
		}),
		f("remark", "Remark", "Remark", "remark", "string", "string", func(f *parser.FieldMeta) { f.DBType = "text"; f.IsSearchable = true; f.Component = "Textarea" }),
		f("level", "Level", "Level", "level", "string", "string", func(f *parser.FieldMeta) { f.DictType = "article_level"; f.Component = "Select" }),
		f("extra_field", "ExtraField", "ExtraField", "extraField", "string", "string", func(f *parser.FieldMeta) { f.Component = "Input" }),
		f("publish_at", "PublishAt", "PublishAt", "publishAt", "*gtime.Time", "string", func(f *parser.FieldMeta) { f.IsTimeField = true; f.Component = "DateTimePicker" }),
		f("expire_at", "ExpireAt", "ExpireAt", "expireAt", "*gtime.Time", "string", func(f *parser.FieldMeta) { f.IsTimeField = true; f.Component = "DateTimePicker" }),
		scopeTenantField(),
		scopeMerchantField(),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return finalizeVerifyMeta(m)
}

// --- 表3: demo_tag（最简表）---
func buildTagMeta() *parser.TableMeta {
	m := &parser.TableMeta{
		TableName: "demo_tag", AppName: "demo", AppNameCamel: "Demo",
		ModelName: "Tag", DaoName: "DemoTag", ModuleName: "tag", PackageName: "tag",
		Comment: "标签",
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		f("name", "Name", "Name", "name", "string", "string", func(f *parser.FieldMeta) {
			f.IsRequired = true
			f.IsSearchable = true
			f.MaxLength = 50
			f.Component = "Input"
			f.ValidationRules = []string{"required", "max-length:50"}
			f.UpdateValidationRules = []string{"max-length:50"}
		}),
		f("color", "Color", "Color", "color", "string", "string", func(f *parser.FieldMeta) { f.Component = "Input"; f.MaxLength = 20 }),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.Component = "InputNumber" }),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Switch"
			f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "禁用", NameIdent: "Disabled"}, {Value: "1", Label: "启用", NameIdent: "Enabled"}}
		}),
		scopeTenantField(),
		scopeMerchantField(),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return finalizeVerifyMeta(m)
}

// --- 表4: demo_user_review（多段模块名 user_review + 跨应用外键）---
// 唯一核心验证点：表名第一个下划线之后保留下划线，moduleName=user_review，
// 包路径 app/demo/internal/logic/user_review/，前端 views/demo/user_review/。
func buildUserReviewMeta() *parser.TableMeta {
	m := &parser.TableMeta{
		TableName: "demo_user_review", AppName: "demo", AppNameCamel: "Demo",
		ModelName: "UserReview", DaoName: "DemoUserReview", ModuleName: "user_review", PackageName: "user_review",
		Comment: "用户审核",
	}
	m.Fields = []parser.FieldMeta{
		f("id", "ID", "Id", "id", "JsonInt64", "string", func(f *parser.FieldMeta) { f.IsID = true; f.IsHidden = true }),
		// 跨应用外键：system_users.username
		f("user_id", "UserID", "UserId", "userId", "JsonInt64", "string", func(f *parser.FieldMeta) {
			f.IsForeignKey = true
			f.Component = "Select"
			f.IsRequired = true
			f.ValidationRules = []string{"required"}
			f.RefTable = "users"
			f.RefTableDB = "system_users"
			f.RefTableApp = "system"
			f.RefTableCamel = "Users"
			f.RefTableLower = "users"
			f.RefDisplayField = "username"
			f.RefDisplayCamel = "Username"
			f.RefDisplayLower = "username"
			f.RefFieldName = "UsersUsername"
			f.RefFieldJSON = "usersUsername"
			f.RefIsTree = false
			f.RefHasDeletedAt = true
		}),
		f("review_type", "ReviewType", "ReviewType", "reviewType", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Radio"
			f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "1", Label: "内容", NameIdent: "Content"}, {Value: "2", Label: "行为", NameIdent: "Behavior"}, {Value: "3", Label: "申诉", NameIdent: "Appeal"}}
		}),
		f("content", "Content", "Content", "content", "string", "string", func(f *parser.FieldMeta) {
			f.IsSearchable = true
			f.MaxLength = 500
			f.Component = "Input"
			f.ValidationRules = []string{"max-length:500"}
			f.UpdateValidationRules = []string{"max-length:500"}
		}),
		f("score", "Score", "Score", "score", "int", "number", func(f *parser.FieldMeta) { f.Component = "InputNumber" }),
		f("is_passed", "IsPassed", "IsPassed", "isPassed", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Switch"
			f.DefaultValue = "0"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "否", NameIdent: "No"}, {Value: "1", Label: "是", NameIdent: "Yes"}}
		}),
		f("sort", "Sort", "Sort", "sort", "int", "number", func(f *parser.FieldMeta) { f.TooltipText = "升序"; f.Component = "InputNumber" }),
		f("status", "Status", "Status", "status", "int", "number", func(f *parser.FieldMeta) {
			f.IsEnum = true
			f.Component = "Switch"
			f.DefaultValue = "1"
			f.EnumValues = []parser.EnumValue{{Value: "0", Label: "禁用", NameIdent: "Disabled"}, {Value: "1", Label: "启用", NameIdent: "Enabled"}}
		}),
		scopeTenantField(),
		scopeMerchantField(),
		hiddenField("created_by", "CreatedBy", "CreatedBy", "createdBy", "JsonInt64", "string"),
		hiddenField("dept_id", "DeptID", "DeptId", "deptID", "JsonInt64", "string"),
		hiddenTimeField("created_at", "CreatedAt"), hiddenTimeField("updated_at", "UpdatedAt"), hiddenTimeField("deleted_at", "DeletedAt"),
	}
	return finalizeVerifyMeta(m)
}

func finalizeVerifyMeta(m *parser.TableMeta) *parser.TableMeta {
	for i := range m.Fields {
		parser.ApplySearchMeta(&m.Fields[i])
	}
	parser.FinalizeTemplateMeta(m)
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

func scopeTenantField() parser.FieldMeta {
	return f("tenant_id", "TenantID", "TenantId", "tenantID", "JsonInt64", "string", func(f *parser.FieldMeta) {
		f.Label = "租户"
		f.ShortLabel = "租户"
		f.Component = parser.ComponentSelect
		f.IsForeignKey = true
		f.RefTable = "tenant"
		f.RefTableDB = "system_tenant"
		f.RefTableApp = "system"
		f.RefTableCamel = "Tenant"
		f.RefTableLower = "tenant"
		f.RefDisplayField = "name"
		f.RefDisplayCamel = "Name"
		f.RefDisplayLower = "name"
		f.RefFieldName = "TenantName"
		f.RefFieldJSON = "tenantName"
	})
}

func scopeMerchantField() parser.FieldMeta {
	return f("merchant_id", "MerchantID", "MerchantId", "merchantID", "JsonInt64", "string", func(f *parser.FieldMeta) {
		f.Label = "商户"
		f.ShortLabel = "商户"
		f.Component = parser.ComponentSelect
		f.IsForeignKey = true
		f.RefTable = "merchant"
		f.RefTableDB = "system_merchant"
		f.RefTableApp = "system"
		f.RefTableCamel = "Merchant"
		f.RefTableLower = "merchant"
		f.RefDisplayField = "name"
		f.RefDisplayCamel = "Name"
		f.RefDisplayLower = "name"
		f.RefFieldName = "MerchantName"
		f.RefFieldJSON = "merchantName"
	})
}

func hiddenTimeField(name, camel string) parser.FieldMeta {
	dao := strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), " ", "")
	lower := strings.ToLower(camel[:1]) + camel[1:]
	return f(name, camel, dao, lower, "*gtime.Time", "string", func(f *parser.FieldMeta) { f.IsHidden = true; f.IsTimeField = true })
}
