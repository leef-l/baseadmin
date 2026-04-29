
package contract

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
	"golang.org/x/crypto/bcrypt"

	"gbaseadmin/app/demo/internal/dao"
	"gbaseadmin/app/demo/internal/middleware"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/model/do"
	"gbaseadmin/app/demo/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterContract(New())
}

func New() *sContract {
	return &sContract{}
}

type sContract struct{}

func normalizeContractIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建体验合同
func (s *sContract) Create(ctx context.Context, in *model.ContractCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	hashedSignPassword, err := bcrypt.GenerateFromPassword([]byte(in.SignPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.DemoContract.Ctx(ctx).Data(do.DemoContract{
		Id:        id,
		ContractNo: in.ContractNo,
		CustomerId: in.CustomerID,
		OrderId: in.OrderID,
		Title: in.Title,
		ContractFile: in.ContractFile,
		SignImage: in.SignImage,
		ContractAmount: in.ContractAmount,
		SignPassword: string(hashedSignPassword),
		SignedAt: in.SignedAt,
		ExpiresAt: in.ExpiresAt,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验合同
func (s *sContract) Update(ctx context.Context, in *model.ContractUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.DemoContract{
		ContractNo: in.ContractNo,
		CustomerId: in.CustomerID,
		OrderId: in.OrderID,
		Title: in.Title,
		ContractFile: in.ContractFile,
		SignImage: in.SignImage,
		ContractAmount: in.ContractAmount,
		SignedAt: in.SignedAt,
		ExpiresAt: in.ExpiresAt,
		Status: in.Status,
	}
	if in.SignPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(in.SignPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.SignPassword = string(hashed)
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.DemoContract.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.DemoContract.Table()).Ctx(ctx).
			Where(dao.DemoContract.Columns().Id, in.ID).
			Where(dao.DemoContract.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("体验合同不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.DemoContract.Table()).Ctx(ctx), in.ID, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().TenantId, dao.DemoContract.Columns().MerchantId, "体验合同"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.DemoContract.Table()).Ctx(ctx), in.ID, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().CreatedBy, dao.DemoContract.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.DemoContract.Table()).Ctx(ctx).
			Where(dao.DemoContract.Columns().Id, in.ID).
			Where(dao.DemoContract.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除体验合同
func (s *sContract) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoContract.Ctx(ctx), id, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().TenantId, dao.DemoContract.Columns().MerchantId, "体验合同"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoContract.Ctx(ctx), id, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().CreatedBy, dao.DemoContract.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoContract.Ctx(ctx).Where(dao.DemoContract.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验合同
func (s *sContract) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeContractIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoContract.Ctx(ctx), normalizedIDs, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().TenantId, dao.DemoContract.Columns().MerchantId, "体验合同"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoContract.Ctx(ctx), normalizedIDs, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().CreatedBy, dao.DemoContract.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoContract.Ctx(ctx).WhereIn(dao.DemoContract.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验合同详情
func (s *sContract) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ContractDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoContract.Ctx(ctx), id, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().TenantId, dao.DemoContract.Columns().MerchantId, "体验合同"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoContract.Ctx(ctx), id, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().CreatedBy, dao.DemoContract.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.ContractDetailOutput{}
	err = dao.DemoContract.Ctx(ctx).Where(dao.DemoContract.Columns().Id, id).Where(dao.DemoContract.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("体验合同不存在或已删除")
	}
	// 查询客户关联显示
	if out.CustomerID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_customer").Where("id", out.CustomerID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.CustomerName = val.String()
		}
	}
	// 查询订单关联显示
	if out.OrderID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_order").Where("id", out.OrderID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("order_no")
		if err == nil {
			out.OrderOrderNo = val.String()
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
func (s *sContract) applyListFilter(ctx context.Context, in *model.ContractListInput) *gdb.Model {
	m := dao.DemoContract.Ctx(ctx).Where(dao.DemoContract.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoContract.Columns().TenantId, dao.DemoContract.Columns().MerchantId)
	if in.ContractNo != "" {
		m = m.Where(dao.DemoContract.Columns().ContractNo, in.ContractNo)
	}
	if in.Title != "" {
		m = m.WhereLike(dao.DemoContract.Columns().Title, "%"+in.Title+"%")
	}
	if in.CustomerID != nil {
		m = m.Where(dao.DemoContract.Columns().CustomerId, *in.CustomerID)
	}
	if in.OrderID != nil {
		m = m.Where(dao.DemoContract.Columns().OrderId, *in.OrderID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoContract.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoContract.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoContract.Columns().Status, *in.Status)
	}
	if in.SignedAtStart != "" {
		m = m.WhereGTE(dao.DemoContract.Columns().SignedAt, in.SignedAtStart)
	}
	if in.SignedAtEnd != "" {
		m = m.WhereLTE(dao.DemoContract.Columns().SignedAt, in.SignedAtEnd)
	}
	if in.ExpiresAtStart != "" {
		m = m.WhereGTE(dao.DemoContract.Columns().ExpiresAt, in.ExpiresAtStart)
	}
	if in.ExpiresAtEnd != "" {
		m = m.WhereLTE(dao.DemoContract.Columns().ExpiresAt, in.ExpiresAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoContract.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoContract.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoContract.Columns().CreatedBy, dao.DemoContract.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sContract) fillRefFields(ctx context.Context, list []*model.ContractListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.CustomerID != 0 {
				idSet[int64(item.CustomerID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("demo_customer").
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
					if val, ok := refMap[int64(item.CustomerID)]; ok {
						item.CustomerName = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.OrderID != 0 {
				idSet[int64(item.OrderID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("demo_order").
				Fields("id", "order_no")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["order_no"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.OrderID)]; ok {
						item.OrderOrderNo = val
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

// List 获取体验合同列表
func (s *sContract) List(ctx context.Context, in *model.ContractListInput) (list []*model.ContractListOutput, total int, err error) {
	if in == nil {
		in = &model.ContractListInput{}
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
func (s *sContract) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoContract.Columns().Id:        true,
		dao.DemoContract.Columns().CreatedAt: true,
		dao.DemoContract.Columns().Status:    true,
		dao.DemoContract.Columns().ContractNo: true,
		dao.DemoContract.Columns().Title: true,
		dao.DemoContract.Columns().ContractAmount: true,
	}
	return allowed[field]
}

func (s *sContract) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoContract.Columns().Id)
}

// Export 导出体验合同（不分页）
func (s *sContract) Export(ctx context.Context, in *model.ContractListInput) (list []*model.ContractListOutput, err error) {
	if in == nil {
		in = &model.ContractListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验合同
func (s *sContract) BatchUpdate(ctx context.Context, in *model.ContractBatchUpdateInput) error {
	data := do.DemoContract{}
	hasChange := false
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeContractIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoContract.Ctx(ctx), normalizedIDs, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().TenantId, dao.DemoContract.Columns().MerchantId, "体验合同"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoContract.Ctx(ctx), normalizedIDs, dao.DemoContract.Columns().Id, dao.DemoContract.Columns().CreatedBy, dao.DemoContract.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoContract.Ctx(ctx).WhereIn(dao.DemoContract.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验合同
func (s *sContract) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.DemoContract{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.ContractNo = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.CustomerId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.OrderId = v
			}
		}
		idx++
		if idx < len(record) {
			data.Title = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.ContractFile = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.SignImage = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.ContractAmount = int64(math.Round(v * 100))
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
		if fkVal, ok := data.CustomerId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("demo_customer").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.OrderId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("demo_order").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.DemoContract.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

