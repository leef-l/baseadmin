package daemon

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/snowflake"
)

const (
	defaultRunUser      = "root"
	defaultNumprocs     = 1
	defaultPriority     = 999
	defaultStartsecs    = 3
	defaultStartretries = 3
	defaultStopSignal   = "QUIT"
)

func init() {
	service.RegisterDaemon(New())
}

func New() *sDaemon {
	return &sDaemon{}
}

type sDaemon struct{}

func (s *sDaemon) Create(ctx context.Context, in *model.DaemonCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return err
	}
	normalizeDaemonCreateInput(in)
	if err := s.validateDaemonConfig(ctx, 0, in.Program, in.Name, in.Command, in.Directory, in.RunUser, in.Numprocs, in.Priority, in.Autostart, in.Autorestart, in.Startsecs, in.Startretries, in.StopSignal, in.Environment); err != nil {
		return err
	}
	if supervisorConfigExists(in.Program) {
		return gerror.New("宝塔Supervisor已存在同名守护进程")
	}

	row := daemonConfigFromCreateInput(in)
	if err := writeSupervisorProgram(row); err != nil {
		return err
	}
	if err := updateSupervisor(ctx); err != nil {
		removeSupervisorProgram(row.Program)
		return err
	}
	if err := upsertBaoTaSupervisorConfig(row); err != nil {
		_ = stopSupervisorProgram(ctx, in.Program)
		removeSupervisorProgram(in.Program)
		_ = updateSupervisor(ctx)
		return err
	}

	_, err := dao.Daemon.Ctx(ctx).Data(do.Daemon{
		Id:           snowflake.Generate(),
		Name:         in.Name,
		Program:      in.Program,
		Command:      in.Command,
		Directory:    in.Directory,
		RunUser:      in.RunUser,
		Numprocs:     in.Numprocs,
		Priority:     in.Priority,
		Autostart:    in.Autostart,
		Autorestart:  in.Autorestart,
		Startsecs:    in.Startsecs,
		Startretries: in.Startretries,
		StopSignal:   in.StopSignal,
		Environment:  in.Environment,
		Remark:       in.Remark,
		CreatedBy:    shared.CurrentActorUserID(ctx),
		DeptId:       shared.CurrentActorDeptID(ctx),
	}).Insert()
	if err != nil {
		_ = stopSupervisorProgram(ctx, in.Program)
		removeSupervisorProgram(in.Program)
		_ = updateSupervisor(ctx)
		removeBaoTaSupervisorConfig(in.Program)
	}
	return err
}

func (s *sDaemon) Update(ctx context.Context, in *model.DaemonUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return err
	}
	row, err := s.loadRow(ctx, in.ID)
	if err != nil {
		return err
	}
	normalizeDaemonUpdateInput(in, row.Program)
	if in.Program != "" && in.Program != row.Program {
		return gerror.New("守护进程名不支持修改，请删除后重新创建")
	}
	if err := s.validateDaemonConfig(ctx, in.ID, row.Program, in.Name, in.Command, in.Directory, in.RunUser, in.Numprocs, in.Priority, in.Autostart, in.Autorestart, in.Startsecs, in.Startretries, in.StopSignal, in.Environment); err != nil {
		return err
	}

	next := daemonConfigFromUpdateInput(in, row.Program)
	if err := writeSupervisorProgram(next); err != nil {
		return err
	}
	if err := updateSupervisor(ctx); err != nil {
		return err
	}
	if err := upsertBaoTaSupervisorConfig(next); err != nil {
		return err
	}

	_, err = dao.Daemon.Ctx(ctx).
		Where(dao.Daemon.Columns().Id, in.ID).
		Data(do.Daemon{
			Name:         in.Name,
			Command:      in.Command,
			Directory:    in.Directory,
			RunUser:      in.RunUser,
			Numprocs:     in.Numprocs,
			Priority:     in.Priority,
			Autostart:    in.Autostart,
			Autorestart:  in.Autorestart,
			Startsecs:    in.Startsecs,
			Startretries: in.Startretries,
			StopSignal:   in.StopSignal,
			Environment:  in.Environment,
			Remark:       in.Remark,
		}).
		Update()
	return err
}

