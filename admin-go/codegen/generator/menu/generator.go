package menu

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

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
	config Config
}

const (
	menuTypeDirectory = 1
	menuTypePage      = 2
	menuTypeButton    = 3
)

// New 创建菜单生成器
func New(cfg Config) *Generator {
	return &Generator{config: cfg}
}

// Generate 为一张表生成菜单数据
func (g *Generator) Generate(meta *parser.TableMeta) (int, error) {
	db, err := sql.Open("mysql", g.config.DSN)
	if err != nil {
		return 0, fmt.Errorf("连接数据库失败: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec("SET NAMES utf8mb4"); err != nil {
		return 0, fmt.Errorf("设置字符集失败: %w", err)
	}

	count := 0

	// 1. 查找或创建应用目录
	dirPath := "/" + meta.AppName
	dirID, err := g.ensureDirectory(db, meta.AppName, dirPath)
	if err != nil {
		return count, err
	}

	// 2. 读取模块级配置
	menuSort := 0
	menuIsShow := 1
	moduleKey := meta.AppName + "/" + meta.ModuleName
	if modCfg, ok := g.config.MenuModules[moduleKey]; ok {
		if modCfg.Sort > 0 {
			menuSort = modCfg.Sort
		}
		if modCfg.IsShow != nil {
			menuIsShow = *modCfg.IsShow
		}
	}

	// 3. 插入菜单页
	menuTitle := buildMenuTitle(meta)
	menuPath := "/" + meta.AppName + "/" + dashCase(meta.ModuleName)
	menuComponent := meta.AppName + "/" + meta.ModuleName + "/index"
	menuPermission := meta.AppName + ":" + meta.ModuleName + ":list"

	menuID, created, err := g.ensureMenu(db, dirID, menuTitle, menuPath, menuComponent, menuPermission, menuSort, menuIsShow)
	if err != nil {
		return count, err
	}
	if created {
		count++
	}

	// 4. 插入按钮权限
	for _, btn := range buildButtonSpecs(meta) {
		btnTitle := menuTitle + btn.suffix
		created, err := g.ensureButton(db, menuID, btnTitle, btn.permission, btn.sort)
		if err != nil {
			return count, err
		}
		if created {
			count++
		}
	}

	return count, nil
}

// ensureDirectory 查找或创建应用目录（type=1）
func (g *Generator) ensureDirectory(db *sql.DB, appName, path string) (int64, error) {
	// 查找已存在的目录
	id, err := findMenuID(db, "path", path, menuTypeDirectory)
	if err == nil {
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
		fmt.Printf("  [dry-run] INSERT 目录: title=%s, path=%s, icon=%s, sort=%d\n", title, path, icon, sortVal)
		return id, nil
	}

	_, err = db.Exec(
		`INSERT INTO system_menu (id, parent_id, title, type, path, component, permission, icon, sort, is_show, is_cache, status, created_by, dept_id, created_at, updated_at)
		 VALUES (?, 0, ?, ?, ?, NULL, '', ?, ?, 1, 0, 1, 0, 0, NOW(), NOW())`,
		id, title, menuTypeDirectory, path, icon, sortVal,
	)
	if err != nil {
		return 0, fmt.Errorf("创建目录失败: %w", err)
	}
	fmt.Printf("  [菜单] 创建目录: %s (ID: %d, sort: %d)\n", title, id, sortVal)
	return id, nil
}

// ensureMenu 查找或创建菜单页（type=2）
func (g *Generator) ensureMenu(db *sql.DB, parentID int64, title, path, component, permission string, sort int, isShow int) (int64, bool, error) {
	id, err := findMenuID(db, "path", path, menuTypePage)
	if err == nil {
		if g.config.Force {
			_, err = db.Exec(
				`UPDATE system_menu SET title=?, component=?, permission=?, sort=?, is_show=?, updated_at=NOW() WHERE id=?`,
				title, component, permission, sort, isShow, id,
			)
			if err != nil {
				return 0, false, fmt.Errorf("更新菜单失败: %w", err)
			}
			fmt.Printf("  [菜单] 更新菜单: %s (%s)\n", title, path)
			return id, false, nil
		}
		fmt.Printf("  [菜单] 跳过（已存在）: %s (%s)\n", title, path)
		return id, false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, false, fmt.Errorf("查询菜单失败: %w", err)
	}

	id = generateID()

	if g.config.DryRun {
		fmt.Printf("  [dry-run] INSERT 菜单: title=%s, path=%s, permission=%s, sort=%d, is_show=%d\n", title, path, permission, sort, isShow)
		return id, true, nil
	}

	_, err = db.Exec(
		`INSERT INTO system_menu (id, parent_id, title, type, path, component, permission, icon, sort, is_show, is_cache, status, created_by, dept_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, ?, 0, 1, 0, 0, NOW(), NOW())`,
		id, parentID, title, menuTypePage, path, component, permission, sort, isShow,
	)
	if err != nil {
		return 0, false, fmt.Errorf("创建菜单失败: %w", err)
	}
	fmt.Printf("  [菜单] 创建菜单: %s (%s)\n", title, path)
	return id, true, nil
}

// ensureButton 查找或创建按钮权限（type=3）
func (g *Generator) ensureButton(db *sql.DB, parentID int64, title, permission string, sort int) (bool, error) {
	id, err := findMenuID(db, "permission", permission, menuTypeButton)
	if err == nil {
		if g.config.Force {
			_, err = db.Exec(
				`UPDATE system_menu SET title=?, sort=?, updated_at=NOW() WHERE id=?`,
				title, sort, id,
			)
			if err != nil {
				return false, fmt.Errorf("更新按钮失败: %w", err)
			}
			fmt.Printf("  [菜单] 更新按钮: %s (%s)\n", title, permission)
			return false, nil
		}
		fmt.Printf("  [菜单] 跳过（已存在）: %s (%s)\n", title, permission)
		return false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("查询按钮失败: %w", err)
	}

	id = generateID()

	if g.config.DryRun {
		fmt.Printf("  [dry-run] INSERT 按钮: title=%s, permission=%s\n", title, permission)
		return true, nil
	}

	_, err = db.Exec(
		`INSERT INTO system_menu (id, parent_id, title, type, path, component, permission, icon, sort, is_show, is_cache, status, created_by, dept_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, NULL, NULL, ?, '', ?, 0, 0, 1, 0, 0, NOW(), NOW())`,
		id, parentID, title, menuTypeButton, permission, sort,
	)
	if err != nil {
		return false, fmt.Errorf("创建按钮失败: %w", err)
	}
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
	}
	if !meta.HasParentID {
		buttons = append(buttons, buttonSpec{suffix: "批量删除", permission: meta.AppName + ":" + meta.ModuleName + ":batch-delete", sort: 4})
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

func findMenuID(db *sql.DB, field string, value any, menuType int) (int64, error) {
	if db == nil {
		return 0, sql.ErrConnDone
	}
	var id int64
	query := fmt.Sprintf("SELECT id FROM system_menu WHERE %s = ? AND type = ? AND deleted_at IS NULL", field)
	err := db.QueryRow(query, value, menuType).Scan(&id)
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
