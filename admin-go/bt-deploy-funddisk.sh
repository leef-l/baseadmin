#!/usr/bin/env bash
set -euo pipefail

# ============================================
# FundDisk 宝塔面板一键部署脚本
# 站点: funddisk.easytestdev.online
# 服务: system + upload + member
# 双数据库: system/upload → baseadmin, member → funddisk
# ============================================

# ---------- 配置区 ----------
DOMAIN="funddisk.easytestdev.online"
DEPLOY_DIR="/www/wwwroot/${DOMAIN}"

# system / upload 使用 baseadmin 库
DB_BASEADMIN_NAME="sql_baseadmin_easytestdev_online"
DB_BASEADMIN_USER="sql_baseadmin_easytestdev_online"
DB_BASEADMIN_PASS="4365d2c5953988"

# member 使用 funddisk 库
DB_FUNDDISK_NAME="sql_funddisk_easytestdev_online"
DB_FUNDDISK_USER="sql_funddisk_easytestdev_online"
DB_FUNDDISK_PASS="daa5129899a218"

DB_HOST="127.0.0.1"
DB_PORT="3306"
JWT_SECRET="${JWT_SECRET:-$(openssl rand -hex 32)}"

APPS=("system" "upload" "member")
PORTS=("10032" "10033" "10027")

SERVICE_PREFIX="funddisk"
SSL_MODE="auto"

# 资源限制
GOMEMLIMIT="512MiB"
GOMAXPROCS="2"
GOGC="100"

# 宝塔路径
BT_NGINX_DIR="/www/server/panel/vhost/nginx"
BT_SUPERVISOR_PROFILE_DIR="/www/server/panel/plugin/supervisor/profile"
BT_SUPERVISOR_LOG_DIR="/www/server/panel/plugin/supervisor/log"
BT_SUPERVISOR_CONFIG_JSON="/www/server/panel/plugin/supervisor/config.json"
SUPERVISOR_CONF="/etc/supervisor/supervisord.conf"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
SOURCE_DIR="$(cd "$SCRIPT_DIR/.." && pwd -P)"

# ---------- 跳过选项 ----------
SKIP_BUILD="0"
SKIP_NGINX="0"
SKIP_SUPERVISOR="0"
SKIP_MIGRATE="0"
DRY_RUN="0"

usage() {
  cat <<'EOF'
Usage: bash bt-deploy-funddisk.sh [OPTIONS]

Options:
  --skip-build        跳过编译步骤
  --skip-nginx        跳过 Nginx 配置
  --skip-supervisor   跳过 Supervisor 配置
  --skip-migrate      跳过数据库迁移
  --dry-run           只打印配置，不执行
  -h, --help          显示帮助
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --skip-build)     SKIP_BUILD="1" ;;
    --skip-nginx)     SKIP_NGINX="1" ;;
    --skip-supervisor) SKIP_SUPERVISOR="1" ;;
    --skip-migrate)   SKIP_MIGRATE="1" ;;
    --dry-run)        DRY_RUN="1" ;;
    -h|--help)        usage; exit 0 ;;
    *)                echo "未知选项: $1" >&2; exit 1 ;;
  esac
  shift
done

# ---------- 工具函数 ----------
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
die()   { echo -e "${RED}[ERROR]${NC} $1" >&2; exit 1; }

b64() { printf '%s' "$1" | base64 | tr -d '\n'; }

find_mysql_bin() {
  for candidate in mysql /www/server/mysql/bin/mysql /www/server/mariadb/bin/mysql /usr/bin/mysql /usr/local/mysql/bin/mysql; do
    if command -v "$candidate" >/dev/null 2>&1 || [ -x "$candidate" ]; then
      printf '%s\n' "$candidate"
      return 0
    fi
  done
  return 1
}