func (s *sDaemon) Delete(ctx context.Context, id snowflake.JsonInt64) (*model.DaemonOperationOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = stopSupervisorProgram(ctx, row.Program)
	removeSupervisorProgram(row.Program)
	if err := updateSupervisor(ctx); err != nil {
		return nil, err
	}
	removeSupervisorLogs(row.Program)
	removeBaoTaSupervisorConfig(row.Program)
	if _, err = dao.Daemon.Ctx(ctx).
		Where(dao.Daemon.Columns().Id, id).
		Delete(); err != nil {
		return nil, err
	}
	return s.operationOutput(ctx, row.Program, "已删除"), nil
}

func (s *sDaemon) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) (*model.DaemonBatchOperationOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return nil, gerror.New("请选择要删除的守护进程")
	}
	return s.runBatch(ctx, ids, s.Delete)
}

func (s *sDaemon) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DaemonDetailOutput, err error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	out = daemonRowToOutput(row)
	s.fillRuntime(ctx, []*model.DaemonListOutput{out})
	return out, nil
}

func (s *sDaemon) List(ctx context.Context, in *model.DaemonListInput) (list []*model.DaemonListOutput, total int, err error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, 0, err
	}
	if in == nil {
		in = &model.DaemonListInput{}
	}
	normalizeDaemonListInput(in)
	m := dao.Daemon.Ctx(ctx).Where(dao.Daemon.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Daemon.Columns().Name, "%"+in.Keyword+"%").
			WhereOrLike(dao.Daemon.Columns().Program, "%"+in.Keyword+"%").
			WhereOrLike(dao.Daemon.Columns().Command, "%"+in.Keyword+"%").
			WhereOrLike(dao.Daemon.Columns().Directory, "%"+in.Keyword+"%").
			WhereOrLike(dao.Daemon.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Program != "" {
		m = m.WhereLike(dao.Daemon.Columns().Program, "%"+in.Program+"%")
	}
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(dao.Daemon.Columns().Id).Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	s.fillRuntime(ctx, list)
	return list, total, nil
}

func (s *sDaemon) Restart(ctx context.Context, id snowflake.JsonInt64) (*model.DaemonOperationOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := restartSupervisorProgram(ctx, row.Program); err != nil {
		return nil, err
	}
	return s.operationOutput(ctx, row.Program, "已重启"), nil
}

func (s *sDaemon) BatchRestart(ctx context.Context, ids []snowflake.JsonInt64) (*model.DaemonBatchOperationOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return nil, gerror.New("请选择要重启的守护进程")
	}
	return s.runBatch(ctx, ids, s.Restart)
}

func (s *sDaemon) Stop(ctx context.Context, id snowflake.JsonInt64) (*model.DaemonOperationOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := stopSupervisorProgram(ctx, row.Program); err != nil {
		return nil, err
	}
	return s.operationOutput(ctx, row.Program, "已暂停"), nil
}

func (s *sDaemon) BatchStop(ctx context.Context, ids []snowflake.JsonInt64) (*model.DaemonBatchOperationOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return nil, gerror.New("请选择要暂停的守护进程")
	}
	return s.runBatch(ctx, ids, s.Stop)
}

