
package contract_template

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

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterContractTemplate(New())
}

func New() *sContractTemplate {
	return &sContractTemplate{}
}

type sContractTemplate struct{}

func normalizeContractTemplateIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建会员合同模板
func (s *sContractTemplate) Create(ctx context.Context, in *model.ContractTemplateCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberContractTemplate.Ctx(ctx).Data(do.MemberContractTemplate{
		Id:        id,
		TemplateName: in.TemplateName,
		TemplateType: in.TemplateType,
		Content: in.Content,
		IsDefault: in.IsDefault,
		Remark: in.Remark,
		Sort: in.Sort,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新会员合同模板
func (s *sContractTemplate) Update(ctx context.Context, in *model.ContractTemplateUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberContractTemplate{
		TemplateName: in.TemplateName,
		TemplateType: in.TemplateType,
		Content: in.Content,
		IsDefault: in.IsDefault,
		Remark: in.Remark,
		Sort: in.Sort,
		Status: in.Status,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), in.ID, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().TenantId, dao.MemberContractTemplate.Columns().MerchantId, "会员合同模板"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), in.ID, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().CreatedBy, dao.MemberContractTemplate.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberContractTemplate.Ctx(ctx).Where(dao.MemberContractTemplate.Columns().Id, in.ID).Where(dao.MemberContractTemplate.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除会员合同模板
func (s *sContractTemplate) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), id, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().TenantId, dao.MemberContractTemplate.Columns().MerchantId, "会员合同模板"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), id, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().CreatedBy, dao.MemberContractTemplate.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberContractTemplate.Ctx(ctx).Where(dao.MemberContractTemplate.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除会员合同模板
func (s *sContractTemplate) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeContractTemplateIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), normalizedIDs, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().TenantId, dao.MemberContractTemplate.Columns().MerchantId, "会员合同模板"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), normalizedIDs, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().CreatedBy, dao.MemberContractTemplate.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberContractTemplate.Ctx(ctx).WhereIn(dao.MemberContractTemplate.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取会员合同模板详情
func (s *sContractTemplate) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ContractTemplateDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), id, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().TenantId, dao.MemberContractTemplate.Columns().MerchantId, "会员合同模板"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), id, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().CreatedBy, dao.MemberContractTemplate.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.ContractTemplateDetailOutput{}
	err = dao.MemberContractTemplate.Ctx(ctx).Where(dao.MemberContractTemplate.Columns().Id, id).Where(dao.MemberContractTemplate.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("会员合同模板不存在或已删除")
	}
	// 查询租户关联显示
	if out.TenantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_tenant").Where("id", out.TenantID)
		val, err := refQuery.Value("name")
		if err == nil {
			out.TenantName = val.String()
		}
	}
	// 查询商户关联显示
	if out.MerchantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_merchant").Where("id", out.MerchantID)
		val, err := refQuery.Value("name")
		if err == nil {
			out.MerchantName = val.String()
		}
	}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *sContractTemplate) applyListFilter(ctx context.Context, in *model.ContractTemplateListInput) *gdb.Model {
	m := dao.MemberContractTemplate.Ctx(ctx).Where(dao.MemberContractTemplate.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberContractTemplate.Columns().TenantId, dao.MemberContractTemplate.Columns().MerchantId)
	if in.TemplateName != "" {
		m = m.WhereLike(dao.MemberContractTemplate.Columns().TemplateName, "%"+in.TemplateName+"%")
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberContractTemplate.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberContractTemplate.Columns().MerchantId, *in.MerchantID)
	}
	if in.IsDefault != nil {
		m = m.Where(dao.MemberContractTemplate.Columns().IsDefault, *in.IsDefault)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberContractTemplate.Columns().Status, *in.Status)
	}
	if in.TemplateType != nil {
		m = m.Where(dao.MemberContractTemplate.Columns().TemplateType, *in.TemplateType)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberContractTemplate.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberContractTemplate.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberContractTemplate.Columns().CreatedBy, dao.MemberContractTemplate.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sContractTemplate) fillRefFields(ctx context.Context, list []*model.ContractTemplateListOutput) {
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

// List 获取会员合同模板列表
func (s *sContractTemplate) List(ctx context.Context, in *model.ContractTemplateListInput) (list []*model.ContractTemplateListOutput, total int, err error) {
	if in == nil {
		in = &model.ContractTemplateListInput{}
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
func (s *sContractTemplate) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberContractTemplate.Columns().Id:        true,
		dao.MemberContractTemplate.Columns().CreatedAt: true,
		dao.MemberContractTemplate.Columns().Sort:      true,
		dao.MemberContractTemplate.Columns().Status:    true,
		dao.MemberContractTemplate.Columns().TemplateName: true,
		dao.MemberContractTemplate.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sContractTemplate) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderAsc(dao.MemberContractTemplate.Columns().Sort).OrderDesc(dao.MemberContractTemplate.Columns().Id)
}

// Export 导出会员合同模板（不分页）
func (s *sContractTemplate) Export(ctx context.Context, in *model.ContractTemplateListInput) (list []*model.ContractTemplateListOutput, err error) {
	if in == nil {
		in = &model.ContractTemplateListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑会员合同模板
func (s *sContractTemplate) BatchUpdate(ctx context.Context, in *model.ContractTemplateBatchUpdateInput) error {
	data := do.MemberContractTemplate{}
	hasChange := false
	if in.IsDefault != nil {
		data.IsDefault = *in.IsDefault
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeContractTemplateIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), normalizedIDs, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().TenantId, dao.MemberContractTemplate.Columns().MerchantId, "会员合同模板"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberContractTemplate.Ctx(ctx), normalizedIDs, dao.MemberContractTemplate.Columns().Id, dao.MemberContractTemplate.Columns().CreatedBy, dao.MemberContractTemplate.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberContractTemplate.Ctx(ctx).WhereIn(dao.MemberContractTemplate.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入会员合同模板
func (s *sContractTemplate) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.MemberContractTemplate{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.TemplateName = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.TemplateType = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Content = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.IsDefault = v
			}
		}
		idx++
		if idx < len(record) {
			data.Remark = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Sort = v
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
		if _, insertErr := dao.MemberContractTemplate.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

