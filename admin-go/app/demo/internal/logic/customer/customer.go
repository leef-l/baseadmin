
package customer

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/demo/internal/dao"
	"gbaseadmin/app/demo/internal/middleware"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/model/do"
	"gbaseadmin/app/demo/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterCustomer(New())
}

func New() *sCustomer {
	return &sCustomer{}
}

type sCustomer struct{}

func normalizeCustomerIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建体验客户
func (s *sCustomer) Create(ctx context.Context, in *model.CustomerCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoCustomer.Ctx(ctx).Data(do.DemoCustomer{
		Id:        id,
		Avatar: in.Avatar,
		Name: in.Name,
		CustomerNo: in.CustomerNo,
		Phone: in.Phone,
		Email: in.Email,
		Gender: in.Gender,
		Level: in.Level,
		SourceType: in.SourceType,
		IsVip: in.IsVip,
		RegisteredAt: in.RegisteredAt,
		Remark: in.Remark,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验客户
func (s *sCustomer) Update(ctx context.Context, in *model.CustomerUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.DemoCustomer{
		Avatar: in.Avatar,
		Name: in.Name,
		CustomerNo: in.CustomerNo,
		Phone: in.Phone,
		Email: in.Email,
		Gender: in.Gender,
		Level: in.Level,
		SourceType: in.SourceType,
		IsVip: in.IsVip,
		RegisteredAt: in.RegisteredAt,
		Remark: in.Remark,
		Status: in.Status,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoCustomer.Ctx(ctx), in.ID, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().TenantId, dao.DemoCustomer.Columns().MerchantId, "体验客户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoCustomer.Ctx(ctx), in.ID, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().CreatedBy, dao.DemoCustomer.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoCustomer.Ctx(ctx).Where(dao.DemoCustomer.Columns().Id, in.ID).Where(dao.DemoCustomer.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除体验客户
func (s *sCustomer) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoCustomer.Ctx(ctx), id, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().TenantId, dao.DemoCustomer.Columns().MerchantId, "体验客户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoCustomer.Ctx(ctx), id, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().CreatedBy, dao.DemoCustomer.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoCustomer.Ctx(ctx).Where(dao.DemoCustomer.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验客户
func (s *sCustomer) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeCustomerIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoCustomer.Ctx(ctx), normalizedIDs, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().TenantId, dao.DemoCustomer.Columns().MerchantId, "体验客户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoCustomer.Ctx(ctx), normalizedIDs, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().CreatedBy, dao.DemoCustomer.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoCustomer.Ctx(ctx).WhereIn(dao.DemoCustomer.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验客户详情
func (s *sCustomer) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.CustomerDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoCustomer.Ctx(ctx), id, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().TenantId, dao.DemoCustomer.Columns().MerchantId, "体验客户"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoCustomer.Ctx(ctx), id, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().CreatedBy, dao.DemoCustomer.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.CustomerDetailOutput{}
	err = dao.DemoCustomer.Ctx(ctx).Where(dao.DemoCustomer.Columns().Id, id).Where(dao.DemoCustomer.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("体验客户不存在或已删除")
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
func (s *sCustomer) applyListFilter(ctx context.Context, in *model.CustomerListInput) *gdb.Model {
	m := dao.DemoCustomer.Ctx(ctx).Where(dao.DemoCustomer.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoCustomer.Columns().TenantId, dao.DemoCustomer.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.DemoCustomer.Columns().Name, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoCustomer.Columns().Phone, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoCustomer.Columns().Email, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoCustomer.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.CustomerNo != "" {
		m = m.Where(dao.DemoCustomer.Columns().CustomerNo, in.CustomerNo)
	}
	if in.Name != "" {
		m = m.WhereLike(dao.DemoCustomer.Columns().Name, "%"+in.Name+"%")
	}
	if in.Phone != "" {
		m = m.WhereLike(dao.DemoCustomer.Columns().Phone, "%"+in.Phone+"%")
	}
	if in.Email != "" {
		m = m.WhereLike(dao.DemoCustomer.Columns().Email, "%"+in.Email+"%")
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoCustomer.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoCustomer.Columns().MerchantId, *in.MerchantID)
	}
	if in.Gender != nil {
		m = m.Where(dao.DemoCustomer.Columns().Gender, *in.Gender)
	}
	if in.Level != nil {
		m = m.Where(dao.DemoCustomer.Columns().Level, *in.Level)
	}
	if in.SourceType != nil {
		m = m.Where(dao.DemoCustomer.Columns().SourceType, *in.SourceType)
	}
	if in.IsVip != nil {
		m = m.Where(dao.DemoCustomer.Columns().IsVip, *in.IsVip)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoCustomer.Columns().Status, *in.Status)
	}
	if in.RegisteredAtStart != "" {
		m = m.WhereGTE(dao.DemoCustomer.Columns().RegisteredAt, in.RegisteredAtStart)
	}
	if in.RegisteredAtEnd != "" {
		m = m.WhereLTE(dao.DemoCustomer.Columns().RegisteredAt, in.RegisteredAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoCustomer.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoCustomer.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoCustomer.Columns().CreatedBy, dao.DemoCustomer.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sCustomer) fillRefFields(ctx context.Context, list []*model.CustomerListOutput) {
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

// List 获取体验客户列表
func (s *sCustomer) List(ctx context.Context, in *model.CustomerListInput) (list []*model.CustomerListOutput, total int, err error) {
	if in == nil {
		in = &model.CustomerListInput{}
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
func (s *sCustomer) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoCustomer.Columns().Id:        true,
		dao.DemoCustomer.Columns().CreatedAt: true,
		dao.DemoCustomer.Columns().Status:    true,
		dao.DemoCustomer.Columns().Name: true,
		dao.DemoCustomer.Columns().CustomerNo: true,
		dao.DemoCustomer.Columns().Phone: true,
		dao.DemoCustomer.Columns().Email: true,
		dao.DemoCustomer.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sCustomer) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoCustomer.Columns().Id)
}

// Export 导出体验客户（不分页）
func (s *sCustomer) Export(ctx context.Context, in *model.CustomerListInput) (list []*model.CustomerListOutput, err error) {
	if in == nil {
		in = &model.CustomerListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验客户
func (s *sCustomer) BatchUpdate(ctx context.Context, in *model.CustomerBatchUpdateInput) error {
	data := do.DemoCustomer{}
	hasChange := false
	if in.Gender != nil {
		data.Gender = *in.Gender
		hasChange = true
	}
	if in.Level != nil {
		data.Level = *in.Level
		hasChange = true
	}
	if in.SourceType != nil {
		data.SourceType = *in.SourceType
		hasChange = true
	}
	if in.IsVip != nil {
		data.IsVip = *in.IsVip
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeCustomerIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoCustomer.Ctx(ctx), normalizedIDs, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().TenantId, dao.DemoCustomer.Columns().MerchantId, "体验客户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoCustomer.Ctx(ctx), normalizedIDs, dao.DemoCustomer.Columns().Id, dao.DemoCustomer.Columns().CreatedBy, dao.DemoCustomer.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoCustomer.Ctx(ctx).WhereIn(dao.DemoCustomer.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验客户
func (s *sCustomer) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.DemoCustomer{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.Avatar = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Name = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.CustomerNo = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Phone = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Email = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Gender = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Level = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.SourceType = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.IsVip = v
			}
		}
		idx++
		if idx < len(record) {
			data.Remark = strings.TrimSpace(record[idx])
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
		if _, insertErr := dao.DemoCustomer.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