find_nginx_bin() {
  if command -v nginx >/dev/null 2>&1; then command -v nginx; return 0; fi
  if [ -x /www/server/nginx/sbin/nginx ]; then printf '%s\n' /www/server/nginx/sbin/nginx; return 0; fi
  return 1
}

find_supervisorctl_bin() {
  if command -v supervisorctl >/dev/null 2>&1; then command -v supervisorctl; return 0; fi
  if [ -x /www/server/panel/pyenv/bin/supervisorctl ]; then printf '%s\n' /www/server/panel/pyenv/bin/supervisorctl; return 0; fi
  return 1
}

find_python_bin() {
  if command -v python3 >/dev/null 2>&1; then command -v python3; return 0; fi
  if command -v python >/dev/null 2>&1; then command -v python; return 0; fi
  if [ -x /www/server/panel/pyenv/bin/python ]; then printf '%s\n' /www/server/panel/pyenv/bin/python; return 0; fi
  return 1
}

# 获取 app 对应的数据库配置
get_db_for_app() {
  local app="$1"
  case "$app" in
    member)
      echo "${DB_FUNDDISK_USER}:${DB_FUNDDISK_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_FUNDDISK_NAME}"
      ;;
    *)
      echo "${DB_BASEADMIN_USER}:${DB_BASEADMIN_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_BASEADMIN_NAME}"
      ;;
  esac
}

# ---------- 加载 Go 环境 ----------
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
if ! command -v go &>/dev/null; then
  die "找不到 go 命令，请确认 Go 已安装"
fi

# ---------- 检查 CPU 负载 ----------
check_cpu() {
  local cpu_idle
  cpu_idle=$(top -bn1 | grep 'Cpu' | awk '{print $8}' | tr -d '%id,')
  local cpu_used=$((100 - ${cpu_idle%.*}))
  if [ "$cpu_used" -gt 80 ]; then
    die "CPU 使用率 ${cpu_used}% 超过 80%，请等待负载降低后再执行"
  fi
  info "当前 CPU 使用率: ${cpu_used}%"
}

# ---------- 打印配置 ----------
info "========================================="
info "FundDisk 部署配置"
info "========================================="
info "域名: ${DOMAIN}"
info "部署目录: ${DEPLOY_DIR}"
info "服务: ${APPS[*]}"
info "端口: ${PORTS[*]}"
info "Supervisor 前缀: ${SERVICE_PREFIX}"
info "system/upload 数据库: ${DB_BASEADMIN_NAME}"
info "member 数据库: ${DB_FUNDDISK_NAME}"
info "资源限制: GOMEMLIMIT=${GOMEMLIMIT} GOMAXPROCS=${GOMAXPROCS}"
info "========================================="

if [ "$DRY_RUN" = "1" ]; then
  info "dry-run 模式，不执行任何操作"
  exit 0
fi

check_cpu

# ---------- 1. 创建目录 ----------
info "创建部署目录..."
mkdir -p "$DEPLOY_DIR"
for app in "${APPS[@]}"; do
  mkdir -p "$DEPLOY_DIR/$app/manifest/config"
  mkdir -p "$DEPLOY_DIR/$app/logs"
done
mkdir -p "$DEPLOY_DIR/upload/resource/upload"
mkdir -p "$DEPLOY_DIR/database/migrations"
mkdir -p "$DEPLOY_DIR/scripts"

# ---------- 2. 编译 ----------
if [ "$SKIP_BUILD" = "1" ]; then
  warn "跳过编译"
else
  info "开始编译..."
  cd "$SCRIPT_DIR"
  for i in "${!APPS[@]}"; do
    app="${APPS[$i]}"
    check_cpu
    info "编译 $app ..."
    cd "$SCRIPT_DIR/app/$app"
    CGO_ENABLED=0 GOOS=linux go build -o "$DEPLOY_DIR/$app/$app" .
    chmod +x "$DEPLOY_DIR/$app/$app"
    cd "$SCRIPT_DIR"
    info "$app 编译完成"
  done
  info "编译 migrate ..."
  cd "$SCRIPT_DIR"
  CGO_ENABLED=0 GOOS=linux go build -o "$DEPLOY_DIR/migrate" ./cmd/migrate
  chmod +x "$DEPLOY_DIR/migrate"
  info "migrate 编译完成"
