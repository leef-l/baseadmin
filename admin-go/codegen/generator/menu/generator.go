package menu

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"

	_ "github.com/go-sql-driver/mysql"

	"gbaseadmin/codegen/parser"
)

// MenuAppConfig 应用目录的标题、图标和排序配置
type MenuAppConfig struct {
	Title string
	Icon  string
	Sort  int
}

// MenuModuleConfig 模块级菜单配置
type MenuModuleConfig struct {
	Sort   int
	IsShow *int
	Icon   string
}

// Config 菜单生成器配置
type Config struct {
	DSN         string
	Force       bool
	DryRun      bool
	MenuApps    map[string]MenuAppConfig    // 从 codegen.yaml 加载的应用目录配置
	MenuModules map[string]MenuModuleConfig // 从 codegen.yaml 加载的模块级配置，key: "appName/moduleName"
}

// Generator 菜单生成器
type Generator struct {
	config     Config
	operations []PreviewOperation
}

type PreviewOperation struct {
	Kind       string `json:"kind"`
	Action     string `json:"action"`
	Title      string `json:"title"`
	Path       string `json:"path,omitempty"`
	Permission string `json:"permission,omitempty"`
	Component  string `json:"component,omitempty"`
	Sort       int    `json:"sort,omitempty"`
	IsShow     int    `json:"isShow,omitempty"`
}

type menuStore interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

const (
	menuTypeDirectory = 1
	menuTypePage      = 2
	menuTypeButton    = 3
)

type menuIconRule struct {
	icon     string
	keywords []string
}

var menuIconByModule = map[string]string{
	"api":      "ApiOutlined",
	"config":   "ControlOutlined",
	"dept":     "ApartmentOutlined",
	"dict":     "TagsOutlined",
	"dir":      "FolderOpenOutlined",
	"dir_rule": "PartitionOutlined",
	"file":     "FileTextOutlined",
	"log":      "ProfileOutlined",
	"menu":     "MenuOutlined",
	"notice":   "NotificationOutlined",
	"role":     "TeamOutlined",
	"table":    "DatabaseOutlined",
	"users":    "UserOutlined",
	"user":     "UserOutlined",
}

var menuIconRules = []menuIconRule{
	{icon: "MenuOutlined", keywords: []string{"菜单", "导航", "menu", "menus"}},
	{icon: "UserOutlined", keywords: []string{"用户", "账号", "成员", "user", "users", "account", "member"}},
	{icon: "TeamOutlined", keywords: []string{"角色", "岗位", "团队", "权限组", "role", "roles", "team"}},
	{icon: "ApartmentOutlined", keywords: []string{"部门", "组织", "机构", "dept", "department", "organization", "org"}},
	{icon: "ControlOutlined", keywords: []string{"配置", "设置", "参数", "config", "setting", "settings", "option", "options"}},
	{icon: "PartitionOutlined", keywords: []string{"规则", "策略", "rule", "rules", "policy"}},
	{icon: "FolderOpenOutlined", keywords: []string{"目录", "文件夹", "folder", "dir", "directory"}},
	{icon: "FileTextOutlined", keywords: []string{"文件", "文档", "记录", "报表", "file", "files", "document", "record", "report"}},
	{icon: "CloudUploadOutlined", keywords: []string{"上传", "存储", "upload", "storage", "bucket", "oss", "cos", "s3"}},
	{icon: "ApiOutlined", keywords: []string{"接口", "api", "openapi", "webhook"}},
	{icon: "NotificationOutlined", keywords: []string{"通知", "消息", "公告", "notice", "notification", "message", "mail"}},
	{icon: "TagsOutlined", keywords: []string{"字典", "标签", "分类", "dict", "dictionary", "tag", "tags", "category"}},
	{icon: "ProfileOutlined", keywords: []string{"日志", "审计", "log", "logs", "audit"}},
	{icon: "DatabaseOutlined", keywords: []string{"数据", "数据库", "schema", "database", "data"}},
}

// New 创建菜单生成器
func New(cfg Config) *Generator {
	return &Generator{config: cfg}
}

func (g *Generator) PreviewOperations() []PreviewOperation {
	return append([]PreviewOperation(nil), g.operations...)
}

// Generate 为一张表生成菜单数据
func (g *Generator) Generate(meta *parser.TableMeta) (int, error) {
	return g.GenerateBatch([]*parser.TableMeta{meta})
}

