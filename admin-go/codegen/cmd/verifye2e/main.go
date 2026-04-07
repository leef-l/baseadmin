package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"

	"gbaseadmin/codegen/generator/backend"
	"gbaseadmin/codegen/generator/frontend"
	"gbaseadmin/codegen/generator/util"
	"gbaseadmin/codegen/internal/runtimeutil"
	"gbaseadmin/codegen/parser"
)

type verifyConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func main() {
	codegenRoot := mustCodegenRoot()
	if err := runtimeutil.LoadEnvFileIfExists(filepath.Join(codegenRoot, "..", ".env")); err != nil {
		fatal(fmt.Errorf("加载环境变量失败: %w", err))
	}

	cfg, err := loadVerifyConfig()
	if err != nil {
		fatal(err)
	}

	db, err := sql.Open("mysql", cfg.dsn())
	if err != nil {
		fatal(fmt.Errorf("连接数据库失败: %w", err))
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fatal(fmt.Errorf("数据库 ping 失败: %w", err))
	}

	cleanupDB, err := prepareVerifySchema(db, cfg, filepath.Join(codegenRoot, "sql", "e2e_verify.sql"))
	if err != nil {
		fatal(err)
	}
	defer cleanupDB()

	tempRoot, err := os.MkdirTemp("", "baseadmin-codegen-e2e-*")
	if err != nil {
		fatal(fmt.Errorf("创建临时目录失败: %w", err))
	}
	keepTemp := false
	defer func() {
		if keepTemp {
			fmt.Printf("[verifye2e] 保留临时目录: %s\n", tempRoot)
			return
		}
		_ = os.RemoveAll(tempRoot)
	}()

	if err := prepareWorkspace(codegenRoot, tempRoot); err != nil {
		keepTemp = true
		fatal(err)
	}

	if err := runE2E(codegenRoot, tempRoot, cfg); err != nil {
		keepTemp = true
		fatal(err)
	}

	fmt.Println("[verifye2e] 验证通过")
}

func runE2E(codegenRoot, tempRoot string, cfg verifyConfig) error {
	tableNames := []string{"verifydemo_category", "verifydemo_article", "verifydemo_tag"}
	p, err := parser.New(cfg.dsn(), []string{"created_at", "updated_at", "deleted_at", "created_by", "dept_id"})
	if err != nil {
		return fmt.Errorf("初始化 parser 失败: %w", err)
	}
	defer p.Close()

	metas, err := p.ParseTables(tableNames)
	if err != nil {
		return fmt.Errorf("解析验证表失败: %w", err)
	}

	templateDir := filepath.Join(codegenRoot, "templates")
	backendOutputRoot := filepath.Join(tempRoot, "app")
	frontendOutputRoot := filepath.Join(tempRoot, "frontend-src")
	cache := util.NewTemplateCache()
	appModules := make(map[string][]string)
	appTables := make(map[string][]string)

	for _, meta := range metas {
		meta.EnableOpLog = false
		backendGen := backend.New(backend.Config{
			TemplateDir: filepath.Join(templateDir, "backend"),
			OutputDir:   filepath.Join(backendOutputRoot, meta.AppName),
			Force:       true,
			Cache:       cache,
		})
		if _, err := backendGen.Generate(meta); err != nil {
			return fmt.Errorf("生成后端失败(%s): %w", meta.TableName, err)
		}

		frontendGen := frontend.New(frontend.Config{
			TemplateDir: filepath.Join(templateDir, "frontend"),
			OutputDir:   frontendOutputRoot,
			Force:       true,
			Cache:       cache,
		})
		if _, err := frontendGen.Generate(meta); err != nil {
			return fmt.Errorf("生成前端失败(%s): %w", meta.TableName, err)
		}

		appModules[meta.AppName] = append(appModules[meta.AppName], meta.ModuleName)
		appTables[meta.AppName] = append(appTables[meta.AppName], meta.TableName)
	}

	for appName, modules := range appModules {
		appDir := filepath.Join(backendOutputRoot, appName)
		if err := os.MkdirAll(appDir, 0o755); err != nil {
			return fmt.Errorf("创建应用目录失败: %w", err)
		}
		sort.Strings(modules)
		sort.Strings(appTables[appName])

		if err := renderTemplate(
			filepath.Join(templateDir, "backend", "hack_config.tpl"),
			filepath.Join(appDir, "hack", "config.yaml"),
			map[string]string{
				"DBLink": cfg.gfDSN(),
				"Tables": strings.Join(appTables[appName], ","),
			},
		); err != nil {
			return fmt.Errorf("生成 hack/config.yaml 失败: %w", err)
		}

		if err := runCommand(appDir, "gf", "gen", "dao"); err != nil {
			return fmt.Errorf("gf gen dao 失败: %w", err)
		}

		if err := renderTemplate(
			filepath.Join(templateDir, "backend", "main.tpl"),
			filepath.Join(appDir, "main.go"),
			map[string]any{
				"AppName": appName,
				"Modules": modules,
			},
		); err != nil {
			return fmt.Errorf("生成 main.go 失败: %w", err)
		}

		if err := renderTemplate(
			filepath.Join(templateDir, "backend", "cmd.tpl"),
			filepath.Join(appDir, "internal", "cmd", "cmd.go"),
			map[string]any{
				"AppName": appName,
				"Modules": modules,
			},
		); err != nil {
			return fmt.Errorf("生成 cmd.go 失败: %w", err)
		}

		if err := copyFile(
			filepath.Join(templateDir, "backend", "middleware_auth.tpl"),
			filepath.Join(appDir, "internal", "middleware", "auth.go"),
		); err != nil {
			return fmt.Errorf("复制 middleware/auth.go 失败: %w", err)
		}
		if err := copyFile(
			filepath.Join(templateDir, "backend", "middleware_context.tpl"),
			filepath.Join(appDir, "internal", "middleware", "context.go"),
		); err != nil {
			return fmt.Errorf("复制 middleware/context.go 失败: %w", err)
		}
		if err := os.MkdirAll(filepath.Join(appDir, "internal", "packed"), 0o755); err != nil {
			return fmt.Errorf("创建 packed 目录失败: %w", err)
		}
		if err := os.WriteFile(filepath.Join(appDir, "internal", "packed", "packed.go"), []byte("package packed\n"), 0o644); err != nil {
			return fmt.Errorf("写 packed.go 失败: %w", err)
		}

		if err := runCommand(tempRoot, "go", "build", "./app/"+appName+"/..."); err != nil {
			return fmt.Errorf("编译生成应用失败(%s): %w", appName, err)
		}
	}

	return nil
}

