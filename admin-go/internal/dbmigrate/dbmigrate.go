package dbmigrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlDriver "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gopkg.in/yaml.v3"
)

const (
	defaultMigrationsDir = "database/migrations"
	defaultDatabaseNode  = "default"
	legacyBaselineVer    = 1
)

var migrationNameSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

type Options struct {
	Action       string
	Name         string
	Dir          string
	DSN          string
	Link         string
	ConfigPath   string
	DatabaseNode string
	Steps        int
	Version      int
	Now          time.Time
}

type gfConfig struct {
	Database map[string]struct {
		Link string `yaml:"link"`
	} `yaml:"database"`
}

func Run(ctx context.Context, opts Options) error {
	opts = normalizeOptions(opts)
	switch opts.Action {
	case "create":
		return createMigrationFiles(opts.Dir, opts.Name, opts.Now)
	case "up", "down", "version", "force":
	default:
		return fmt.Errorf("unsupported migrate action: %s", opts.Action)
	}

	dir, err := resolveMigrationsDir(opts.Dir)
	if err != nil {
		return err
	}
	if !hasMigrationFiles(dir) {
		switch opts.Action {
		case "version":
			fmt.Println("migration version: none")
			return nil
		case "up":
			fmt.Printf("no migration files found under %s\n", dir)
			return nil
		default:
			return fmt.Errorf("no migration files found under %s", dir)
		}
	}

	dsn, err := resolveDSN(opts)
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		return err
	}

	driver, err := mysqlDriver.WithInstance(db, &mysqlDriver.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fileSourceURL(dir), "mysql", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	switch opts.Action {
	case "up":
		if err := baselineExistingSchema(ctx, db, m); err != nil {
			return err
		}
		err = m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("migration up: no change")
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println("migration up: success")
		return nil
	case "down":
		steps := opts.Steps
		if steps <= 0 {
			steps = 1
		}
		err = m.Steps(-steps)
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("migration down: no change")
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("migration down: rolled back %d step(s)\n", steps)
		return nil
	case "force":
		if opts.Version < 0 {
			return errors.New("force action requires a non-negative --version")
		}
		if err := m.Force(opts.Version); err != nil {
			return err
		}
		fmt.Printf("migration force: set version to %d\n", opts.Version)
		return nil
	case "version":
		version, dirty, err := m.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("migration version: none")
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("migration version: %d dirty=%t\n", version, dirty)
		return nil
	default:
		return fmt.Errorf("unsupported migrate action: %s", opts.Action)
	}
}

func normalizeOptions(opts Options) Options {
	opts.Action = strings.ToLower(strings.TrimSpace(opts.Action))
	opts.Name = strings.TrimSpace(opts.Name)
	opts.Dir = strings.TrimSpace(opts.Dir)
	opts.DSN = strings.TrimSpace(opts.DSN)
	opts.Link = strings.TrimSpace(opts.Link)
	opts.ConfigPath = strings.TrimSpace(opts.ConfigPath)
	opts.DatabaseNode = strings.TrimSpace(opts.DatabaseNode)
	if opts.Dir == "" {
		opts.Dir = defaultMigrationsDir
	}
	if opts.DatabaseNode == "" {
		opts.DatabaseNode = defaultDatabaseNode
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	return opts
}

func resolveDSN(opts Options) (string, error) {
	if opts.DSN != "" {
		return ensureMySQLDSNParams(opts.DSN), nil
	}
	if value := strings.TrimSpace(os.Getenv("GBASEADMIN_DB_DSN")); value != "" {
		return ensureMySQLDSNParams(value), nil
	}
	if opts.Link != "" {
		return goFrameLinkToMySQLDSN(opts.Link)
	}
	if value := strings.TrimSpace(os.Getenv("GBASEADMIN_DB_LINK")); value != "" {
		return goFrameLinkToMySQLDSN(value)
	}
	if dsn := envMySQLDSN(); dsn != "" {
		return dsn, nil
	}

	configPath, err := resolveConfigPath(opts.ConfigPath)
	if err != nil {
		return "", err
	}
	link, err := loadDatabaseLink(configPath, opts.DatabaseNode)
	if err != nil {
		return "", err
	}
	return goFrameLinkToMySQLDSN(link)
}

func envMySQLDSN() string {
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASS"))
	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	if host == "" || user == "" || name == "" {
		return ""
	}
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	if port == "" {
		port = "3306"
	}
	return ensureMySQLDSNParams(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name))
}

