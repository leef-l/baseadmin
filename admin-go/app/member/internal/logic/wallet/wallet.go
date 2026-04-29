
package wallet

import (
	"context"
	"encoding/csv"
	"io"
	"math"
	"strconv"
	"strings"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterWallet(New())
}

func New() *sWallet {
	return &sWallet{}
}

type sWallet struct{}

func normalizeWalletIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	normalized := make([]snowflake.JsonInt64, 0, len(ids))
	for _, id := range ids {
		value := int64(id)
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, id)
		if len(normalized) >= 500 {
			break
		}
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

// Create 创建会员钱包
func (s *sWallet) Create(ctx context.Context, in *model.WalletCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberWallet.Ctx(ctx).Data(do.MemberWallet{
		Id:        id,
		UserId: in.UserID,
		WalletType: in.WalletType,
		Balance: in.Balance,
		TotalIncome: in.TotalIncome,
		TotalExpense: in.TotalExpense,
		FrozenAmount: in.FrozenAmount,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新会员钱包
func (s *sWallet) Update(ctx context.Context, in *model.WalletUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberWallet{
		UserId: in.UserID,
		WalletType: in.WalletType,
		Balance: in.Balance,
		TotalIncome: in.TotalIncome,
		TotalExpense: in.TotalExpense,
		FrozenAmount: in.FrozenAmount,
		Status: in.Status,
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.MemberWallet.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
			Where(dao.MemberWallet.Columns().Id, in.ID).
			Where(dao.MemberWallet.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("会员钱包不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.MemberWallet.Table()).Ctx(ctx), in.ID, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().TenantId, dao.MemberWallet.Columns().MerchantId, "会员钱包"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.MemberWallet.Table()).Ctx(ctx), in.ID, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().CreatedBy, dao.MemberWallet.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.MemberWallet.Table()).Ctx(ctx).
			Where(dao.MemberWallet.Columns().Id, in.ID).
			Where(dao.MemberWallet.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除会员钱包
func (s *sWallet) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberWallet.Ctx(ctx), id, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().TenantId, dao.MemberWallet.Columns().MerchantId, "会员钱包"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberWallet.Ctx(ctx), id, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().CreatedBy, dao.MemberWallet.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberWallet.Ctx(ctx).Where(dao.MemberWallet.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除会员钱包
func (s *sWallet) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeWalletIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberWallet.Ctx(ctx), normalizedIDs, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().TenantId, dao.MemberWallet.Columns().MerchantId, "会员钱包"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberWallet.Ctx(ctx), normalizedIDs, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().CreatedBy, dao.MemberWallet.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberWallet.Ctx(ctx).WhereIn(dao.MemberWallet.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取会员钱包详情
func (s *sWallet) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WalletDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberWallet.Ctx(ctx), id, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().TenantId, dao.MemberWallet.Columns().MerchantId, "会员钱包"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberWallet.Ctx(ctx), id, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().CreatedBy, dao.MemberWallet.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.WalletDetailOutput{}
	err = dao.MemberWallet.Ctx(ctx).Where(dao.MemberWallet.Columns().Id, id).Where(dao.MemberWallet.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("会员钱包不存在或已删除")
	}
	// 查询会员关联显示
	if out.UserID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.UserID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("username")
		if err == nil {
			out.UserUsername = val.String()
		}
	}
	// 查询租户关联显示
	if out.TenantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_tenant").Where("id", out.TenantID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.TenantName = val.String()
		}
	}
	// 查询商户关联显示
	if out.MerchantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_merchant").Where("id", out.MerchantID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.MerchantName = val.String()
		}
	}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *sWallet) applyListFilter(ctx context.Context, in *model.WalletListInput) *gdb.Model {
	m := dao.MemberWallet.Ctx(ctx).Where(dao.MemberWallet.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberWallet.Columns().TenantId, dao.MemberWallet.Columns().MerchantId)
	if in.UserID != nil {
		m = m.Where(dao.MemberWallet.Columns().UserId, *in.UserID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberWallet.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberWallet.Columns().MerchantId, *in.MerchantID)
	}
	if in.WalletType != nil {
		m = m.Where(dao.MemberWallet.Columns().WalletType, *in.WalletType)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberWallet.Columns().Status, *in.Status)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberWallet.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberWallet.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberWallet.Columns().CreatedBy, dao.MemberWallet.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sWallet) fillRefFields(ctx context.Context, list []*model.WalletListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.UserID != 0 {
				idSet[int64(item.UserID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_user").
				Fields("id", "username")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["username"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.UserID)]; ok {
						item.UserUsername = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.TenantID != 0 {
				idSet[int64(item.TenantID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_tenant").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.TenantID)]; ok {
						item.TenantName = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.MerchantID != 0 {
				idSet[int64(item.MerchantID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_merchant").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.MerchantID)]; ok {
						item.MerchantName = val
					}
				}
			}
		}
	}
}

// List 获取会员钱包列表
func (s *sWallet) List(ctx context.Context, in *model.WalletListInput) (list []*model.WalletListOutput, total int, err error) {
	if in == nil {
		in = &model.WalletListInput{}
	}
	// PageSize 上限保护
	if in.PageSize <= 0 {
		in.PageSize = 10
	} else if in.PageSize > 500 {
		in.PageSize = 500
	}
	if in.PageNum <= 0 {
		in.PageNum = 1
	}
	m := s.applyListFilter(ctx, in)
	total, err = m.Count()
	if err != nil {
		return
	}
	// 动态排序（白名单校验防止 SQL 注入）
	m = s.applyListOrder(m, in.OrderBy, in.OrderDir)
	err = m.Page(in.PageNum, in.PageSize).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// isAllowedOrderField 校验排序字段是否在允许列表中
func (s *sWallet) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberWallet.Columns().Id:        true,
		dao.MemberWallet.Columns().CreatedAt: true,
		dao.MemberWallet.Columns().Status:    true,
		dao.MemberWallet.Columns().Balance: true,
		dao.MemberWallet.Columns().TotalIncome: true,
		dao.MemberWallet.Columns().FrozenAmount: true,
	}
	return allowed[field]
}

func (s *sWallet) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.MemberWallet.Columns().Id)
}

// Export 导出会员钱包（不分页）
func (s *sWallet) Export(ctx context.Context, in *model.WalletListInput) (list []*model.WalletListOutput, err error) {
	if in == nil {
		in = &model.WalletListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑会员钱包
func (s *sWallet) BatchUpdate(ctx context.Context, in *model.WalletBatchUpdateInput) error {
	data := do.MemberWallet{}
	hasChange := false
	if in.WalletType != nil {
		data.WalletType = *in.WalletType
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeWalletIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberWallet.Ctx(ctx), normalizedIDs, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().TenantId, dao.MemberWallet.Columns().MerchantId, "会员钱包"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberWallet.Ctx(ctx), normalizedIDs, dao.MemberWallet.Columns().Id, dao.MemberWallet.Columns().CreatedBy, dao.MemberWallet.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberWallet.Ctx(ctx).WhereIn(dao.MemberWallet.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入会员钱包
func (s *sWallet) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
	const maxImportFileSize = 10 << 20 // 10MB
	const maxImportRows = 5000
	if file.Size > maxImportFileSize {
		return 0, 0, fmt.Errorf("文件大小超过限制（最大10MB）")
	}
	f, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	// 跳过表头
	if _, err = reader.Read(); err != nil {
		return 0, 0, fmt.Errorf("读取CSV表头失败: %w", err)
	}

	rowCount := 0
	for {
		record, readErr := reader.Read()
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return success, fail, fmt.Errorf("读取CSV数据失败: %w", readErr)
		}
		if len(record) == 0 {
			continue
		}
		rowCount++
		if rowCount > maxImportRows {
			return success, fail, fmt.Errorf("导入数据超过 %d 行上限，已处理 %d 条成功、%d 条失败", maxImportRows, success, fail)
		}
		// 逐行插入
		id := snowflake.Generate()
		data := do.MemberWallet{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.UserId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.WalletType = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.Balance = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.TotalIncome = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.TotalExpense = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.FrozenAmount = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Status = v
			}
		}
		idx++
		tenantID := snowflake.JsonInt64(0)
		merchantID := snowflake.JsonInt64(0)
		middleware.ApplyTenantScopeToWrite(ctx, &tenantID, &merchantID)
		if err := middleware.EnsureTenantMerchantAccessible(ctx, tenantID, merchantID); err != nil {
			fail++
			continue
		}
		data.TenantId = tenantID
		data.MerchantId = merchantID
		if fkVal, ok := data.UserId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.MemberWallet.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