func (s *sDaemon) Log(ctx context.Context, id snowflake.JsonInt64, logType string, lines int) (*model.DaemonLogOutput, error) {
	if err := ensurePlatformSuperAdmin(ctx); err != nil {
		return nil, err
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	logType = normalizeLogType(logType)
	content, err := readSupervisorLog(row.Program, logType, lines)
	if err != nil {
		return nil, err
	}
	return &model.DaemonLogOutput{
		Program: row.Program,
		LogType: logType,
		Content: content,
	}, nil
}

type daemonRow struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Program      string              `json:"program"`
	Command      string              `json:"command"`
	Directory    string              `json:"directory"`
	RunUser      string              `json:"runUser"`
	Numprocs     int                 `json:"numprocs"`
	Priority     int                 `json:"priority"`
	Autostart    int                 `json:"autostart"`
	Autorestart  int                 `json:"autorestart"`
	Startsecs    int                 `json:"startsecs"`
	Startretries int                 `json:"startretries"`
	StopSignal   string              `json:"stopSignal"`
	Environment  string              `json:"environment"`
	Remark       string              `json:"remark"`
}

func normalizeDaemonCreateInput(in *model.DaemonCreateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Program = normalizeProgramName(in.Program)
	in.Command = strings.TrimSpace(in.Command)
	in.Directory = strings.TrimSpace(in.Directory)
	in.RunUser = normalizeDefault(in.RunUser, defaultRunUser)
	in.Numprocs = normalizePositiveDefault(in.Numprocs, defaultNumprocs)
	in.Priority = normalizePositiveDefault(in.Priority, defaultPriority)
	in.Autostart = normalizeBinaryDefault(in.Autostart, 1)
	in.Autorestart = normalizeBinaryDefault(in.Autorestart, 1)
	in.Startsecs = normalizePositiveDefault(in.Startsecs, defaultStartsecs)
	in.Startretries = normalizePositiveDefault(in.Startretries, defaultStartretries)
	in.StopSignal = strings.ToUpper(normalizeDefault(in.StopSignal, defaultStopSignal))
	in.Environment = strings.TrimSpace(in.Environment)
	in.Remark = strings.TrimSpace(in.Remark)
}

func normalizeDaemonUpdateInput(in *model.DaemonUpdateInput, program string) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Program = normalizeProgramName(in.Program)
	if in.Program == "" {
		in.Program = program
	}
	in.Command = strings.TrimSpace(in.Command)
	in.Directory = strings.TrimSpace(in.Directory)
	in.RunUser = normalizeDefault(in.RunUser, defaultRunUser)
	in.Numprocs = normalizePositiveDefault(in.Numprocs, defaultNumprocs)
	in.Priority = normalizePositiveDefault(in.Priority, defaultPriority)
	in.Autostart = normalizeBinaryDefault(in.Autostart, 1)
	in.Autorestart = normalizeBinaryDefault(in.Autorestart, 1)
	in.Startsecs = normalizePositiveDefault(in.Startsecs, defaultStartsecs)
	in.Startretries = normalizePositiveDefault(in.Startretries, defaultStartretries)
	in.StopSignal = strings.ToUpper(normalizeDefault(in.StopSignal, defaultStopSignal))
	in.Environment = strings.TrimSpace(in.Environment)
	in.Remark = strings.TrimSpace(in.Remark)
}

func normalizeDaemonListInput(in *model.DaemonListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Program = normalizeProgramName(in.Program)
}