func prepareWorkspace(codegenRoot, tempRoot string) error {
	adminRoot := filepath.Clean(filepath.Join(codegenRoot, ".."))
	if err := copyFile(filepath.Join(adminRoot, "go.mod"), filepath.Join(tempRoot, "go.mod")); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(adminRoot, "go.sum"), filepath.Join(tempRoot, "go.sum")); err != nil {
		return err
	}
	if err := copyDir(filepath.Join(adminRoot, "utility"), filepath.Join(tempRoot, "utility")); err != nil {
		return err
	}
	return nil
}

func prepareVerifySchema(db *sql.DB, cfg verifyConfig, sqlPath string) (func(), error) {
	systemUsersExists, err := tableExists(db, cfg.DBName, "system_users")
	if err != nil {
		return nil, err
	}

	dropTables := []string{"verifydemo_article", "verifydemo_category", "verifydemo_tag"}
	for _, tableName := range dropTables {
		if _, err := db.Exec("DROP TABLE IF EXISTS `" + tableName + "`"); err != nil {
			return nil, fmt.Errorf("清理旧验证表失败(%s): %w", tableName, err)
		}
	}

	sqlData, err := os.ReadFile(sqlPath)
	if err != nil {
		return nil, fmt.Errorf("读取 e2e SQL 失败: %w", err)
	}
	for _, stmt := range runtimeutil.SplitSQLStatements(string(sqlData)) {
		if _, err := db.Exec(stmt); err != nil {
			return nil, fmt.Errorf("执行 e2e SQL 失败: %w\nSQL: %s", err, stmt)
		}
	}

	return func() {
		for _, tableName := range dropTables {
			_, _ = db.Exec("DROP TABLE IF EXISTS `" + tableName + "`")
		}
		if !systemUsersExists {
			_, _ = db.Exec("DROP TABLE IF EXISTS `system_users`")
		}
	}, nil
}

func tableExists(db *sql.DB, dbName, tableName string) (bool, error) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?",
		dbName, tableName,
	).Scan(&count)
	return count > 0, err
}

func renderTemplate(tplPath, outPath string, data any) error {
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(util.SharedFuncMap).ParseFiles(tplPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return tpl.Execute(file, data)
}

func runCommand(workdir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GOCACHE=/tmp/go-build-cache")
	return cmd.Run()
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func mustCodegenRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fatal(errors.New("无法定位 verifye2e 源码路径"))
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func loadVerifyConfig() (verifyConfig, error) {
	cfg := verifyConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName:   os.Getenv("MYSQL_DATABASE"),
	}
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.DBName == "" {
		return cfg, errors.New("缺少 MySQL 环境变量，请先配置 admin-go/.env")
	}
	return cfg, nil
}

func (c verifyConfig) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (c verifyConfig) gfDSN() string {
	return fmt.Sprintf("mysql:%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[verifye2e] %v\n", err)
	os.Exit(1)
}
