package domain

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

const baoTaApplySSLScript = `
import json
import os
import sys
import traceback

domain = sys.argv[1]
site_root = sys.argv[2]
panel_path = "/www/server/panel"
os.chdir(panel_path)
if "class/" not in sys.path:
    sys.path.insert(0, "class/")

try:
    import public
    from acme_v2 import acme_v2

    acme_dir = os.path.join(site_root, ".well-known", "acme-challenge")
    os.makedirs(acme_dir, exist_ok=True)
    try:
        public.set_own(acme_dir, "www")
    except Exception:
        pass

    result = acme_v2().apply_cert([domain], "http", site_root)
    if not isinstance(result, dict) or not result.get("status"):
        print(json.dumps({
            "status": False,
            "msg": result.get("msg") if isinstance(result, dict) else str(result),
        }, ensure_ascii=False))
        sys.exit(2)

    fullchain = (result.get("cert") or "") + (result.get("root") or "")
    private_key = result.get("private_key") or ""
    save_path = result.get("save_path") or ""
    if save_path and not save_path.startswith("/"):
        save_path = os.path.join(panel_path, save_path)
    if (not fullchain or not private_key) and save_path:
        try:
            with open(os.path.join(save_path, "fullchain.pem"), "r", encoding="utf-8") as f:
                fullchain = f.read()
            with open(os.path.join(save_path, "privkey.pem"), "r", encoding="utf-8") as f:
                private_key = f.read()
        except Exception:
            pass
    if not fullchain or not private_key:
        print(json.dumps({"status": False, "msg": "宝塔ACME申请成功但未返回证书内容"}, ensure_ascii=False))
        sys.exit(3)

    cert_dir = os.path.join("/www/server/panel/vhost/cert", domain)
    os.makedirs(cert_dir, exist_ok=True)
    fullchain_path = os.path.join(cert_dir, "fullchain.pem")
    privkey_path = os.path.join(cert_dir, "privkey.pem")
    with open(fullchain_path, "w", encoding="utf-8") as f:
        f.write(fullchain)
    with open(privkey_path, "w", encoding="utf-8") as f:
        f.write(private_key)
    os.chmod(fullchain_path, 0o600)
    os.chmod(privkey_path, 0o600)

    print(json.dumps({
        "status": True,
        "cert_path": cert_dir,
        "source_path": save_path,
    }, ensure_ascii=False))
except Exception as exc:
    print(json.dumps({
        "status": False,
        "msg": str(exc),
        "trace": traceback.format_exc(limit=3),
    }, ensure_ascii=False))
    sys.exit(1)
`

type baoTaSSLResult struct {
	Status     bool   `json:"status"`
	Message    any    `json:"msg"`
	CertPath   string `json:"cert_path"`
	SourcePath string `json:"source_path"`
	Trace      string `json:"trace"`
}

func applyBaoTaSSL(ctx context.Context, domainName string) (string, error) {
	if _, _, ok := domainCertPaths(domainName); ok {
		return filepath.Join(nginxCertDir, domainName), nil
	}

	sslCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(sslCtx, panelPythonBinary(), "-c", baoTaApplySSLScript, domainName, siteRoot)
	out, err := cmd.CombinedOutput()
	result, parseErr := parseBaoTaSSLResult(out)
	if parseErr != nil {
		return "", parseErr
	}
	if err != nil || !result.Status {
		message := formatBaoTaSSLMessage(result)
		if message == "" {
			message = trimCommandOutput(out)
		}
		if message == "" && err != nil {
			message = err.Error()
		}
		return "", gerror.Newf("宝塔SSL证书申请失败: %s", message)
	}
	if result.CertPath == "" {
		return "", gerror.New("宝塔SSL证书申请成功但证书路径为空")
	}
	return result.CertPath, nil
}

func parseBaoTaSSLResult(out []byte) (*baoTaSSLResult, error) {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
			continue
		}
		var result baoTaSSLResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}
		return &result, nil
	}
	return nil, gerror.Newf("无法解析宝塔SSL申请结果: %s", trimCommandOutput(out))
}

func formatBaoTaSSLMessage(result *baoTaSSLResult) string {
	if result == nil {
		return ""
	}
	switch value := result.Message.(type) {
	case string:
		return strings.TrimSpace(value)
	case []any, map[string]any:
		data, _ := json.Marshal(value)
		return strings.TrimSpace(string(data))
	default:
		if value == nil {
			return strings.TrimSpace(result.Trace)
		}
		data, _ := json.Marshal(value)
		return strings.TrimSpace(string(data))
	}
}

func panelPythonBinary() string {
	if fileExists("/www/server/panel/pyenv/bin/python3") {
		return "/www/server/panel/pyenv/bin/python3"
	}
	return "/usr/bin/python3"
}