fi

# 复制迁移文件
cp -R "$SCRIPT_DIR/database/migrations/." "$DEPLOY_DIR/database/migrations/"

# ---------- 3. 生成配置文件 ----------
# member 服务额外需要 jwt.memberSecret / redis / sms / member.* 配置块
# 仅当目标 config.yaml 不存在时才生成；已有配置不覆盖（避免运营手动改的值丢失）
JWT_MEMBER_SECRET="${JWT_MEMBER_SECRET:-$(openssl rand -hex 32)}"
REDIS_ADDR="${REDIS_ADDR:-127.0.0.1:6379}"
REDIS_PASS="${REDIS_PASS:-}"
REDIS_DB="${REDIS_DB:-0}"

info "生成配置文件..."
for i in "${!APPS[@]}"; do
  app="${APPS[$i]}"
  port="${PORTS[$i]}"
  db_link="$(get_db_for_app "$app")"
  conf="$DEPLOY_DIR/$app/manifest/config/config.yaml"

  if [ -f "$conf" ]; then
    info "$app 配置已存在，保留不覆盖: $conf"
    continue
  fi

  cat > "$conf" <<EOF
server:
  address: ":${port}"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

logger:
  level: "warning"
  stdout: true
  path: "${DEPLOY_DIR}/${app}/logs"

database:
  default:
    link: "mysql:${db_link}"
    debug: false
    charset: "utf8mb4"

redis:
  default:
    address: "${REDIS_ADDR}"
    pass: "${REDIS_PASS}"
    db: ${REDIS_DB}

jwt:
  secret: "${JWT_SECRET}"
  memberSecret: "${JWT_MEMBER_SECRET}"
  expire: 24
EOF

  if [ "$app" = "member" ]; then
    cat >> "$conf" <<EOF

# 短信配置（默认 mock；生产改 aliyun 并填 AK/SK/signName/templateCode）
sms:
  provider: "mock"
  codeExpireSeconds: 300
  limitSeconds: 60
  verifyMaxAttempts: 5
  providers:
    mock:
      kind: "mock"
      fixedCode: "123456"
    aliyun:
      kind: "aliyun"
      region: "cn-hangzhou"
      accessKeyId: ""
      accessKeySecret: ""
      signName: ""
      templateCode: ""

# 会员业务配置
member:
  h5BaseURL: "https://${DOMAIN}/h5"
  inviteCodeLength: 8
  registerRequireInviteCode: true
  contractStorageDir: "${DEPLOY_DIR}/member/runtime/contracts"
  contractChromiumPath: ""
  teamExportRoot: "${DEPLOY_DIR}/member/runtime/team-export"
EOF
  fi

  info "$app 配置已生成 (端口: $port)"
done

# ---------- 4. 生成启动脚本 ----------
info "生成启动脚本..."

cat > "$DEPLOY_DIR/scripts/start-service.sh" <<'STARTEOF'
#!/bin/sh
set -eu

APP_NAME="${1:?missing app name}"
DEPLOY_DIR="DEPLOY_DIR_PLACEHOLDER"
APP_ROOT="${DEPLOY_DIR}/${APP_NAME}"
APP_BINARY="${APP_ROOT}/${APP_NAME}"
MIGRATE_SCRIPT="${DEPLOY_DIR}/scripts/run-db-migrations.sh"

if [ ! -f "$APP_BINARY" ]; then
  echo "[start] 找不到二进制文件: ${APP_BINARY}" >&2
  exit 1
fi

export GOGC="${GOGC:-100}"
export GOMEMLIMIT="${GOMEMLIMIT:-512MiB}"
export GOMAXPROCS="${GOMAXPROCS:-2}"

