
package appointment

import (
	"context"
	"encoding/csv"
	"io"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
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
	service.RegisterAppointment(New())
}

func New() *sAppointment {
	return &sAppointment{}
}

type sAppointment struct{}

func normalizeAppointmentIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

// Create 创建体验预约
func (s *sAppointment) Create(ctx context.Context, in *model.AppointmentCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoAppointment.Ctx(ctx).Data(do.DemoAppointment{
		Id:        id,
		AppointmentNo: in.AppointmentNo,
		CustomerId: in.CustomerID,
		Subject: in.Subject,
		AppointmentAt: in.AppointmentAt,
		ContactPhone: in.ContactPhone,
		Address: in.Address,
		Remark: in.Remark,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验预约
func (s *sAppointment) Update(ctx context.Context, in *model.AppointmentUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoAppointment.Ctx(ctx), in.ID, dao.DemoAppointment.Columns().Id, dao.DemoAppointment.Columns().TenantId, dao.DemoAppointment.Columns().MerchantId, "体验预约"); err != nil {
		return err
	}
	data := do.DemoAppointment{
		AppointmentNo: in.AppointmentNo,
		CustomerId: in.CustomerID,
		Subject: in.Subject,
		AppointmentAt: in.AppointmentAt,
		ContactPhone: in.ContactPhone,
		Address: in.Address,
		Remark: in.Remark,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
	}
	_, err := dao.DemoAppointment.Ctx(ctx).Where(dao.DemoAppointment.Columns().Id, in.ID).Data(data).Update()
	return err
}

// Delete 软删除体验预约
func (s *sAppointment) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoAppointment.Ctx(ctx), id, dao.DemoAppointment.Columns().Id, dao.DemoAppointment.Columns().TenantId, dao.DemoAppointment.Columns().MerchantId, "体验预约"); err != nil {
		return err
	}
	_, err := dao.DemoAppointment.Ctx(ctx).Where(dao.DemoAppointment.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验预约
func (s *sAppointment) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeAppointmentIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoAppointment.Ctx(ctx), normalizedIDs, dao.DemoAppointment.Columns().Id, dao.DemoAppointment.Columns().TenantId, dao.DemoAppointment.Columns().MerchantId, "体验预约"); err != nil {
		return err
	}
	_, err := dao.DemoAppointment.Ctx(ctx).WhereIn(dao.DemoAppointment.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验预约详情
func (s *sAppointment) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.AppointmentDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoAppointment.Ctx(ctx), id, dao.DemoAppointment.Columns().Id, dao.DemoAppointment.Columns().TenantId, dao.DemoAppointment.Columns().MerchantId, "体验预约"); err != nil {
		return nil, err
	}
	out = &model.AppointmentDetailOutput{}
	err = dao.DemoAppointment.Ctx(ctx).Where(dao.DemoAppointment.Columns().Id, id).Where(dao.DemoAppointment.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out.ID == 0 {
		return nil, nil
	}
	// 查询客户关联显示
	if out.CustomerID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_customer").Where("id", out.CustomerID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("name")
		if err == nil {
			out.CustomerName = val.String()
		}
	}
	// 查询租户关联显示
	if out.TenantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_tenant").Where("id", out.TenantID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("name")
		if err == nil {
			out.TenantName = val.String()
		}
	}
	// 查询商户关联显示
	if out.MerchantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_merchant").Where("id", out.MerchantID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("name")
		if err == nil {
			out.MerchantName = val.String()
		}
	}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *sAppointment) applyListFilter(ctx context.Context, in *model.AppointmentListInput) *gdb.Model {
	m := dao.DemoAppointment.Ctx(ctx).Where(dao.DemoAppointment.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoAppointment.Columns().TenantId, dao.DemoAppointment.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.DemoAppointment.Columns().Subject, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoAppointment.Columns().ContactPhone, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoAppointment.Columns().Address, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoAppointment.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.AppointmentNo != "" {
		m = m.Where(dao.DemoAppointment.Columns().AppointmentNo, in.AppointmentNo)
	}
	if in.Subject != "" {
		m = m.WhereLike(dao.DemoAppointment.Columns().Subject, "%"+in.Subject+"%")
	}
	if in.ContactPhone != "" {
		m = m.WhereLike(dao.DemoAppointment.Columns().ContactPhone, "%"+in.ContactPhone+"%")
	}
	if in.CustomerID != nil {
		m = m.Where(dao.DemoAppointment.Columns().CustomerId, *in.CustomerID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoAppointment.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoAppointment.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoAppointment.Columns().Status, *in.Status)
	}
	if in.AppointmentAtStart != "" {
		m = m.WhereGTE(dao.DemoAppointment.Columns().AppointmentAt, in.AppointmentAtStart)
	}
	if in.AppointmentAtEnd != "" {
		m = m.WhereLTE(dao.DemoAppointment.Columns().AppointmentAt, in.AppointmentAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoAppointment.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoAppointment.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoAppointment.Columns().CreatedBy, dao.DemoAppointment.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sAppointment) fillRefFields(ctx context.Context, list []*model.AppointmentListOutput) {
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

// List 获取体验预约列表
func (s *sAppointment) List(ctx context.Context, in *model.AppointmentListInput) (list []*model.AppointmentListOutput, total int, err error) {
	if in == nil {
		in = &model.AppointmentListInput{}
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
func (s *sAppointment) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoAppointment.Columns().Id:        true,
		dao.DemoAppointment.Columns().CreatedAt: true,
		dao.DemoAppointment.Columns().Status:    true,
		dao.DemoAppointment.Columns().AppointmentNo: true,
		dao.DemoAppointment.Columns().ContactPhone: true,
		dao.DemoAppointment.Columns().Address: true,
		dao.DemoAppointment.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sAppointment) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoAppointment.Columns().Id)
}

// Export 导出体验预约（不分页）
func (s *sAppointment) Export(ctx context.Context, in *model.AppointmentListInput) (list []*model.AppointmentListOutput, err error) {
	if in == nil {
		in = &model.AppointmentListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验预约
func (s *sAppointment) BatchUpdate(ctx context.Context, in *model.AppointmentBatchUpdateInput) error {
	data := do.DemoAppointment{}
	if in.Status != nil {
		data.Status = *in.Status
	}
	normalizedIDs := normalizeAppointmentIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoAppointment.Ctx(ctx), normalizedIDs, dao.DemoAppointment.Columns().Id, dao.DemoAppointment.Columns().TenantId, dao.DemoAppointment.Columns().MerchantId, "体验预约"); err != nil {
		return err
	}
	_, err := dao.DemoAppointment.Ctx(ctx).WhereIn(dao.DemoAppointment.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验预约
func (s *sAppointment) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
	f, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// 跳过表头
	if _, err = reader.Read(); err != nil {
		return 0, 0, fmt.Errorf("读取CSV表头失败: %w", err)
	}

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
		// 逐行插入
		id := snowflake.Generate()
		data := do.DemoAppointment{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.AppointmentNo = record[idx]
		}
		idx++
		if idx < len(record) {
			data.CustomerId = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Subject = record[idx]
		}
		idx++
		if idx < len(record) {
			data.ContactPhone = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Address = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Remark = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Status = record[idx]
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
		if _, insertErr := dao.DemoAppointment.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