// GenerateBatch 为多张表批量生成菜单数据。
// 非 dry-run 模式下目录、菜单和按钮会在同一个事务中提交，避免留下半套菜单。
func (g *Generator) GenerateBatch(metas []*parser.TableMeta) (int, error) {
	g.operations = nil
	if len(metas) == 0 {
		return 0, nil
	}
	if err := validateBatchMetas(metas); err != nil {
		return 0, err
	}

	db, err := g.openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	return g.generateBatchWithDB(db, metas)
}

func (g *Generator) openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", g.config.DSN)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	if _, err := db.Exec("SET NAMES utf8mb4"); err != nil {
		db.Close()
		return nil, fmt.Errorf("设置字符集失败: %w", err)
	}
	return db, nil
}

func (g *Generator) generateBatchWithDB(db *sql.DB, metas []*parser.TableMeta) (int, error) {
	if g.config.DryRun {
		total := 0
		for _, meta := range metas {
			if meta == nil {
				continue
			}
			count, err := g.generateWithStore(db, meta)
			if err != nil {
				return total, err
			}
			total += count
		}
		return total, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("开启菜单事务失败: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	total := 0
	for _, meta := range metas {
		if meta == nil {
			continue
		}
		count, err := g.generateWithStore(tx, meta)
		if err != nil {
			return total, err
		}
		total += count
	}

	if err := tx.Commit(); err != nil {
		return total, fmt.Errorf("提交菜单事务失败: %w", err)
	}
	committed = true
	return total, nil
}

func (g *Generator) generateWithStore(store menuStore, meta *parser.TableMeta) (int, error) {
	count := 0

	// 1. 查找或创建应用目录
	dirPath := "/" + meta.AppName
	dirID, err := g.ensureDirectory(store, meta.AppName, dirPath)
	if err != nil {
		return count, err
	}

	// 2. 读取模块级配置
	menuSort := 0
	menuIsShow := 1
	menuIcon := guessMenuIcon(meta)
	moduleKey := meta.AppName + "/" + meta.ModuleName
	if modCfg, ok := g.config.MenuModules[moduleKey]; ok {
		if modCfg.Sort > 0 {
			menuSort = modCfg.Sort
		}
		if modCfg.IsShow != nil {
			menuIsShow = *modCfg.IsShow
		}
		if modCfg.Icon != "" {
			menuIcon = modCfg.Icon
		}
	}

	// 3. 插入菜单页
	menuTitle := buildMenuTitle(meta)
	menuPath := "/" + meta.AppName + "/" + dashCase(meta.ModuleName)
	menuComponent := meta.AppName + "/" + meta.ModuleName + "/index"
	menuPermission := meta.AppName + ":" + meta.ModuleName + ":list"

	menuID, created, err := g.ensureMenu(store, dirID, menuTitle, menuPath, menuComponent, menuPermission, menuIcon, menuSort, menuIsShow)
	if err != nil {
		return count, err
	}
	if created {
		count++
	}

	// 4. 插入按钮权限
	for _, btn := range buildButtonSpecs(meta) {
		btnTitle := menuTitle + btn.suffix
		created, err := g.ensureButton(store, menuID, btnTitle, btn.permission, btn.sort)
		if err != nil {
			return count, err
		}
		if created {
			count++
		}
	}

	return count, nil
}

func (g *Generator) recordOperation(op PreviewOperation) {
	g.operations = append(g.operations, op)
}

func validateBatchMetas(metas []*parser.TableMeta) error {
	pathOwners := make(map[string]string, len(metas))
	permissionOwners := make(map[string]string, len(metas)*4)

	for idx, meta := range metas {
		if meta == nil {
			return fmt.Errorf("第 %d 个菜单元数据为空", idx+1)
		}
		if strings.TrimSpace(meta.AppName) == "" {
			return fmt.Errorf("模块 %q 缺少应用名，无法生成菜单", meta.TableName)
		}
		if strings.TrimSpace(meta.ModuleName) == "" {
			return fmt.Errorf("表 %q 缺少模块名，无法生成菜单", meta.TableName)
		}

		title := buildMenuTitle(meta)
		if strings.TrimSpace(title) == "" {
			return fmt.Errorf("表 %q 生成的菜单标题为空", meta.TableName)
		}

		moduleKey := meta.AppName + "/" + meta.ModuleName
		menuPath := "/" + meta.AppName + "/" + dashCase(meta.ModuleName)
		if owner, exists := pathOwners[menuPath]; exists {
			return fmt.Errorf("菜单 path 冲突: %s 同时来自 %s 和 %s", menuPath, owner, moduleKey)
		}
		pathOwners[menuPath] = moduleKey

		menuPermission := meta.AppName + ":" + meta.ModuleName + ":list"
		if owner, exists := permissionOwners[menuPermission]; exists {
			return fmt.Errorf("菜单 permission 冲突: %s 同时来自 %s 和 %s", menuPermission, owner, moduleKey)
		}
		permissionOwners[menuPermission] = moduleKey

		for _, btn := range buildButtonSpecs(meta) {
			if strings.TrimSpace(btn.permission) == "" {
				return fmt.Errorf("模块 %s 生成了空按钮权限", moduleKey)
			}
			if owner, exists := permissionOwners[btn.permission]; exists {
				return fmt.Errorf("按钮 permission 冲突: %s 同时来自 %s 和 %s", btn.permission, owner, moduleKey)
			}
			permissionOwners[btn.permission] = moduleKey
		}
	}

	return nil
}

// ensureDirectory 查找或创建应用目录（type=1）
func (g *Generator) ensureDirectory(store menuStore, appName, path string) (int64, error) {
	// 查找已存在的目录
	id, err := findMenuID(store, "path", path, menuTypeDirectory)
	if err == nil {
		g.recordOperation(PreviewOperation{Kind: "directory", Action: "existing", Title: path, Path: path, Sort: 0, IsShow: 1})
		fmt.Printf("  [菜单] 目录已存在: %s (ID: %d)\n", path, id)
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("查询目录失败: %w", err)
	}

	// 创建目录
	id = generateID()
	title := appName + "管理"
	icon := "AppstoreOutlined"
	sortVal := 50
	if cfg, ok := g.config.MenuApps[appName]; ok {
		if cfg.Title != "" {
			title = cfg.Title
		}
		if cfg.Icon != "" {
			icon = cfg.Icon
		}
		if cfg.Sort > 0 {
			sortVal = cfg.Sort
		}
	}

	if g.config.DryRun {
		g.recordOperation(PreviewOperation{Kind: "directory", Action: "create", Title: title, Path: path, Sort: sortVal, IsShow: 1})
		fmt.Printf("  [dry-run] INSERT 目录: title=%s, path=%s, icon=%s, sort=%d\n", title, path, icon, sortVal)
		return id, nil
	}

	_, err = store.Exec(
		`INSERT INTO system_menu (id, parent_id, title, type, path, component, permission, icon, sort, is_show, is_cache, status, created_by, dept_id, tenant_id, merchant_id, created_at, updated_at)
		 VALUES (?, 0, ?, ?, ?, NULL, '', ?, ?, 1, 0, 1, 0, 0, 0, 0, NOW(), NOW())`,
		id, title, menuTypeDirectory, path, icon, sortVal,
	)
	if err != nil {
		return 0, fmt.Errorf("创建目录失败: %w", err)
	}
	g.recordOperation(PreviewOperation{Kind: "directory", Action: "create", Title: title, Path: path, Sort: sortVal, IsShow: 1})
	fmt.Printf("  [菜单] 创建目录: %s (ID: %d, sort: %d)\n", title, id, sortVal)
	return id, nil
}

// ensureMenu 查找或创建菜单页（type=2）
func (g *Generator) ensureMenu(store menuStore, parentID int64, title, path, component, permission, icon string, sort int, isShow int) (int64, bool, error) {
	id, err := findMenuID(store, "path", path, menuTypePage)
	if err == nil {
		if g.config.Force {
			_, err = store.Exec(
				`UPDATE system_menu SET title=?, component=?, permission=?, icon=?, sort=?, is_show=?, updated_at=NOW() WHERE id=?`,
				title, component, permission, icon, sort, isShow, id,
			)
			if err != nil {
				return 0, false, fmt.Errorf("更新菜单失败: %w", err)
			}
			g.recordOperation(PreviewOperation{Kind: "menu", Action: "update", Title: title, Path: path, Permission: permission, Component: component, Sort: sort, IsShow: isShow})
			fmt.Printf("  [菜单] 更新菜单: %s (%s)\n", title, path)
			return id, false, nil
		}
		g.recordOperation(PreviewOperation{Kind: "menu", Action: "skip", Title: title, Path: path, Permission: permission, Component: component, Sort: sort, IsShow: isShow})
		fmt.Printf("  [菜单] 跳过（已存在）: %s (%s)\n", title, path)
		return id, false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, false, fmt.Errorf("查询菜单失败: %w", err)
	}

	id = generateID()

	if g.config.DryRun {
		g.recordOperation(PreviewOperation{Kind: "menu", Action: "create", Title: title, Path: path, Permission: permission, Component: component, Sort: sort, IsShow: isShow})
		fmt.Printf("  [dry-run] INSERT 菜单: title=%s, path=%s, permission=%s, icon=%s, sort=%d, is_show=%d\n", title, path, permission, icon, sort, isShow)
		return id, true, nil
	}

	_, err = store.Exec(
		`INSERT INTO system_menu (id, parent_id, title, type, path, component, permission, icon, sort, is_show, is_cache, status, created_by, dept_id, tenant_id, merchant_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 1, 0, 0, 0, 0, NOW(), NOW())`,
		id, parentID, title, menuTypePage, path, component, permission, icon, sort, isShow,
	)
	if err != nil {
		return 0, false, fmt.Errorf("创建菜单失败: %w", err)
	}
	g.recordOperation(PreviewOperation{Kind: "menu", Action: "create", Title: title, Path: path, Permission: permission, Component: component, Sort: sort, IsShow: isShow})
	fmt.Printf("  [菜单] 创建菜单: %s (%s)\n", title, path)
	return id, true, nil
}

// ensureButton 查找或创建按钮权限（type=3）
func (g *Generator) ensureButton(store menuStore, parentID int64, title, permission string, sort int) (bool, error) {
	id, err := findMenuID(store, "permission", permission, menuTypeButton)
	if err == nil {
		if g.config.Force {
			_, err = store.Exec(
				`UPDATE system_menu SET title=?, sort=?, updated_at=NOW() WHERE id=?`,
				title, sort, id,
			)
			if err != nil {
				return false, fmt.Errorf("更新按钮失败: %w", err)
			}
			g.recordOperation(PreviewOperation{Kind: "button", Action: "update", Title: title, Permission: permission, Sort: sort})
			fmt.Printf("  [菜单] 更新按钮: %s (%s)\n", title, permission)
			return false, nil
		}
		g.recordOperation(PreviewOperation{Kind: "button", Action: "skip", Title: title, Permission: permission, Sort: sort})
		fmt.Printf("  [菜单] 跳过（已存在）: %s (%s)\n", title, permission)
		return false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("查询按钮失败: %w", err)
	}

	id = generateID()

	if g.config.DryRun {
		g.recordOperation(PreviewOperation{Kind: "button", Action: "create", Title: title, Permission: permission, Sort: sort})
		fmt.Printf("  [dry-run] INSERT 按钮: title=%s, permission=%s\n", title, permission)
		return true, nil
	}

	_, err = store.Exec(
		`INSERT INTO system_menu (id, parent_id, title, type, path, component, permission, icon, sort, is_show, is_cache, status, created_by, dept_id, tenant_id, merchant_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, NULL, NULL, ?, '', ?, 0, 0, 1, 0, 0, 0, 0, NOW(), NOW())`,
		id, parentID, title, menuTypeButton, permission, sort,
	)
	if err != nil {
		return false, fmt.Errorf("创建按钮失败: %w", err)
	}
	g.recordOperation(PreviewOperation{Kind: "button", Action: "create", Title: title, Permission: permission, Sort: sort})
	fmt.Printf("  [菜单] 创建按钮: %s (%s)\n", title, permission)
	return true, nil
}

// cleanTitle 从表注释中提取简短标题
func cleanTitle(comment string) string {
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return ""
	}
	// 去掉常见后缀
	for _, suffix := range []string{"表", "管理"} {
		for strings.HasSuffix(comment, suffix) {
			comment = strings.TrimSpace(strings.TrimSuffix(comment, suffix))
		}
	}
	return comment
}

func buildMenuTitle(meta *parser.TableMeta) string {
	title := cleanTitle(meta.Comment)
	if title != "" {
		return title
	}
	if meta != nil && meta.ModelName != "" {
		return meta.ModelName
	}
	if meta != nil {
		return meta.ModuleName
	}
	return ""
}

func guessMenuIcon(meta *parser.TableMeta) string {
	if meta == nil {
		return "AppstoreOutlined"
	}

	moduleName := strings.ToLower(strings.TrimSpace(meta.ModuleName))
	if icon, ok := menuIconByModule[moduleName]; ok {
		return icon
	}

	searchText := normalizeMenuIconText(
		meta.ModuleName,
		meta.ModelName,
		cleanTitle(meta.Comment),
		meta.Comment,
	)
	for _, rule := range menuIconRules {
		for _, keyword := range rule.keywords {
			if matchMenuKeyword(searchText, keyword) {
				return rule.icon
			}
		}
	}

	return "AppstoreOutlined"
}

func normalizeMenuIconText(parts ...string) string {
	var builder strings.Builder
	for _, part := range parts {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "" {
			continue
		}
		part = strings.NewReplacer("_", " ", "-", " ", "/", " ").Replace(part)
		if builder.Len() > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(part)
	}
	return builder.String()
}

func matchMenuKeyword(searchText, keyword string) bool {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return false
	}

	if isASCIIKeyword(keyword) {
		for _, token := range strings.Fields(searchText) {
			if token == keyword {
				return true
			}
		}
		return false
	}

	return strings.Contains(searchText, keyword)
}