if [ "$APP_NAME" = "upload" ]; then
  export UPLOAD_LOCAL_ROOT="${UPLOAD_LOCAL_ROOT:-${DEPLOY_DIR}/upload}"
fi

chmod +x "$APP_BINARY" 2>/dev/null || true

if [ -f "$MIGRATE_SCRIPT" ]; then
  /bin/sh "$MIGRATE_SCRIPT"
fi

cd "$APP_ROOT"
if command -v nice >/dev/null 2>&1; then
  exec nice -n 5 "$APP_BINARY"
fi
exec "$APP_BINARY"
STARTEOF

sed -i "s|DEPLOY_DIR_PLACEHOLDER|${DEPLOY_DIR}|g" "$DEPLOY_DIR/scripts/start-service.sh"
chmod 755 "$DEPLOY_DIR/scripts/start-service.sh"

# ---------- 5. 生成迁移脚本（双数据库） ----------
info "生成迁移脚本..."

DB_BASEADMIN_HOST_B64="$(b64 "$DB_HOST")"
DB_BASEADMIN_PORT_B64="$(b64 "$DB_PORT")"
DB_BASEADMIN_NAME_B64="$(b64 "$DB_BASEADMIN_NAME")"
DB_BASEADMIN_USER_B64="$(b64 "$DB_BASEADMIN_USER")"
DB_BASEADMIN_PASS_B64="$(b64 "$DB_BASEADMIN_PASS")"

DB_FUNDDISK_NAME_B64="$(b64 "$DB_FUNDDISK_NAME")"
DB_FUNDDISK_USER_B64="$(b64 "$DB_FUNDDISK_USER")"
DB_FUNDDISK_PASS_B64="$(b64 "$DB_FUNDDISK_PASS")"

cat > "$DEPLOY_DIR/scripts/run-db-migrations.sh" <<EOF
#!/bin/sh
set -eu

DEPLOY_DIR="${DEPLOY_DIR}"
MIGRATIONS_DIR="\${DEPLOY_DIR}/database/migrations"
LOCK_DIR="\${DEPLOY_DIR}/scripts/.db-migrate.lock"
LOCK_TIMEOUT_SECONDS=120

decode_b64() {
  if [ -z "\$1" ]; then return 0; fi
  printf '%s' "\$1" | base64 -d
}

# baseadmin 库（菜单数据）
BA_HOST="\$(decode_b64 '${DB_BASEADMIN_HOST_B64}')"
BA_PORT="\$(decode_b64 '${DB_BASEADMIN_PORT_B64}')"
BA_NAME="\$(decode_b64 '${DB_BASEADMIN_NAME_B64}')"
BA_USER="\$(decode_b64 '${DB_BASEADMIN_USER_B64}')"
BA_PASS="\$(decode_b64 '${DB_BASEADMIN_PASS_B64}')"

# funddisk 库（业务表）
FD_HOST="\$(decode_b64 '${DB_BASEADMIN_HOST_B64}')"
FD_PORT="\$(decode_b64 '${DB_BASEADMIN_PORT_B64}')"
FD_NAME="\$(decode_b64 '${DB_FUNDDISK_NAME_B64}')"
FD_USER="\$(decode_b64 '${DB_FUNDDISK_USER_B64}')"
FD_PASS="\$(decode_b64 '${DB_FUNDDISK_PASS_B64}')"

find_mysql_bin() {
  for candidate in mysql /www/server/mysql/bin/mysql /www/server/mariadb/bin/mysql /usr/bin/mysql /usr/local/mysql/bin/mysql; do
    if command -v "\$candidate" >/dev/null 2>&1 || [ -x "\$candidate" ]; then
      printf '%s\n' "\$candidate"
      return 0
    fi
  done
  echo "[migrate] mysql 客户端未找到" >&2
  return 1
}

MYSQL_BIN="\$(find_mysql_bin)"

