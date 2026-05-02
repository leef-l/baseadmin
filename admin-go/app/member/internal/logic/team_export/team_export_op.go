package team_export

// 这里是"团队数据导出 + 站点裂变"的真实业务（非 codegen CRUD）。
//
// 当前实现：
//   - Export(userID)
//     1. 递归收集 userID 的全部下级（含本人），最多 5000 人
//     2. 把 member_user / member_wallet / member_warehouse_trade / member_shop_order 中
//        相关数据导出为 SQL（INSERT 语句串），打 gzip 包
//     3. 写入 admin-go/.../resource/team-export/{userID}-{ts}.sql.gz
//     4. 写 member_team_export 记录（file_url=本地相对路径，deploy_status=0）
//   - Deploy(exportID, domain) — 占位实现：
//     更新 deploy_status=1 部署中，调 bt-deploy-funddisk.sh（暂时只记录命令到日志，不真实执行），
//     成功后写 deploy_status=2、deploy_domain、deployed_at。
//     生产环境真实部署接入后再扩展为 exec.Command 调用 + 异常回滚。
//
// 安全：
//   - 导出文件不进 git，仅在服务器本地生成；后续可加签名 URL 下载。
//   - 部署脚本调用涉及 root 权限，只允许后台超级管理员账号触发（路由层做权限码校验）。

import (
	"compress/gzip"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/utility/snowflake"
)

const (
	maxExportSubtreeSize = 5000
	deployStatusInit     = 0
	deployStatusRunning  = 1
	deployStatusDone     = 2
	deployStatusFailed   = 3
)

// ExportInput 导出入参。
type ExportInput struct {
	UserID     int64
	ExportType int    // 1=手动导出 2=自动升级导出
	OperatorID int64  // 后台操作人
	Remark     string // 备注
}

// ExportOutput 导出结果。
type ExportOutput struct {
	ExportID    string
	FileURL     string
	FileSize    int64
	MemberCount int
}

// 导出任务 deploy_status 复用语义：
//   0=排队  1=运行中  2=完成（可下载）  3=失败
// 这是为了不引入新字段。前端展示时按导出场景翻译为"排队/导出中/已就绪/失败"。

// Export 立即创建一条 member_team_export 记录（deploy_status=0 排队），并启动异步任务执行真实导出。
// 业务调用方拿到 ExportID 后用 GetExport(ctx, id) 轮询状态，状态=2 时调 download 接口拿文件。
func Export(ctx context.Context, in *ExportInput) (*ExportOutput, error) {
	if in == nil || in.UserID <= 0 {
		return nil, gerror.New("目标会员 ID 不能为空")
	}

	exportID := snowflake.Generate()
	now := gtime.Now()
	if _, err := dao.MemberTeamExport.Ctx(ctx).Data(do.MemberTeamExport{
		Id:              exportID,
		UserId:          in.UserID,
		TeamMemberCount: 0,
		ExportType:      normalizeExportType(in.ExportType),
		FileUrl:         "",
		FileSize:        0,
		DeployStatus:    deployStatusInit,
		Status:          1,
		Remark:          strings.TrimSpace(in.Remark),
		CreatedBy:       in.OperatorID,
		DeptId:          0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}).Insert(); err != nil {
		return nil, err
	}

	exportRoot := getExportRoot(ctx)
	go runExportAsync(int64(exportID), in.UserID, exportRoot)

	return &ExportOutput{
		ExportID:    fmt.Sprintf("%d", int64(exportID)),
		FileURL:     "",
		FileSize:    0,
		MemberCount: 0,
	}, nil
}