func (s *sDaemon) validateDaemonConfig(
	ctx context.Context,
	currentID snowflake.JsonInt64,
	program string,
	name string,
	command string,
	directory string,
	runUser string,
	numprocs int,
	priority int,
	autostart int,
	autorestart int,
	startsecs int,
	startretries int,
	stopSignal string,
	environment string,
) error {
	if name == "" {
		return gerror.New("显示名称不能为空")
	}
	if len(name) > 80 {
		return gerror.New("显示名称长度不能超过80位")
	}
	if err := validateProgramName(program); err != nil {
		return err
	}
	if command == "" {
		return gerror.New("启动命令不能为空")
	}
	if len(command) > 1000 || strings.ContainsAny(command, "\r\n") {
		return gerror.New("启动命令格式不正确")
	}
	if containsShellMeta(command) {
		return gerror.New("启动命令包含不安全的 shell 元字符")
	}
	if err := validateDirectory(directory); err != nil {
		return err
	}
	if runUser == "" || len(runUser) > 80 || strings.ContainsAny(runUser, " \t\r\n") {
		return gerror.New("运行用户格式不正确")
	}
	if !systemUserExists(runUser) {
		return gerror.New("运行用户不存在")
	}
	if numprocs < 1 || numprocs > 64 {
		return gerror.New("进程数量必须在1到64之间")
	}
	if priority < 1 || priority > 9999 {
		return gerror.New("启动优先级必须在1到9999之间")
	}
	if err := fieldvalid.Binary("随Supervisor启动", autostart); err != nil {
		return err
	}
	if err := fieldvalid.Binary("自动重启", autorestart); err != nil {
		return err
	}
	if startsecs < 0 || startsecs > 3600 {
		return gerror.New("启动稳定秒数必须在0到3600之间")
	}
	if startretries < 0 || startretries > 100 {
		return gerror.New("启动重试次数必须在0到100之间")
	}
	if err := validateStopSignal(stopSignal); err != nil {
		return err
	}
	if len(environment) > 1000 || strings.ContainsAny(environment, "\r\n") {
		return gerror.New("环境变量格式不正确")
	}
	return s.ensureProgramUnique(ctx, currentID, program)
}

func validateProgramName(value string) error {
	if value == "" {
		return gerror.New("进程名不能为空")
	}
	if len(value) > 80 {
		return gerror.New("进程名长度不能超过80位")
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return gerror.New("进程名仅支持字母、数字、中横线和下划线")
	}
	return nil
}

func validateDirectory(value string) error {
	if value == "" {
		return gerror.New("运行目录不能为空")
	}
	if len(value) > 500 || strings.ContainsAny(value, "\r\n") {
		return gerror.New("运行目录格式不正确")
	}
	if !directoryExists(value) {
		return gerror.New("运行目录不存在")
	}
	return nil
}

func validateStopSignal(value string) error {
	switch value {
	case "TERM", "HUP", "INT", "QUIT", "KILL", "USR1", "USR2":
		return nil
	default:
		return gerror.New("停止信号不正确")
	}
}

func normalizeProgramName(value string) string {
	return strings.TrimSpace(value)
}

func normalizeDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func normalizePositiveDefault(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func normalizeBinaryDefault(value, fallback int) int {
	if value != 0 && value != 1 {
		return fallback
	}
	return value
}

func (s *sDaemon) ensureProgramUnique(ctx context.Context, currentID snowflake.JsonInt64, program string) error {
	m := dao.Daemon.Ctx(ctx).
		Where(dao.Daemon.Columns().Program, program).
		Where(dao.Daemon.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Daemon.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("进程名已存在")
	}
	return nil
}

func (s *sDaemon) loadRow(ctx context.Context, id snowflake.JsonInt64) (*daemonRow, error) {
	if id <= 0 {
		return nil, gerror.New("守护进程不存在或已删除")
	}
	row := &daemonRow{}
	if err := dao.Daemon.Ctx(ctx).
		Where(dao.Daemon.Columns().Id, id).
		Where(dao.Daemon.Columns().DeletedAt, nil).
		Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerror.New("守护进程不存在或已删除")
		}
		return nil, err
	}
	if row.ID == 0 {
		return nil, gerror.New("守护进程不存在或已删除")
	}
	return row, nil
}

func (s *sDaemon) fillRuntime(ctx context.Context, list []*model.DaemonListOutput) {
	if len(list) == 0 {
		return
	}
	programs := make([]string, 0, len(list))
	for _, item := range list {
		if item.Program != "" {
			programs = append(programs, item.Program)
		}
		item.ConfigPath = supervisorConfigPath(item.Program)
		item.OutLogPath = supervisorLogPath(item.Program, "normal")
		item.ErrLogPath = supervisorLogPath(item.Program, "error")
	}
	statuses := loadSupervisorRuntimeMap(ctx, programs)
	for _, item := range list {
		runtime := statuses[item.Program]
		if runtime.RunStatus == "" {
			runtime.RunStatus = "MISSING"
			runtime.StatusText = "未加载"
		}
		item.RunStatus = runtime.RunStatus
		item.Pid = runtime.Pid
		item.Uptime = runtime.Uptime
		item.StatusText = runtime.StatusText
	}
}