func isASCIIKeyword(keyword string) bool {
	for _, r := range keyword {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

type buttonSpec struct {
	suffix     string
	permission string
	sort       int
}

func buildButtonSpecs(meta *parser.TableMeta) []buttonSpec {
	buttons := []buttonSpec{
		{suffix: "新增", permission: meta.AppName + ":" + meta.ModuleName + ":create", sort: 1},
		{suffix: "修改", permission: meta.AppName + ":" + meta.ModuleName + ":update", sort: 2},
		{suffix: "删除", permission: meta.AppName + ":" + meta.ModuleName + ":delete", sort: 3},
		{suffix: "批量删除", permission: meta.AppName + ":" + meta.ModuleName + ":batch-delete", sort: 4},
	}
	buttons = append(buttons,
		buttonSpec{suffix: "查看", permission: meta.AppName + ":" + meta.ModuleName + ":detail", sort: 5},
		buttonSpec{suffix: "导出", permission: meta.AppName + ":" + meta.ModuleName + ":export", sort: 6},
	)
	if meta.HasImport {
		buttons = append(buttons, buttonSpec{suffix: "导入", permission: meta.AppName + ":" + meta.ModuleName + ":import", sort: 7})
	}
	if meta.HasBatchEdit {
		buttons = append(buttons, buttonSpec{suffix: "批量编辑", permission: meta.AppName + ":" + meta.ModuleName + ":batch-update", sort: 8})
	}
	return buttons
}

// dashCase 将 snake_case 的模块名转为 dash-case（用于 URL path）
func dashCase(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

var allowedMenuFields = map[string]bool{
	"path":       true,
	"permission": true,
}

func findMenuID(store menuStore, field string, value any, menuType int) (int64, error) {
	if store == nil {
		return 0, sql.ErrConnDone
	}
	if !allowedMenuFields[field] {
		return 0, fmt.Errorf("非法查询字段: %s", field)
	}
	var id int64
	query := fmt.Sprintf("SELECT id FROM system_menu WHERE %s = ? AND type = ? AND deleted_at IS NULL", field)
	err := store.QueryRow(query, value, menuType).Scan(&id)
	return id, err
}

// --- 内联 Snowflake ID 生成（与项目 utility/snowflake 算法一致）---

const (
	sfEpoch          = int64(1700000000000)
	sfWorkerBits     = uint(10)
	sfSequenceBits   = uint(12)
	sfSequenceMax    = int64(-1) ^ (int64(-1) << sfSequenceBits)
	sfWorkerShift    = sfSequenceBits
	sfTimestampShift = sfSequenceBits + sfWorkerBits
)

var sfGen = newSFGenerator(1)

type sfGenerator struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
	nowMillis func() int64
}

func newSFGenerator(workerID int64) *sfGenerator {
	return &sfGenerator{
		workerID: workerID,
		nowMillis: func() int64 {
			return time.Now().UnixMilli()
		},
	}
}

func generateID() int64 {
	return sfGen.nextID()
}

func (g *sfGenerator) nextID() int64 {
	if g == nil {
		return 0
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.currentTimestamp()
	if now < g.timestamp {
		now = g.timestamp
	}
	if now == g.timestamp {
		g.sequence = (g.sequence + 1) & sfSequenceMax
		if g.sequence == 0 {
			for now <= g.timestamp {
				now = g.currentTimestamp()
			}
		}
	} else {
		g.sequence = 0
	}
	g.timestamp = now

	return (now << sfTimestampShift) | (g.workerID << sfWorkerShift) | g.sequence
}

func (g *sfGenerator) currentTimestamp() int64 {
	now := time.Now().UnixMilli()
	if g != nil && g.nowMillis != nil {
		now = g.nowMillis()
	}
	return now - sfEpoch
}