// runExportAsync 后台异步生成 SQL 文件并更新记录。
func runExportAsync(exportID, userID int64, exportRoot string) {
	ctx := context.Background()
	cols := dao.MemberTeamExport.Columns()

	// 标记运行中
	_, _ = dao.MemberTeamExport.Ctx(ctx).
		Where(cols.Id, exportID).
		Data(g.Map{cols.DeployStatus: deployStatusRunning}).Update()

	subtree, err := collectSubtreeIDs(ctx, userID)
	if err != nil {
		markExportFailed(ctx, exportID, err)
		return
	}
	if len(subtree) == 0 {
		markExportFailed(ctx, exportID, gerror.New("目标会员不存在"))
		return
	}

	sqlText, err := buildSQLDump(ctx, subtree)
	if err != nil {
		markExportFailed(ctx, exportID, err)
		return
	}

	if err := os.MkdirAll(exportRoot, 0o755); err != nil {
		markExportFailed(ctx, exportID, fmt.Errorf("创建导出目录失败: %w", err))
		return
	}
	ts := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("%d-%s.sql.gz", userID, ts)
	fullPath := filepath.Join(exportRoot, fileName)
	fileSize, err := writeGzipFile(fullPath, sqlText)
	if err != nil {
		markExportFailed(ctx, exportID, err)
		return
	}

	relativeURL := "/team-export/" + fileName
	if _, err := dao.MemberTeamExport.Ctx(ctx).
		Where(cols.Id, exportID).
		Data(g.Map{
			cols.FileUrl:         relativeURL,
			cols.FileSize:        fileSize,
			cols.TeamMemberCount: len(subtree),
			cols.DeployStatus:    deployStatusDone,
		}).Update(); err != nil {
		g.Log().Errorf(ctx, "[team_export] update record err id=%d err=%v", exportID, err)
		return
	}
	g.Log().Infof(ctx, "[team_export] export done user=%d size=%d members=%d file=%s",
		userID, fileSize, len(subtree), relativeURL)
}

func markExportFailed(ctx context.Context, exportID int64, e error) {
	cols := dao.MemberTeamExport.Columns()
	_, _ = dao.MemberTeamExport.Ctx(ctx).
		Where(cols.Id, exportID).
		Data(g.Map{
			cols.DeployStatus: deployStatusFailed,
			cols.Remark:       "导出失败：" + e.Error(),
		}).Update()
	g.Log().Errorf(ctx, "[team_export] failed id=%d err=%v", exportID, e)
}

// GetExport 查询导出状态（前端轮询）。
// 失败时 errReason 取 remark 字段（markExportFailed 写入 "导出失败：原因"）。
func GetExport(ctx context.Context, exportID int64) (status int, fileURL string, fileSize int64, members int, errReason string, err error) {
	cols := dao.MemberTeamExport.Columns()
	row, err := dao.MemberTeamExport.Ctx(ctx).
		Where(cols.Id, exportID).
		Where(cols.DeletedAt, nil).
		One()
	if err != nil {
		return 0, "", 0, 0, "", err
	}
	if row.IsEmpty() {
		return 0, "", 0, 0, "", gerror.New("导出记录不存在")
	}
	st := row[cols.DeployStatus].Int()
	reason := ""
	if st == deployStatusFailed {
		reason = row[cols.Remark].String()
	}
	return st,
		row[cols.FileUrl].String(),
		row[cols.FileSize].Int64(),
		row[cols.TeamMemberCount].Int(),
		reason,
		nil
}

// DeployInput 站点裂变部署入参。
type DeployInput struct {
	ExportID   int64
	Domain     string // 目标域名（必填）
	OperatorID int64
}

// DeployOutput 部署结果。
type DeployOutput struct {
	DeployStatus int
	DeployDomain string
}

// Deploy 触发站点裂变（当前为占位）。
func Deploy(ctx context.Context, in *DeployInput) (*DeployOutput, error) {
	if in == nil || in.ExportID <= 0 {
		return nil, gerror.New("导出记录 ID 不能为空")
	}
	domain := strings.TrimSpace(in.Domain)
	if domain == "" {
		return nil, gerror.New("部署域名不能为空")
	}

	// 标记部署中
	if _, err := dao.MemberTeamExport.Ctx(ctx).
		Where(dao.MemberTeamExport.Columns().Id, in.ExportID).
		Data(g.Map{
			dao.MemberTeamExport.Columns().DeployStatus: deployStatusRunning,
			dao.MemberTeamExport.Columns().DeployDomain: domain,
		}).Update(); err != nil {
		return nil, err
	}

	// 占位：实际生产环境会执行 bt-deploy-funddisk.sh 起新站。
	// 这里只把命令记到日志，避免误触发部署脚本（root 权限）。
	g.Log().Warningf(ctx,
		"[team_export] DEPLOY PLACEHOLDER — 真实部署需手动 ssh 到服务器执行：\n"+
			"  DOMAIN=%s bash /www/wwwroot/project/fund-disk/admin-go/bt-deploy-funddisk.sh\n"+
			"  并把导出的 sql.gz 灌入新站数据库。export_id=%d",
		domain, in.ExportID,
	)

	// 暂时直接标记成功，让流程跑通；真实接入后改为部署脚本退出码判断
	now := gtime.Now()
	if _, err := dao.MemberTeamExport.Ctx(ctx).
		Where(dao.MemberTeamExport.Columns().Id, in.ExportID).
		Data(g.Map{
			dao.MemberTeamExport.Columns().DeployStatus: deployStatusDone,
			dao.MemberTeamExport.Columns().DeployedAt:   now,
		}).Update(); err != nil {
		return nil, err
	}

	return &DeployOutput{
		DeployStatus: deployStatusDone,
		DeployDomain: domain,
	}, nil
}