mysql_baseadmin() {
  "\$MYSQL_BIN" --host="\$BA_HOST" --port="\$BA_PORT" --user="\$BA_USER" --password="\$BA_PASS" --default-character-set=utf8mb4 "\$BA_NAME" "\$@"
}

mysql_funddisk() {
  "\$MYSQL_BIN" --host="\$FD_HOST" --port="\$FD_PORT" --user="\$FD_USER" --password="\$FD_PASS" --default-character-set=utf8mb4 "\$FD_NAME" "\$@"
}

cleanup_lock() {
  rm -f "\${LOCK_DIR}/pid" 2>/dev/null || true
  rmdir "\$LOCK_DIR" 2>/dev/null || true
}

acquire_lock() {
  wait_seconds=0
  while ! mkdir "\$LOCK_DIR" 2>/dev/null; do
    if [ -f "\${LOCK_DIR}/pid" ]; then
      stale_pid="\$(cat "\${LOCK_DIR}/pid" 2>/dev/null || true)"
      if [ -n "\$stale_pid" ] && ! kill -0 "\$stale_pid" 2>/dev/null; then
        rm -rf "\$LOCK_DIR"
        continue
      fi
    fi
    if [ "\$wait_seconds" -ge "\$LOCK_TIMEOUT_SECONDS" ]; then
      echo "[migrate] 等待迁移锁超时: \$LOCK_DIR" >&2
      return 1
    fi
    wait_seconds=\$((wait_seconds + 1))
    sleep 1
  done
  printf '%s' "\$\$" > "\${LOCK_DIR}/pid"
}

ensure_state_table() {
  mysql_baseadmin -e "CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT NOT NULL PRIMARY KEY, dirty BOOLEAN NOT NULL);"
}

load_state() {
  state_row="\$(mysql_baseadmin -Nse "SELECT version, dirty FROM schema_migrations LIMIT 1;" 2>/dev/null || true)"
  if [ -z "\$state_row" ]; then CURRENT_VERSION=0; CURRENT_DIRTY=0; return 0; fi
  set -- \$state_row
  CURRENT_VERSION="\${1:-0}"
  CURRENT_DIRTY="\${2:-0}"
}

mark_dirty() {
  mysql_baseadmin -e "DELETE FROM schema_migrations; INSERT INTO schema_migrations (version, dirty) VALUES (\$1, TRUE);"
}

mark_clean() {
  mysql_baseadmin -e "UPDATE schema_migrations SET version = \$1, dirty = FALSE;"
}

