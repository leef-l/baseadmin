package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	supervisorProfileDir = "/www/server/panel/plugin/supervisor/profile"
	supervisorLogDir     = "/www/server/panel/plugin/supervisor/log"
	baoTaSupervisorConf  = "/www/server/panel/plugin/supervisor/config.json"
)

type supervisorProgram struct {
	Name         string
	Program      string
	Command      string
	Directory    string
	RunUser      string
	Numprocs     int
	Priority     int
	Autostart    int
	Autorestart  int
	Startsecs    int
	Startretries int
	StopSignal   string
	Environment  string
	Remark       string
}

type supervisorRuntime struct {
	RunStatus  string
	Pid        string
	Uptime     string
	StatusText string
}

type baoTaSupervisorItem struct {
	Program   string `json:"program"`
	Command   string `json:"command"`
	Directory string `json:"directory"`
	User      string `json:"user"`
	Priority  string `json:"priority"`
	Numprocs  string `json:"numprocs"`
	Ps        string `json:"ps"`
}

func writeSupervisorProgram(item supervisorProgram) error {
	if err := os.MkdirAll(supervisorProfileDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(supervisorLogDir, 0o755); err != nil {
		return err
	}
	content := buildSupervisorProgramContent(item)
	return os.WriteFile(supervisorConfigPath(item.Program), []byte(content), 0o644)
}

func buildSupervisorProgramContent(item supervisorProgram) string {
	var b strings.Builder
	fmt.Fprintf(&b, "[program:%s]\n", item.Program)
	fmt.Fprintf(&b, "directory=%s\n", item.Directory)
	fmt.Fprintf(&b, "command=%s\n", item.Command)
	fmt.Fprintf(&b, "autostart=%s\n", supervisorBool(item.Autostart))
	fmt.Fprintf(&b, "autorestart=%s\n", supervisorBool(item.Autorestart))
	fmt.Fprintf(&b, "startsecs=%d\n", item.Startsecs)
	fmt.Fprintf(&b, "startretries=%d\n", item.Startretries)
	fmt.Fprintf(&b, "stopasgroup=true\n")
	fmt.Fprintf(&b, "killasgroup=true\n")
	fmt.Fprintf(&b, "stopsignal=%s\n", item.StopSignal)
	fmt.Fprintf(&b, "stdout_logfile=%s\n", supervisorLogPath(item.Program, "normal"))
	fmt.Fprintf(&b, "stderr_logfile=%s\n", supervisorLogPath(item.Program, "error"))
	fmt.Fprintf(&b, "stdout_logfile_maxbytes=20MB\n")
	fmt.Fprintf(&b, "stderr_logfile_maxbytes=20MB\n")
	fmt.Fprintf(&b, "stdout_logfile_backups=5\n")
	fmt.Fprintf(&b, "stderr_logfile_backups=5\n")
	fmt.Fprintf(&b, "user=%s\n", item.RunUser)
	fmt.Fprintf(&b, "priority=%d\n", item.Priority)
	if item.Environment != "" {
		fmt.Fprintf(&b, "environment=%s\n", item.Environment)
	}
	fmt.Fprintf(&b, "numprocs=%d\n", item.Numprocs)
	fmt.Fprintf(&b, "process_name=%%(program_name)s_%%(process_num)02d\n")
	return b.String()
}

func supervisorBool(value int) string {
	if value == 1 {
		return "true"
	}
	return "false"
}

func supervisorConfigExists(program string) bool {
	return fileExists(supervisorConfigPath(program))
}

func supervisorConfigPath(program string) string {
	return filepath.Join(supervisorProfileDir, program+".ini")
}

func supervisorLogPath(program, logType string) string {
	suffix := ".out.log"
	if normalizeLogType(logType) == "error" {
		suffix = ".err.log"
	}
	return filepath.Join(supervisorLogDir, program+suffix)
}

func updateSupervisor(ctx context.Context) error {
	out, err := runSupervisorctl(ctx, "update")
	if err == nil && !strings.Contains(strings.ToLower(out), "error") {
		return nil
	}
	if err != nil {
		return gerror.Newf("宝塔Supervisor更新失败: %s", trimOutput(out, err.Error()))
	}
	return gerror.Newf("宝塔Supervisor更新失败: %s", trimOutput(out, ""))
}

func restartSupervisorProgram(ctx context.Context, program string) error {
	_ = stopSupervisorProgram(ctx, program)
	if err := startSupervisorProgram(ctx, program); err != nil {
		return err
	}
	return nil
}

func startSupervisorProgram(ctx context.Context, program string) error {
	out, err := runSupervisorctl(ctx, "start", program+":")
	if err == nil || strings.Contains(out, "already started") {
		return nil
	}
	return gerror.Newf("守护进程启动失败: %s", trimOutput(out, err.Error()))
}

func stopSupervisorProgram(ctx context.Context, program string) error {
	out, err := runSupervisorctl(ctx, "stop", program+":")
	if err == nil || strings.Contains(out, "not running") || strings.Contains(out, "no such group") {
		return nil
	}
	return gerror.Newf("守护进程暂停失败: %s", trimOutput(out, err.Error()))
}

func runSupervisorctl(ctx context.Context, args ...string) (string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cmdCtx, supervisorctlBinary(), args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func supervisorctlBinary() string {
	if fileExists("/www/server/panel/pyenv/bin/supervisorctl") {
		return "/www/server/panel/pyenv/bin/supervisorctl"
	}
	return "/usr/bin/supervisorctl"
}

func loadSupervisorRuntimeMap(ctx context.Context, programs []string) map[string]supervisorRuntime {
	out := make(map[string]supervisorRuntime, len(programs))
	if len(programs) == 0 {
		return out
	}
	for _, program := range programs {
		if program != "" {
			out[program] = supervisorRuntime{}
		}
	}
	statusOutput, err := runSupervisorctl(ctx, "status")
	if err != nil && strings.TrimSpace(statusOutput) == "" {
		for program := range out {
			out[program] = supervisorRuntime{RunStatus: "UNKNOWN", StatusText: "Supervisor状态获取失败"}
		}
		return out
	}
	for _, line := range strings.Split(statusOutput, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		program := supervisorProgramFromStatusName(fields[0])
		if _, ok := out[program]; !ok {
			continue
		}
		current := parseSupervisorRuntime(fields[1], strings.Join(fields[2:], " "))
		out[program] = mergeSupervisorRuntime(out[program], current)
	}
	return out
}

func supervisorProgramFromStatusName(value string) string {
	if index := strings.Index(value, ":"); index > 0 {
		return value[:index]
	}
	if index := strings.LastIndex(value, "_"); index > 0 {
		return value[:index]
	}
	return value
}

func parseSupervisorRuntime(status, detail string) supervisorRuntime {
	status = strings.TrimSpace(status)
	detail = strings.TrimSpace(detail)
	runtime := supervisorRuntime{
		RunStatus:  status,
		StatusText: statusDisplayText(status, detail),
	}
	if status == "RUNNING" {
		if pid := parsePid(detail); pid != "" {
			runtime.Pid = pid
		}
		if uptime := parseUptime(detail); uptime != "" {
			runtime.Uptime = uptime
		}
	}
	return runtime
}

func mergeSupervisorRuntime(old, next supervisorRuntime) supervisorRuntime {
	if old.RunStatus == "" {
		return next
	}
	if old.RunStatus == "RUNNING" {
		return old
	}
	if next.RunStatus == "RUNNING" {
		return next
	}
	if old.StatusText == "" && next.StatusText != "" {
		return next
	}
	return old
}

func parsePid(detail string) string {
	index := strings.Index(detail, "pid ")
	if index < 0 {
		return ""
	}
	rest := detail[index+4:]
	if comma := strings.Index(rest, ","); comma >= 0 {
		rest = rest[:comma]
	}
	return strings.TrimSpace(rest)
}

func parseUptime(detail string) string {
	index := strings.Index(detail, "uptime ")
	if index < 0 {
		return ""
	}
	return strings.TrimSpace(detail[index+7:])
}

func statusDisplayText(status, detail string) string {
	switch status {
	case "RUNNING":
		return "运行中"
	case "STOPPED":
		return "已暂停"
	case "STARTING":
		return "启动中"
	case "STOPPING":
		return "暂停中"
	case "BACKOFF":
		return "启动失败重试中"
	case "FATAL":
		if detail != "" {
			return "异常: " + detail
		}
		return "异常"
	case "EXITED":
		return "已退出"
	default:
		if status == "" {
			return "未知"
		}
		return status
	}
}

func removeSupervisorProgram(program string) {
	_ = os.Remove(supervisorConfigPath(program))
}

func removeSupervisorLogs(program string) {
	_ = os.Remove(supervisorLogPath(program, "normal"))
	_ = os.Remove(supervisorLogPath(program, "error"))
}

func upsertBaoTaSupervisorConfig(item supervisorProgram) error {
	list, err := readBaoTaSupervisorConfig()
	if err != nil {
		return err
	}
	next := baoTaSupervisorItem{
		Program:   item.Program,
		Command:   item.Command,
		Directory: item.Directory,
		User:      item.RunUser,
		Priority:  strconv.Itoa(item.Priority),
		Numprocs:  strconv.Itoa(item.Numprocs),
		Ps:        item.Remark,
	}
	replaced := false
	for index := range list {
		if list[index].Program == item.Program {
			list[index] = next
			replaced = true
			break
		}
	}
	if !replaced {
		list = append(list, next)
	}
	return writeBaoTaSupervisorConfig(list)
}

func removeBaoTaSupervisorConfig(program string) {
	list, err := readBaoTaSupervisorConfig()
	if err != nil {
		return
	}
	filtered := make([]baoTaSupervisorItem, 0, len(list))
	for _, item := range list {
		if item.Program != program {
			filtered = append(filtered, item)
		}
	}
	_ = writeBaoTaSupervisorConfig(filtered)
}

func readBaoTaSupervisorConfig() ([]baoTaSupervisorItem, error) {
	data, err := os.ReadFile(baoTaSupervisorConf)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, nil
	}
	var list []baoTaSupervisorItem
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func writeBaoTaSupervisorConfig(list []baoTaSupervisorItem) error {
	if err := os.MkdirAll(filepath.Dir(baoTaSupervisorConf), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(baoTaSupervisorConf, data, 0o644)
}

func readSupervisorLog(program, logType string, lines int) (string, error) {
	if lines <= 0 || lines > 2000 {
		lines = 500
	}
	path := supervisorLogPath(program, logType)
	if !fileExists(path) {
		return "", nil
	}
	cmd := exec.Command("tail", "-n", strconv.Itoa(lines), path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", gerror.Newf("日志读取失败: %s", trimOutput(string(out), err.Error()))
	}
	return string(out), nil
}

func normalizeLogType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "error" || value == "err" {
		return "error"
	}
	return "normal"
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func systemUserExists(username string) bool {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return true
	}
	prefix := username + ":"
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func trimOutput(out string, fallback string) string {
	out = strings.TrimSpace(out)
	if out == "" {
		out = strings.TrimSpace(fallback)
	}
	if len(out) > 1200 {
		out = out[:1200]
	}
	return out
}