func resolveConfigPath(provided string) (string, error) {
	candidates := make([]string, 0, 6)
	if provided != "" {
		candidates = append(candidates, provided)
	}
	if value := strings.TrimSpace(os.Getenv("GBASEADMIN_MIGRATE_CONFIG")); value != "" {
		candidates = append(candidates, value)
	}
	if cfgPath := strings.TrimSpace(os.Getenv("GF_GCFG_PATH")); cfgPath != "" {
		cfgFile := strings.TrimSpace(os.Getenv("GF_GCFG_FILE"))
		if cfgFile == "" {
			cfgFile = "config.yaml"
		}
		if filepath.IsAbs(cfgFile) {
			candidates = append(candidates, cfgFile)
		} else {
			candidates = append(candidates, filepath.Join(cfgPath, cfgFile))
		}
	}
	candidates = append(candidates,
		filepath.Join("manifest", "config", "config.yaml"),
		filepath.Join("app", "system", "manifest", "config", "config.yaml"),
		filepath.Join("app", "upload", "manifest", "config", "config.yaml"),
		filepath.Join("app", "system", "hack", "config.yaml"),
		filepath.Join("app", "upload", "hack", "config.yaml"),
	)
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	return "", errors.New("unable to locate database config file for migrations")
}

func loadDatabaseLink(path, node string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var cfg gfConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	if cfg.Database == nil {
		return "", fmt.Errorf("database config missing in %s", path)
	}
	item, ok := cfg.Database[node]
	if !ok || strings.TrimSpace(item.Link) == "" {
		return "", fmt.Errorf("database.%s.link missing in %s", node, path)
	}
	return item.Link, nil
}

func goFrameLinkToMySQLDSN(link string) (string, error) {
	link = strings.TrimSpace(link)
	if link == "" {
		return "", errors.New("empty database link")
	}
	if strings.HasPrefix(link, "mysql:") {
		return ensureMySQLDSNParams(strings.TrimPrefix(link, "mysql:")), nil
	}
	if strings.Contains(link, "://") {
		return "", fmt.Errorf("unsupported database link: %s", link)
	}
	return ensureMySQLDSNParams(link), nil
}

func ensureMySQLDSNParams(dsn string) string {
	dsn = strings.TrimSpace(dsn)
	base := dsn
	rawQuery := ""
	if idx := strings.Index(dsn, "?"); idx >= 0 {
		base = dsn[:idx]
		rawQuery = dsn[idx+1:]
	}
	values, _ := url.ParseQuery(rawQuery)
	if values.Get("charset") == "" {
		values.Set("charset", "utf8mb4")
	}
	if values.Get("parseTime") == "" {
		values.Set("parseTime", "true")
	}
	if values.Get("multiStatements") == "" {
		values.Set("multiStatements", "true")
	}
	encoded := values.Encode()
	if encoded == "" {
		return base
	}
	return base + "?" + encoded
}

func createMigrationFiles(dir, name string, now time.Time) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("create action requires a migration name")
	}
	absDir, err := resolveMigrationsDir(dir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return err
	}
	version := now.UTC().Format("20060102150405")
	safeName := normalizeMigrationName(name)
	upPath := filepath.Join(absDir, fmt.Sprintf("%s_%s.up.sql", version, safeName))
	downPath := filepath.Join(absDir, fmt.Sprintf("%s_%s.down.sql", version, safeName))
	if fileExists(upPath) || fileExists(downPath) {
		return fmt.Errorf("migration files already exist for version %s and name %s", version, safeName)
	}
	if err := os.WriteFile(upPath, []byte("-- Write your migration here.\n"), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(downPath, []byte("-- Write your rollback here.\n"), 0o644); err != nil {
		return err
	}
	fmt.Printf("created migration files:\n%s\n%s\n", upPath, downPath)
	return nil
}

func normalizeMigrationName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = migrationNameSanitizer.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		return "migration"
	}
	return name
}

func resolveMigrationsDir(dir string) (string, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = defaultMigrationsDir
	}
	return filepath.Abs(dir)
}

func hasMigrationFiles(dir string) bool {
	matches, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	return err == nil && len(matches) > 0
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func fileSourceURL(dir string) string {
	return "file://" + filepath.ToSlash(dir)
}

func baselineExistingSchema(ctx context.Context, db *sql.DB, m *migrate.Migrate) error {
	version, _, err := m.Version()
	if err == nil && version > 0 {
		return nil
	}
	if !errors.Is(err, migrate.ErrNilVersion) {
		return err
	}
	exists, err := legacySchemaExists(ctx, db)
	if err != nil || !exists {
		return err
	}
	if err := m.Force(legacyBaselineVer); err != nil {
		return err
	}
	fmt.Printf("existing schema detected without migration metadata, forced baseline version %d\n", legacyBaselineVer)
	return nil
}

func legacySchemaExists(ctx context.Context, db *sql.DB) (bool, error) {
	const query = `
SELECT COUNT(*)
FROM information_schema.tables
WHERE table_schema = DATABASE()
  AND table_name IN ('system_users', 'system_menu', 'upload_config')`
	var count int
	if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