apply_pending_migrations() {
  if [ ! -d "\$MIGRATIONS_DIR" ]; then
    echo "[migrate] 迁移目录不存在，跳过: \${MIGRATIONS_DIR}"
    return 0
  fi
  set -- "\$MIGRATIONS_DIR"/*.up.sql
  if [ ! -e "\$1" ]; then echo "[migrate] 无迁移文件: \${MIGRATIONS_DIR}"; return 0; fi
  applied_count=0
  for migration_file in "\$@"; do
    migration_name="\$(basename "\$migration_file")"
    migration_version="\${migration_name%%_*}"
    case "\$migration_version" in ''|*[!0-9]*) echo "[migrate] 跳过无效文件名: \${migration_name}" >&2; continue ;; esac
    migration_version_num="\$(expr "\$migration_version" + 0)"
    if [ "\$migration_version_num" -le "\$CURRENT_VERSION" ]; then continue; fi
    echo "[migrate] 执行 \${migration_name}"
    mark_dirty "\$migration_version_num"
    # 路由规则：
    #   - 含 _admin_menus 或 _system_*.sql → baseadmin 库（系统级 / 菜单）
    #   - 其它 _member_* → funddisk 库（业务表）
    #   - 其它（兼容历史） → baseadmin 库
    case "\$migration_name" in
      *_member_admin_menus*) target_db="baseadmin" ;;
      *_member_*)            target_db="funddisk"  ;;
      *)                     target_db="baseadmin" ;;
    esac
    if [ "\$target_db" = "funddisk" ]; then
      echo "[migrate]   → funddisk 库"
      mysql_funddisk < "\$migration_file"
    else
      echo "[migrate]   → baseadmin 库"
      mysql_baseadmin < "\$migration_file"
    fi
    mark_clean "\$migration_version_num"
    CURRENT_VERSION="\$migration_version_num"
    CURRENT_DIRTY=0
    applied_count=\$((applied_count + 1))
  done
  if [ "\$applied_count" -eq 0 ]; then echo "[migrate] 无新迁移需要执行"; else echo "[migrate] 执行了 \${applied_count} 个迁移，当前版本=\${CURRENT_VERSION}"; fi
}

trap cleanup_lock EXIT INT TERM
acquire_lock
ensure_state_table
load_state
if [ "\$CURRENT_DIRTY" != "0" ]; then
  echo "[migrate] schema_migrations 在版本 \${CURRENT_VERSION} 处于 dirty 状态，需要手动处理" >&2
  exit 1
fi
apply_pending_migrations
EOF

chmod 755 "$DEPLOY_DIR/scripts/run-db-migrations.sh"

# ---------- 6. 数据库迁移 ----------
if [ "$SKIP_MIGRATE" = "1" ]; then
  warn "跳过数据库迁移"
else
  info "执行数据库迁移..."
  MYSQL_BIN="$(find_mysql_bin)" || die "找不到 mysql 客户端"
  # 使用 system 配置执行迁移（迁移文件走 baseadmin 库）
  if [ -f "$DEPLOY_DIR/migrate" ]; then
    "$DEPLOY_DIR/migrate" --config "$DEPLOY_DIR/system/manifest/config/config.yaml" --dir "$DEPLOY_DIR/database/migrations" up || warn "migrate 工具执行失败，尝试脚本迁移"
  fi
  info "数据库迁移完成"
fi

# ---------- 7. Supervisor 配置 ----------
if [ "$SKIP_SUPERVISOR" = "1" ]; then
  warn "跳过 Supervisor 配置"
else
  info "配置 Supervisor..."

  SUPERVISORCTL_BIN="$(find_supervisorctl_bin)" || die "找不到 supervisorctl"
  PYTHON_BIN="$(find_python_bin)" || die "找不到 python"

  # 确保 supervisor 包含 profile 目录
  if [ -f "$SUPERVISOR_CONF" ]; then
    if ! grep -Fq "$BT_SUPERVISOR_PROFILE_DIR/*.ini" "$SUPERVISOR_CONF"; then
      cat >> "$SUPERVISOR_CONF" <<SUPEOF

[include]
files = ${BT_SUPERVISOR_PROFILE_DIR}/*.ini
SUPEOF
    fi
  fi

  mkdir -p "$BT_SUPERVISOR_PROFILE_DIR" "$BT_SUPERVISOR_LOG_DIR"

  # 停止旧的同名服务
  for app in "${APPS[@]}"; do
    program="${SERVICE_PREFIX}-${app}"
    "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" stop "${program}:" >/dev/null 2>&1 || true
  done

  # 生成每个服务的 ini
  for i in "${!APPS[@]}"; do
    app="${APPS[$i]}"
    program="${SERVICE_PREFIX}-${app}"
    app_root="$DEPLOY_DIR/$app"
    ini_file="$BT_SUPERVISOR_PROFILE_DIR/${program}.ini"
    environment="GOGC=\"${GOGC}\",GOMEMLIMIT=\"${GOMEMLIMIT}\",GOMAXPROCS=\"${GOMAXPROCS}\""
    if [ "$app" = "upload" ]; then
      environment="${environment},UPLOAD_LOCAL_ROOT=\"${DEPLOY_DIR}/upload\""
    fi

    cat > "$ini_file" <<INIEOF
[program:${program}]
directory=${app_root}
command=/bin/sh ${DEPLOY_DIR}/scripts/start-service.sh ${app}
autostart=true
autorestart=true
startsecs=5
startretries=3
stopasgroup=true
killasgroup=true
stopsignal=QUIT
stdout_logfile=${BT_SUPERVISOR_LOG_DIR}/${program}.out.log
stderr_logfile=${BT_SUPERVISOR_LOG_DIR}/${program}.err.log
stdout_logfile_maxbytes=20MB
stderr_logfile_maxbytes=20MB
stdout_logfile_backups=3
stderr_logfile_backups=3
user=root
priority=999
environment=${environment}
numprocs=1
process_name=%(program_name)s_%(process_num)02d
INIEOF
    info "Supervisor 配置: $program → 端口 ${PORTS[$i]}"
  done

  # 更新宝塔 Supervisor config.json
  "$PYTHON_BIN" - "$BT_SUPERVISOR_CONFIG_JSON" "$SERVICE_PREFIX" "$DEPLOY_DIR" <<'PYEOF'
import json
import os
import sys

path, prefix, deploy_dir = sys.argv[1:4]
try:
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)
except Exception:
    data = []

apps = ["system", "upload", "member"]
commands = {app: f"/bin/sh {deploy_dir}/scripts/start-service.sh {app}" for app in apps}
programs = {app: f"{prefix}-{app}" for app in apps}
data = [
    item for item in data
    if item.get("program") not in programs.values()
    and item.get("command") not in commands.values()
]
for app in apps:
    data.append({
        "program": programs[app],
        "command": commands[app],
        "directory": f"{deploy_dir}/{app}",
        "user": "root",
        "priority": "999",
        "numprocs": "1",
        "runStatus": "",
        "ps": programs[app],
    })

os.makedirs(os.path.dirname(path), exist_ok=True)
with open(path, "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False)
PYEOF

  # 重载 Supervisor
  "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" reread || true
  "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" update || true
  for app in "${APPS[@]}"; do
    program="${SERVICE_PREFIX}-${app}"
    "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" restart "${program}:" >/dev/null 2>&1 || {
      "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" stop "${program}:" >/dev/null 2>&1 || true
      "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" start "${program}:" >/dev/null 2>&1 || warn "启动失败: $program"
    }
  done

  info "Supervisor 状态:"
  "$SUPERVISORCTL_BIN" -c "$SUPERVISOR_CONF" status | grep -E "^${SERVICE_PREFIX}-" || true
fi

# ---------- 8. Nginx 配置 ----------
if [ "$SKIP_NGINX" = "1" ]; then
  warn "跳过 Nginx 配置"
else
  info "配置 Nginx..."

  NGINX_BIN="$(find_nginx_bin)" || die "找不到 nginx"
  mkdir -p "$BT_NGINX_DIR"

  # 检测 SSL
  SSL_ENABLED="0"
  if [ "$SSL_MODE" = "auto" ]; then
    if [ -f "/www/server/panel/vhost/cert/${DOMAIN}/fullchain.pem" ] && [ -f "/www/server/panel/vhost/cert/${DOMAIN}/privkey.pem" ]; then
      SSL_ENABLED="1"
    fi
  fi

  NGINX_CONF="$BT_NGINX_DIR/${DOMAIN}.conf"
  if [ -f "$NGINX_CONF" ]; then
    cp -a "$NGINX_CONF" "${NGINX_CONF}.bak.$(date +%Y%m%d%H%M%S)"
    info "已备份旧 Nginx 配置"
  fi

  if [ "$SSL_ENABLED" = "1" ]; then
    cat > "$NGINX_CONF" <<NGXEOF
server {
    listen 80;
    listen [::]:80;
    server_name ${DOMAIN};
    return 301 https://\$host\$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name ${DOMAIN};

    root ${DEPLOY_DIR};
    index index.html;
    client_max_body_size 50m;

    ssl_certificate /www/server/panel/vhost/cert/${DOMAIN}/fullchain.pem;
    ssl_certificate_key /www/server/panel/vhost/cert/${DOMAIN}/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers on;

    # 管理端入口
    location = /admin {
        return 301 /admin/;
    }

    location = /admin/index.html {
        alias ${DEPLOY_DIR}/admin/index.html;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }

    location ^~ /admin/ {
        alias ${DEPLOY_DIR}/admin/;
        try_files \$uri \$uri/ /admin/index.html;
    }

    # 静态资源缓存
    location ~* \.(js|css|woff2?|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public";
        access_log off;
    }

    # 上传文件静态访问
    location ^~ /upload/ {
        alias ${DEPLOY_DIR}/upload/resource/upload/;
        expires 30d;
        add_header Cache-Control "public";
        access_log off;
    }

    # member 后台 API（admin 用）
    location ^~ /api/member/ {
        proxy_pass http://127.0.0.1:${PORTS[2]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # member portal API（H5 会员端用，路径不能合并到 /api/member/）
    location ^~ /api/member-portal/ {
        proxy_pass http://127.0.0.1:${PORTS[2]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # upload API
    location ^~ /api/upload/ {
        proxy_pass http://127.0.0.1:${PORTS[1]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # system API（兜底，匹配其余 /api/ 请求）
    location ^~ /api/ {
        proxy_pass http://127.0.0.1:${PORTS[0]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 前端 Vue Router history 模式
    location / {
        try_files \$uri \$uri/ /index.html;
        location = /index.html {
            add_header Cache-Control "no-cache, no-store, must-revalidate";
        }
    }

    # 安全 Header
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
}
NGXEOF
  else
    cat > "$NGINX_CONF" <<NGXEOF
server {
    listen 80;
    listen [::]:80;
    server_name ${DOMAIN};

    root ${DEPLOY_DIR};
    index index.html;
    client_max_body_size 50m;

    location = /admin {
        return 301 /admin/;
    }

    location = /admin/index.html {
        alias ${DEPLOY_DIR}/admin/index.html;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }

    location ^~ /admin/ {
        alias ${DEPLOY_DIR}/admin/;
        try_files \$uri \$uri/ /admin/index.html;
    }

    location ~* \.(js|css|woff2?|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public";
        access_log off;
    }

    location ^~ /upload/ {
        alias ${DEPLOY_DIR}/upload/resource/upload/;
        expires 30d;
        add_header Cache-Control "public";
        access_log off;
    }

    location ^~ /api/member/ {
        proxy_pass http://127.0.0.1:${PORTS[2]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    location ^~ /api/upload/ {
        proxy_pass http://127.0.0.1:${PORTS[1]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    location ^~ /api/ {
        proxy_pass http://127.0.0.1:${PORTS[0]};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 10s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    location / {
        try_files \$uri \$uri/ /index.html;
        location = /index.html {
            add_header Cache-Control "no-cache, no-store, must-revalidate";
        }
    }

    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
}
NGXEOF
  fi

  "$NGINX_BIN" -t && "$NGINX_BIN" -s reload
  info "Nginx 配置已生效"
fi

# ---------- 9. 完成 ----------
echo ""
info "========================================="
info "部署完成！"
info "========================================="
info "部署目录: $DEPLOY_DIR"
info "服务管理（Supervisor）:"
for app in "${APPS[@]}"; do
  info "  supervisorctl status ${SERVICE_PREFIX}-${app}"
  info "  supervisorctl restart ${SERVICE_PREFIX}-${app}:"
done
info ""
info "前端构建后将 dist 内容复制到 $DEPLOY_DIR:"
info "  cp -rf vue-vben-admin/apps/web-antd/dist/* $DEPLOY_DIR/"
info ""
info "管理后台地址: https://${DOMAIN}/admin/"
info "========================================="
