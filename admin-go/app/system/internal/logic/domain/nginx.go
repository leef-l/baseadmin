package domain

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/system/internal/model"
)

const (
	managedConfigPrefix = "baseadmin-domain-"
	nginxVhostDir       = "/www/server/panel/vhost/nginx"
	nginxCertDir        = "/www/server/panel/vhost/cert"
	siteRoot            = "/www/wwwroot/baseadmin.easytestdev.online"
	uploadAliasRoot     = "/www/wwwroot/baseadmin.easytestdev.online/upload/resource/upload/"
)

func applyNginxConfig(ctx context.Context, row *domainRow) (*model.DomainApplyNginxOutput, error) {
	if row == nil || row.Domain == "" {
		return nil, gerror.New("域名不存在或已删除")
	}
	if err := ensureNginxDomainAvailable(row.Domain); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(nginxVhostDir, 0o755); err != nil {
		return nil, err
	}

	configPath := managedConfigPath(row.Domain)
	fullchain, privkey, hasCert := domainCertPaths(row.Domain)
	content := buildNginxConfig(row.Domain, fullchain, privkey, hasCert)

	previous, readErr := os.ReadFile(configPath)
	hadPrevious := readErr == nil
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		return nil, err
	}
	if err := runNginxTest(ctx); err != nil {
		restoreNginxConfig(configPath, previous, hadPrevious)
		return nil, err
	}
	if err := reloadNginx(ctx); err != nil {
		return nil, err
	}

	sslStatus := 0
	if hasCert {
		sslStatus = 1
	}
	return &model.DomainApplyNginxOutput{
		ConfigPath:  configPath,
		NginxStatus: 1,
		SslStatus:   sslStatus,
	}, nil
}

func managedConfigPath(domainName string) string {
	return filepath.Join(nginxVhostDir, managedConfigPrefix+domainName+".conf")
}

func domainCertPaths(domainName string) (fullchain string, privkey string, ok bool) {
	fullchain = filepath.Join(nginxCertDir, domainName, "fullchain.pem")
	privkey = filepath.Join(nginxCertDir, domainName, "privkey.pem")
	if fileExists(fullchain) && fileExists(privkey) {
		return fullchain, privkey, true
	}
	return fullchain, privkey, false
}

func buildNginxConfig(domainName, fullchain, privkey string, hasCert bool) string {
	blocks := []string{buildNginxServerBlock(domainName, "", "", false)}
	if hasCert {
		blocks = append(blocks, buildNginxServerBlock(domainName, fullchain, privkey, true))
	}
	return strings.Join(blocks, "\n")
}

func buildNginxServerBlock(domainName, fullchain, privkey string, ssl bool) string {
	listen := "listen 80;"
	sslLines := ""
	if ssl {
		listen = "listen 443 ssl http2;"
		sslLines = fmt.Sprintf(`
    ssl_certificate %s;
    ssl_certificate_key %s;
    ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
    ssl_ciphers EECDH+CHACHA20:EECDH+CHACHA20-draft:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    add_header Strict-Transport-Security "max-age=31536000";
    error_page 497 https://$host$request_uri;
`, fullchain, privkey)
	}
	return fmt.Sprintf(`# managed by baseadmin system_domain
server {
    %s
    server_name %s;
    root %s;
    index index.html index.htm;
%s
    location = / {
        return 302 /admin/;
    }

    location /admin/ {
        try_files $uri $uri/ /admin/index.html;
    }

    location /api/system/ {
        proxy_pass http://127.0.0.1:10022/api/system/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api/upload/ {
        proxy_pass http://127.0.0.1:10023/api/upload/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location ^~ /upload/ {
        alias %s;
        expires 30d;
        add_header Cache-Control "public";
        access_log off;
    }

    location ~ ^/(\.user.ini|\.htaccess|\.git|\.env|\.svn|\.project|LICENSE|README.md) {
        return 404;
    }

    location ~ \.well-known {
        allow all;
    }

    access_log /www/wwwlogs/baseadmin-domain-%s.log;
    error_log /www/wwwlogs/baseadmin-domain-%s.error.log;
}
`, listen, domainName, siteRoot, sslLines, uploadAliasRoot, domainName, domainName)
}

func ensureNginxDomainAvailable(domainName string) error {
	targetPath := filepath.Clean(managedConfigPath(domainName))
	entries, err := os.ReadDir(nginxVhostDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".conf") {
			continue
		}
		path := filepath.Join(nginxVhostDir, entry.Name())
		if filepath.Clean(path) == targetPath {
			continue
		}
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if configDeclaresServerName(string(content), domainName) {
			return gerror.Newf("域名已在Nginx配置中存在: %s", entry.Name())
		}
	}
	return nil
}

func configDeclaresServerName(content, domainName string) bool {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "server_name ") {
			continue
		}
		line = strings.TrimSuffix(strings.TrimPrefix(line, "server_name "), ";")
		for _, item := range strings.Fields(line) {
			if strings.TrimSpace(item) == domainName {
				return true
			}
		}
	}
	return false
}

func runNginxTest(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, nginxBinary(), "-t").CombinedOutput()
	if err == nil {
		return nil
	}
	return gerror.Newf("Nginx配置检测失败: %s", trimCommandOutput(out))
}

func reloadNginx(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, nginxBinary(), "-s", "reload").CombinedOutput()
	if err == nil {
		return nil
	}
	return gerror.Newf("Nginx重载失败: %s", trimCommandOutput(out))
}

func nginxBinary() string {
	if fileExists("/usr/bin/nginx") {
		return "/usr/bin/nginx"
	}
	return "/www/server/nginx/sbin/nginx"
}

func restoreNginxConfig(path string, previous []byte, hadPrevious bool) {
	if hadPrevious {
		_ = os.WriteFile(path, previous, 0o644)
		return
	}
	_ = os.Remove(path)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func trimCommandOutput(out []byte) string {
	value := strings.TrimSpace(string(out))
	if len(value) > 1200 {
		value = value[:1200]
	}
	return value
}