func (s *sDaemon) operationOutput(ctx context.Context, program, message string) *model.DaemonOperationOutput {
	runtime := loadSupervisorRuntimeMap(ctx, []string{program})[program]
	if runtime.RunStatus == "" {
		runtime.RunStatus = "MISSING"
	}
	return &model.DaemonOperationOutput{
		Program:   program,
		RunStatus: runtime.RunStatus,
		Message:   message,
	}
}

func (s *sDaemon) runBatch(
	ctx context.Context,
	ids []snowflake.JsonInt64,
	fn func(context.Context, snowflake.JsonInt64) (*model.DaemonOperationOutput, error),
) (*model.DaemonBatchOperationOutput, error) {
	results := make([]*model.DaemonOperationOutput, 0, len(ids))
	for _, id := range ids {
		out, err := fn(ctx, id)
		if err != nil {
			return nil, err
		}
		results = append(results, out)
	}
	return &model.DaemonBatchOperationOutput{Results: results}, nil
}

func daemonRowToOutput(row *daemonRow) *model.DaemonDetailOutput {
	if row == nil {
		return nil
	}
	return &model.DaemonDetailOutput{
		ID:           row.ID,
		Name:         row.Name,
		Program:      row.Program,
		Command:      row.Command,
		Directory:    row.Directory,
		RunUser:      row.RunUser,
		Numprocs:     row.Numprocs,
		Priority:     row.Priority,
		Autostart:    row.Autostart,
		Autorestart:  row.Autorestart,
		Startsecs:    row.Startsecs,
		Startretries: row.Startretries,
		StopSignal:   row.StopSignal,
		Environment:  row.Environment,
		Remark:       row.Remark,
	}
}

func daemonConfigFromCreateInput(in *model.DaemonCreateInput) supervisorProgram {
	return supervisorProgram{
		Name:         in.Name,
		Program:      in.Program,
		Command:      in.Command,
		Directory:    in.Directory,
		RunUser:      in.RunUser,
		Numprocs:     in.Numprocs,
		Priority:     in.Priority,
		Autostart:    in.Autostart,
		Autorestart:  in.Autorestart,
		Startsecs:    in.Startsecs,
		Startretries: in.Startretries,
		StopSignal:   in.StopSignal,
		Environment:  in.Environment,
		Remark:       in.Remark,
	}
}

func daemonConfigFromUpdateInput(in *model.DaemonUpdateInput, program string) supervisorProgram {
	return supervisorProgram{
		Name:         in.Name,
		Program:      program,
		Command:      in.Command,
		Directory:    in.Directory,
		RunUser:      in.RunUser,
		Numprocs:     in.Numprocs,
		Priority:     in.Priority,
		Autostart:    in.Autostart,
		Autorestart:  in.Autorestart,
		Startsecs:    in.Startsecs,
		Startretries: in.Startretries,
		StopSignal:   in.StopSignal,
		Environment:  in.Environment,
		Remark:       in.Remark,
	}
}

func containsShellMeta(cmd string) bool {
	for _, ch := range cmd {
		switch ch {
		case '|', ';', '&', '$', '`', '(', ')', '{', '}', '<', '>', '!', '\\':
			return true
		}
	}
	return false
}

func ensurePlatformSuperAdmin(ctx context.Context) error {
	if !shared.ResolveTenantAccessScope(ctx).All {
		return gerror.New("仅平台超级管理员可管理守护进程")
	}
	isAdmin, err := shared.HasCurrentActorAdminRole(ctx)
	if err != nil {
		return err
	}
	if !isAdmin {
		return gerror.New("仅平台超级管理员可管理守护进程")
	}
	return nil
}