// ----- helpers -----

// collectSubtreeIDs 递归收集 root + 下级 ID。
func collectSubtreeIDs(ctx context.Context, root int64) ([]int64, error) {
	seen := map[int64]struct{}{root: {}}
	out := []int64{root}
	queue := []int64{root}
	for len(queue) > 0 {
		batch := queue
		queue = nil
		var rows []struct {
			Id int64 `json:"id"`
		}
		if err := dao.MemberUser.Ctx(ctx).
			Fields(dao.MemberUser.Columns().Id).
			WhereIn(dao.MemberUser.Columns().ParentId, batch).
			Where(dao.MemberUser.Columns().DeletedAt, nil).
			Scan(&rows); err != nil {
			return nil, err
		}
		for _, r := range rows {
			if r.Id <= 0 {
				continue
			}
			if _, ok := seen[r.Id]; ok {
				continue
			}
			seen[r.Id] = struct{}{}
			out = append(out, r.Id)
			queue = append(queue, r.Id)
			if len(out) > maxExportSubtreeSize {
				return nil, gerror.Newf("子树规模超过 %d 人，导出中止", maxExportSubtreeSize)
			}
		}
	}
	return out, nil
}

// buildSQLDump 把目标会员及关联数据生成 INSERT SQL 文本。
//
// 当前生成的表：
//   - member_user
//   - member_wallet
//   - member_wallet_log
//   - member_shop_order
//   - member_warehouse_trade
//
// 简单实现：每行 INSERT IGNORE，使用 16 进制转义字符串避免引号问题。
// 实际产品级导出建议改用 mysqldump，但 mysqldump 依赖外部 binary，这里先求自包含。
func buildSQLDump(ctx context.Context, ids []int64) (string, error) {
	if len(ids) == 0 {
		return "", nil
	}
	var sb strings.Builder
	sb.WriteString("-- funddisk team_export dump\n")
	sb.WriteString("-- generated_at=" + time.Now().UTC().Format(time.RFC3339) + "\n")
	sb.WriteString(fmt.Sprintf("-- members=%d\n\n", len(ids)))

	// member_user（按 id）
	if err := dumpTable(ctx, &sb, "member_user", "id", ids); err != nil {
		return "", err
	}
	// member_wallet（按 user_id）
	if err := dumpTable(ctx, &sb, "member_wallet", "user_id", ids); err != nil {
		return "", err
	}
	// member_wallet_log
	if err := dumpTable(ctx, &sb, "member_wallet_log", "user_id", ids); err != nil {
		return "", err
	}
	// member_shop_order
	if err := dumpTable(ctx, &sb, "member_shop_order", "user_id", ids); err != nil {
		return "", err
	}
	// member_warehouse_goods（按 owner_id）
	if err := dumpTable(ctx, &sb, "member_warehouse_goods", "owner_id", ids); err != nil {
		return "", err
	}
	// member_warehouse_listing（按 seller_id）
	if err := dumpTable(ctx, &sb, "member_warehouse_listing", "seller_id", ids); err != nil {
		return "", err
	}
	// member_warehouse_trade（buyer_id 或 seller_id 命中即导）
	if err := dumpTradeTable(ctx, &sb, ids); err != nil {
		return "", err
	}
	// member_level_log
	if err := dumpTable(ctx, &sb, "member_level_log", "user_id", ids); err != nil {
		return "", err
	}
	// member_rebind_log
	if err := dumpTable(ctx, &sb, "member_rebind_log", "user_id", ids); err != nil {
		return "", err
	}
	// member_contract（按 user_id）
	if err := dumpTable(ctx, &sb, "member_contract", "user_id", ids); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// dumpTable 把 table 中 column IN (ids) 的所有行 dump 成 INSERT IGNORE。
func dumpTable(ctx context.Context, sb *strings.Builder, table, column string, ids []int64) error {
	rows, err := g.DB().Ctx(ctx).Model(table).WhereIn(column, ids).All()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	fields := recordKeys(rows[0])
	colList := "`" + strings.Join(fields, "`,`") + "`"
	sb.WriteString(fmt.Sprintf("-- %s rows=%d\n", table, len(rows)))
	for _, row := range rows {
		values := make([]string, 0, len(fields))
		for _, f := range fields {
			values = append(values, sqlValue(row[f]))
		}
		sb.WriteString(fmt.Sprintf("INSERT IGNORE INTO `%s` (%s) VALUES (%s);\n", table, colList, strings.Join(values, ",")))
	}
	sb.WriteString("\n")
	return nil
}

// dumpTradeTable 单独处理 trade（buyer_id 或 seller_id 任一命中）。
func dumpTradeTable(ctx context.Context, sb *strings.Builder, ids []int64) error {
	rows, err := g.DB().Ctx(ctx).Model("member_warehouse_trade").
		Where("buyer_id IN (?) OR seller_id IN (?)", ids, ids).
		All()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	fields := recordKeys(rows[0])
	colList := "`" + strings.Join(fields, "`,`") + "`"
	sb.WriteString(fmt.Sprintf("-- member_warehouse_trade rows=%d\n", len(rows)))
	for _, row := range rows {
		values := make([]string, 0, len(fields))
		for _, f := range fields {
			values = append(values, sqlValue(row[f]))
		}
		sb.WriteString(fmt.Sprintf("INSERT IGNORE INTO `member_warehouse_trade` (%s) VALUES (%s);\n", colList, strings.Join(values, ",")))
	}
	sb.WriteString("\n")
	return nil
}

// recordKeys 提取 gdb.Record 的 key 列表（map 遍历无序，但每次 dump 内部一致即可）。
func recordKeys(row map[string]*gvar.Var) []string {
	keys := make([]string, 0, len(row))
	for k := range row {
		keys = append(keys, k)
	}
	return keys
}

// sqlValue 把 gvar.Var 转字符串字面量。NULL → NULL；非空 → 用 X'...' 16 进制（MySQL 5.7+），避免引号转义。
func sqlValue(v *gvar.Var) string {
	if v == nil || v.IsNil() {
		return "NULL"
	}
	s := v.String()
	if s == "" {
		// 区分 NULL vs ''：gvar.Var 在 NULL 时 IsNil()==true，已被上面分支拦截，这里 s == "" 视为空字符串。
		return "''"
	}
	return "X'" + hex.EncodeToString([]byte(s)) + "'"
}

// writeGzipFile 把 sql 文本 gzip 压缩写文件，返回压缩后字节数。
func writeGzipFile(path, content string) (int64, error) {
	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("创建导出文件失败: %w", err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	if _, err := gw.Write([]byte(content)); err != nil {
		gw.Close()
		return 0, err
	}
	if err := gw.Close(); err != nil {
		return 0, err
	}
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func normalizeExportType(t int) int {
	if t == 1 || t == 2 {
		return t
	}
	return 1
}

// getExportRoot 决定导出根目录（始终返回绝对路径，避免 cwd 漂移）。
//   - 配置 member.teamExportRoot 优先
//   - 否则用 ./resource/team-export 然后转绝对
func getExportRoot(ctx context.Context) string {
	v := strings.TrimSpace(g.Cfg().MustGet(ctx, "member.teamExportRoot").String())
	if v == "" {
		v = filepath.Join("resource", "team-export")
	}
	if abs, err := filepath.Abs(v); err == nil {
		return abs
	}
	return v
}

