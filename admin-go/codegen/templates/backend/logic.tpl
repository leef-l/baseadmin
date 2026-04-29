{{- $hasRef := false}}
{{- $hasNumericImport := false}}
{{- range .Fields}}
{{- if and .RefFieldName (not .IsHidden)}}
{{- $hasRef = true}}
{{- end}}
{{- if and (not .IsHidden) (not .IsID) (not .IsPassword) (not .IsTimeField) (or (eq .GoType "int") (eq .GoType "int64") (eq .GoType "float64") (eq .GoType "JsonInt64") .IsMoney) (ne .Name "tenant_id") (ne .Name "merchant_id")}}
{{- $hasNumericImport = true}}
{{- end}}
{{- end}}
package {{.PackageName}}

import (
	"context"
{{- if .HasImport}}
	"encoding/csv"
	"io"
{{- if .HasMoney}}
	"math"
{{- end}}
{{- if $hasNumericImport}}
	"strconv"
{{- end}}
	"strings"
{{- end}}
{{- if or .EnableOpLog .HasImport}}
	"fmt"
{{- end}}

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
{{- if or $hasRef (and .HasImport .HasForeignKey)}}
	"github.com/gogf/gf/v2/frame/g"
{{- end}}
{{- if .HasImport}}
	"github.com/gogf/gf/v2/net/ghttp"
{{- end}}
{{- if .HasPassword}}
	"golang.org/x/crypto/bcrypt"
{{- end}}

	"gbaseadmin/app/{{.AppName}}/internal/dao"
{{- if or .HasCreatedBy .HasDeptID .HasTenantScope}}
	"gbaseadmin/app/{{.AppName}}/internal/middleware"
{{- end}}
	"gbaseadmin/app/{{.AppName}}/internal/model"
	"gbaseadmin/app/{{.AppName}}/internal/model/do"
	"gbaseadmin/app/{{.AppName}}/internal/service"
	"gbaseadmin/utility/snowflake"
{{- if .EnableOpLog}}
	"gbaseadmin/utility/oplog"
{{- end}}
)

func init() {
	service.Register{{.ModelName}}(New())
}

func New() *s{{.ModelName}} {
	return &s{{.ModelName}}{}
}

type s{{.ModelName}} struct{}

func normalize{{.ModelName}}IDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建{{.Comment}}
func (s *s{{.ModelName}}) Create(ctx context.Context, in *model.{{.ModelName}}CreateInput) error {
	id := snowflake.Generate()
{{- if .HasTenantScope}}
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID{{if .HasMerchantID}}, &in.MerchantID{{else}}, nil{{end}})
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID{{if .HasMerchantID}}, in.MerchantID{{else}}, 0{{end}}); err != nil {
		return err
	}
{{- end}}
{{- range .Fields}}
{{- if .IsPassword}}
	hashed{{.NameCamel}}, err := bcrypt.GenerateFromPassword([]byte(in.{{.NameCamel}}), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
{{- end}}
{{- end}}
	_, err{{if not .HasPassword}} :={{else}} ={{end}} dao.{{.DaoName}}.Ctx(ctx).Data(do.{{.DaoName}}{
		Id:        id,
{{- range .Fields}}
{{- if and (not .IsID) (not .IsHidden)}}
{{- if .IsPassword}}
		{{.NameDao}}: string(hashed{{.NameCamel}}),
{{- else}}
		{{.NameDao}}: in.{{.NameCamel}},
{{- end}}
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
		CreatedBy: middleware.GetUserID(ctx),
{{- end}}
{{- if .HasDeptID}}
		DeptId: middleware.GetDeptID(ctx),
{{- end}}
	}).Insert()
{{- if .EnableOpLog}}
	if err == nil {
		oplog.Record(ctx, "{{.ModuleName}}", "create", fmt.Sprintf("%v", id), "")
	}
{{- end}}
	return err
}

// Update 更新{{.Comment}}
func (s *s{{.ModelName}}) Update(ctx context.Context, in *model.{{.ModelName}}UpdateInput) error {
{{- if .HasTenantScope}}
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID{{if .HasMerchantID}}, &in.MerchantID{{else}}, nil{{end}})
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID{{if .HasMerchantID}}, in.MerchantID{{else}}, 0{{end}}); err != nil {
		return err
	}
{{- end}}
{{- if .HasParentID}}
	if in.ParentID == in.ID {
		return gerror.New("不能将自身设为父节点")
	}
	if int64(in.ParentID) != 0 {
		childIDs, collectErr := s.collectChildIDs(ctx, in.ID)
		if collectErr != nil {
			return collectErr
		}
		for _, cid := range childIDs {
			if cid == in.ParentID {
				return gerror.New("不能将子节点设为父节点，会形成循环引用")
			}
		}
	}
{{- end}}
	data := do.{{.DaoName}}{
{{- range .Fields}}
{{- if and (not .IsID) (not .IsHidden) (not .IsPassword) (ne .Name "tenant_id") (ne .Name "merchant_id")}}
		{{.NameDao}}: in.{{.NameCamel}},
{{- end}}
{{- end}}
	}
{{- range .Fields}}
{{- if .IsPassword}}
	if in.{{.NameCamel}} != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(in.{{.NameCamel}}), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.{{.NameDao}} = string(hashed)
	}
{{- end}}
{{- end}}
{{- if .HasMoney}}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.{{.DaoName}}.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.{{.DaoName}}.Table()).Ctx(ctx).
			Where(dao.{{.DaoName}}.Columns().Id, in.ID).
			Where(dao.{{.DaoName}}.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("{{.Comment}}不存在或已删除")
		}
{{- if .HasTenantScope}}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.{{.DaoName}}.Table()).Ctx(ctx), in.ID, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
			return err
		}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.{{.DaoName}}.Table()).Ctx(ctx), in.ID, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
			return err
		}
{{- end}}
		_, err = tx.Model(dao.{{.DaoName}}.Table()).Ctx(ctx).
			Where(dao.{{.DaoName}}.Columns().Id, in.ID).
			Where(dao.{{.DaoName}}.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
{{- else}}
{{- if .HasTenantScope}}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), in.ID, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), in.ID, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return err
	}
{{- end}}
	_, err := dao.{{.DaoName}}.Ctx(ctx).Where(dao.{{.DaoName}}.Columns().Id, in.ID).Where(dao.{{.DaoName}}.Columns().DeletedAt, nil).Data(data).Update()
{{- end}}
{{- if .EnableOpLog}}
	if err == nil {
		oplog.Record(ctx, "{{.ModuleName}}", "update", fmt.Sprintf("%v", in.ID), "")
	}
{{- end}}
	return err
}

// Delete 软删除{{.Comment}}
func (s *s{{.ModelName}}) Delete(ctx context.Context, id snowflake.JsonInt64) error {
{{- if .HasParentID}}
	deleteIDs, err := s.collectDeleteIDs(ctx, []snowflake.JsonInt64{id})
	if err != nil {
		return err
	}
	if len(deleteIDs) == 0 {
		return nil
	}
{{- if .HasTenantScope}}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), deleteIDs, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), deleteIDs, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return err
	}
{{- end}}
	_, err = dao.{{.DaoName}}.Ctx(ctx).WhereIn(dao.{{.DaoName}}.Columns().Id, deleteIDs).Delete()
{{- else}}
{{- if .HasTenantScope}}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), id, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), id, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return err
	}
{{- end}}
	_, err := dao.{{.DaoName}}.Ctx(ctx).Where(dao.{{.DaoName}}.Columns().Id, id).Delete()
{{- end}}
{{- if .EnableOpLog}}
	if err == nil {
		oplog.Record(ctx, "{{.ModuleName}}", "delete", fmt.Sprintf("%v", id), "")
	}
{{- end}}
	return err
}
{{- if .HasParentID}}

// BatchDelete 批量软删除{{.Comment}}
func (s *s{{.ModelName}}) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	deleteIDs, err := s.collectDeleteIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(deleteIDs) == 0 {
		return nil
	}
{{- if .HasTenantScope}}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), deleteIDs, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), deleteIDs, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return err
	}
{{- end}}
	_, err = dao.{{.DaoName}}.Ctx(ctx).WhereIn(dao.{{.DaoName}}.Columns().Id, deleteIDs).Delete()
{{- if .EnableOpLog}}
	if err == nil {
		oplog.Record(ctx, "{{.ModuleName}}", "batch-delete", fmt.Sprintf("%v", deleteIDs), "")
	}
{{- end}}
	return err
}

// collectDeleteIDs 汇总批量删除所需的节点 ID，并补齐所有子节点
func (s *s{{.ModelName}}) collectDeleteIDs(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := normalize{{.ModelName}}IDs(ids)
	if len(normalized) == 0 {
		return nil, nil
	}
	const maxCollect = 10000
	collected := make([]snowflake.JsonInt64, 0, len(normalized))
	seen := make(map[int64]struct{}, len(normalized))
	for _, id := range normalized {
		if _, ok := seen[int64(id)]; !ok {
			seen[int64(id)] = struct{}{}
			collected = append(collected, id)
		}
		childIDs, err := s.collectChildIDs(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, childID := range childIDs {
			if _, ok := seen[int64(childID)]; ok {
				continue
			}
			seen[int64(childID)] = struct{}{}
			collected = append(collected, childID)
			if len(collected) > maxCollect {
				return nil, gerror.Newf("子树节点过多（超过 %d），请分批删除", maxCollect)
			}
		}
	}
	return collected, nil
}

// collectChildIDs 递归收集所有子节点 ID（最大深度 20 层防止无限递归）
func (s *s{{.ModelName}}) collectChildIDs(ctx context.Context, parentID snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	return s.doCollectChildIDs(ctx, parentID, 0)
}

func (s *s{{.ModelName}}) doCollectChildIDs(ctx context.Context, parentID snowflake.JsonInt64, depth int) ([]snowflake.JsonInt64, error) {
	if depth > 20 {
		return nil, gerror.New("子树深度超过 20 层上限，请检查数据完整性")
	}
	var childIDs []snowflake.JsonInt64
	m := dao.{{.DaoName}}.Ctx(ctx).
		Where(dao.{{.DaoName}}.Columns().ParentId, parentID).
		Where(dao.{{.DaoName}}.Columns().DeletedAt, nil)
{{- if .HasTenantScope}}
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}})
{{- end}}
	result, err := m.Fields(dao.{{.DaoName}}.Columns().Id).
		Array()
	if err != nil || len(result) == 0 {
		return childIDs, err
	}
	for _, v := range result {
		cid := snowflake.JsonInt64(v.Int64())
		childIDs = append(childIDs, cid)
		grandChildren, err := s.doCollectChildIDs(ctx, cid, depth+1)
		if err != nil {
			return childIDs, err
		}
		childIDs = append(childIDs, grandChildren...)
	}
	return childIDs, nil
}
{{- end}}
{{- if not .HasParentID}}

// BatchDelete 批量软删除{{.Comment}}
func (s *s{{.ModelName}}) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalize{{.ModelName}}IDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
{{- if .HasTenantScope}}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), normalizedIDs, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), normalizedIDs, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return err
	}
{{- end}}
	_, err := dao.{{.DaoName}}.Ctx(ctx).WhereIn(dao.{{.DaoName}}.Columns().Id, normalizedIDs).Delete()
{{- if .EnableOpLog}}
	if err == nil {
		oplog.Record(ctx, "{{.ModuleName}}", "batch-delete", fmt.Sprintf("%v", normalizedIDs), "")
	}
{{- end}}
	return err
}
{{- end}}

// Detail 获取{{.Comment}}详情
func (s *s{{.ModelName}}) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.{{.ModelName}}DetailOutput, err error) {
{{- if .HasTenantScope}}
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), id, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return nil, err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), id, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return nil, err
	}
{{- end}}
	out = &model.{{.ModelName}}DetailOutput{}
	err = dao.{{.DaoName}}.Ctx(ctx).Where(dao.{{.DaoName}}.Columns().Id, id).Where(dao.{{.DaoName}}.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("{{.Comment}}不存在或已删除")
	}
{{- range .Fields}}
{{- if and .RefFieldName (not .IsHidden)}}
	// 查询{{.Label}}关联显示
	if out.{{.NameCamel}} != 0 {
		refQuery := g.DB().Ctx(ctx).Model("{{.RefTableDB}}").Where("id", out.{{.NameCamel}})
{{- if .RefHasDeletedAt}}
		refQuery = refQuery.Where("deleted_at", nil)
{{- end}}
{{- if $.HasTenantScope}}
{{- if .RefHasTenantID}}
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", {{if .RefHasMerchantID}}"merchant_id"{{else}}""{{end}})
{{- end}}
{{- end}}
		val, err := refQuery.Value("{{.RefDisplayField}}")
		if err == nil {
			out.{{.RefFieldName}} = val.String()
		}
	}
{{- end}}
{{- end}}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *s{{.ModelName}}) applyListFilter(ctx context.Context, in *model.{{.ModelName}}ListInput) *gdb.Model {
	m := dao.{{.DaoName}}.Ctx(ctx).Where(dao.{{.DaoName}}.Columns().DeletedAt, nil)
{{- if .HasTenantScope}}
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}})
{{- end}}
{{- if .HasKeywordSearch}}
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
{{- range $index, $field := .KeywordSearchFields}}
{{- if eq $index 0}}
		keywordBuilder = keywordBuilder.WhereLike(dao.{{$.DaoName}}.Columns().{{.NameDao}}, "%"+in.Keyword+"%")
{{- else}}
		keywordBuilder = keywordBuilder.WhereOrLike(dao.{{$.DaoName}}.Columns().{{.NameDao}}, "%"+in.Keyword+"%")
{{- end}}
{{- end}}
		m = m.Where(keywordBuilder)
	}
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
	if in.{{.NameCamel}}Start != "" {
		m = m.WhereGTE(dao.{{$.DaoName}}.Columns().{{.NameDao}}, in.{{.NameCamel}}Start)
	}
	if in.{{.NameCamel}}End != "" {
		m = m.WhereLTE(dao.{{$.DaoName}}.Columns().{{.NameDao}}, in.{{.NameCamel}}End)
	}
{{- else if .SearchPointer}}
	if in.{{.NameCamel}} != nil {
		m = m.Where(dao.{{$.DaoName}}.Columns().{{.NameDao}}, *in.{{.NameCamel}})
	}
{{- else}}
	if in.{{.NameCamel}} != "" {
{{- if eq .SearchOperator "eq"}}
		m = m.Where(dao.{{$.DaoName}}.Columns().{{.NameDao}}, in.{{.NameCamel}})
{{- else}}
		m = m.WhereLike(dao.{{$.DaoName}}.Columns().{{.NameDao}}, "%"+in.{{.NameCamel}}+"%")
{{- end}}
	}
{{- end}}
{{- end}}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.{{.DaoName}}.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.{{.DaoName}}.Columns().CreatedAt, in.EndTime)
	}
{{- if or .HasCreatedBy .HasDeptID}}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, {{if .HasCreatedBy}}dao.{{.DaoName}}.Columns().CreatedBy{{else}}""{{end}}, {{if .HasDeptID}}dao.{{.DaoName}}.Columns().DeptId{{else}}""{{end}})
{{- end}}
	return m
}

{{- if $hasRef}}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *s{{.ModelName}}) fillRefFields(ctx context.Context, list []*model.{{.ModelName}}ListOutput) {
{{- range .Fields}}
{{- if and .RefFieldName (not .IsHidden)}}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.{{.NameCamel}} != 0 {
				idSet[int64(item.{{.NameCamel}})] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("{{.RefTableDB}}").
				Fields("id", "{{.RefDisplayField}}")
{{- if .RefHasDeletedAt}}
			refQuery = refQuery.Where("deleted_at", nil)
{{- end}}
{{- if $.HasTenantScope}}
{{- if .RefHasTenantID}}
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", {{if .RefHasMerchantID}}"merchant_id"{{else}}""{{end}})
{{- end}}
{{- end}}
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["{{.RefDisplayField}}"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.{{.NameCamel}})]; ok {
						item.{{.RefFieldName}} = val
					}
				}
			}
		}
	}
{{- end}}
{{- end}}
}
{{- end}}

// List 获取{{.Comment}}列表
func (s *s{{.ModelName}}) List(ctx context.Context, in *model.{{.ModelName}}ListInput) (list []*model.{{.ModelName}}ListOutput, total int, err error) {
	if in == nil {
		in = &model.{{.ModelName}}ListInput{}
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
{{- if $hasRef}}
	s.fillRefFields(ctx, list)
{{- end}}
	return
}

// isAllowedOrderField 校验排序字段是否在允许列表中
func (s *s{{.ModelName}}) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.{{.DaoName}}.Columns().Id:        true,
		dao.{{.DaoName}}.Columns().CreatedAt: true,
{{- if .HasSort}}
		dao.{{.DaoName}}.Columns().Sort:      true,
{{- end}}
{{- if .HasStatus}}
		dao.{{.DaoName}}.Columns().Status:    true,
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsPassword) (or .IsMoney .IsSearchable)}}
		dao.{{$.DaoName}}.Columns().{{.NameDao}}: true,
{{- end}}
{{- end}}
	}
	return allowed[field]
}

func (s *s{{.ModelName}}) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
{{- if .HasSort}}
	return m.OrderAsc(dao.{{.DaoName}}.Columns().Sort).OrderDesc(dao.{{.DaoName}}.Columns().Id)
{{- else}}
	return m.OrderDesc(dao.{{.DaoName}}.Columns().Id)
{{- end}}
}

// Export 导出{{.Comment}}（不分页）
func (s *s{{.ModelName}}) Export(ctx context.Context, in *model.{{.ModelName}}ListInput) (list []*model.{{.ModelName}}ListOutput, err error) {
	if in == nil {
		in = &model.{{.ModelName}}ListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
{{- if $hasRef}}
	s.fillRefFields(ctx, list)
{{- end}}
	return
}
{{- if .HasParentID}}

// Tree 获取{{.Comment}}树形结构
func (s *s{{.ModelName}}) Tree(ctx context.Context, in *model.{{.ModelName}}TreeInput) (tree []*model.{{.ModelName}}TreeOutput, err error) {
	var list []*model.{{.ModelName}}TreeOutput
	if in == nil {
		in = &model.{{.ModelName}}TreeInput{}
	}
	m := dao.{{.DaoName}}.Ctx(ctx).Where(dao.{{.DaoName}}.Columns().DeletedAt, nil)
{{- if .HasTenantScope}}
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}})
{{- end}}
{{- if .HasKeywordSearch}}
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
{{- range $index, $field := .KeywordSearchFields}}
{{- if eq $index 0}}
		keywordBuilder = keywordBuilder.WhereLike(dao.{{$.DaoName}}.Columns().{{.NameDao}}, "%"+in.Keyword+"%")
{{- else}}
		keywordBuilder = keywordBuilder.WhereOrLike(dao.{{$.DaoName}}.Columns().{{.NameDao}}, "%"+in.Keyword+"%")
{{- end}}
{{- end}}
		m = m.Where(keywordBuilder)
	}
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
	if in.{{.NameCamel}}Start != "" {
		m = m.WhereGTE(dao.{{$.DaoName}}.Columns().{{.NameDao}}, in.{{.NameCamel}}Start)
	}
	if in.{{.NameCamel}}End != "" {
		m = m.WhereLTE(dao.{{$.DaoName}}.Columns().{{.NameDao}}, in.{{.NameCamel}}End)
	}
{{- else if .SearchPointer}}
	if in.{{.NameCamel}} != nil {
		m = m.Where(dao.{{$.DaoName}}.Columns().{{.NameDao}}, *in.{{.NameCamel}})
	}
{{- else}}
	if in.{{.NameCamel}} != "" {
{{- if eq .SearchOperator "eq"}}
		m = m.Where(dao.{{$.DaoName}}.Columns().{{.NameDao}}, in.{{.NameCamel}})
{{- else}}
		m = m.WhereLike(dao.{{$.DaoName}}.Columns().{{.NameDao}}, "%"+in.{{.NameCamel}}+"%")
{{- end}}
	}
{{- end}}
{{- end}}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.{{.DaoName}}.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.{{.DaoName}}.Columns().CreatedAt, in.EndTime)
	}
{{- if or .HasCreatedBy .HasDeptID}}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, {{if .HasCreatedBy}}dao.{{.DaoName}}.Columns().CreatedBy{{else}}""{{end}}, {{if .HasDeptID}}dao.{{.DaoName}}.Columns().DeptId{{else}}""{{end}})
{{- end}}
	err = m.{{if .HasSort}}OrderAsc(dao.{{.DaoName}}.Columns().Sort).{{end}}Limit(5000).Scan(&list)
	if err != nil {
		return
	}

	// 使用 map 迭代方式组装树
	nodeMap := make(map[int64]*model.{{.ModelName}}TreeOutput, len(list))
	for _, item := range list {
		item.Children = make([]*model.{{.ModelName}}TreeOutput, 0)
		nodeMap[int64(item.ID)] = item
	}

	tree = make([]*model.{{.ModelName}}TreeOutput, 0)
	for _, item := range list {
		if int64(item.ParentID) == 0 {
			tree = append(tree, item)
		} else if parent, ok := nodeMap[int64(item.ParentID)]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			tree = append(tree, item)
		}
	}
{{- /* 填充非 parent_id 的外键关联字段 */}}
{{- range .Fields}}
{{- if and .RefFieldName (not .IsHidden) (not .IsParentID)}}
	// 批量填充{{.Label}}关联显示
	{
		idSet := make(map[int64]struct{})
		var collectIDs func(items []*model.{{$.ModelName}}TreeOutput)
		collectIDs = func(items []*model.{{$.ModelName}}TreeOutput) {
			for _, item := range items {
				if item.{{.NameCamel}} != 0 {
					idSet[int64(item.{{.NameCamel}})] = struct{}{}
				}
				if len(item.Children) > 0 {
					collectIDs(item.Children)
				}
			}
		}
		collectIDs(list)
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("{{.RefTableDB}}").
				Fields("id", "{{.RefDisplayField}}")
{{- if .RefHasDeletedAt}}
			refQuery = refQuery.Where("deleted_at", nil)
{{- end}}
{{- if $.HasTenantScope}}
{{- if .RefHasTenantID}}
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", {{if .RefHasMerchantID}}"merchant_id"{{else}}""{{end}})
{{- end}}
{{- end}}
			rows, queryErr := refQuery.WhereIn("id", ids).All()
			if queryErr == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["{{.RefDisplayField}}"].String()
				}
				var fillRef func(items []*model.{{$.ModelName}}TreeOutput)
				fillRef = func(items []*model.{{$.ModelName}}TreeOutput) {
					for _, item := range items {
						if val, ok := refMap[int64(item.{{.NameCamel}})]; ok {
							item.{{.RefFieldName}} = val
						}
						if len(item.Children) > 0 {
							fillRef(item.Children)
						}
					}
				}
				fillRef(list)
			}
		}
	}
{{- end}}
{{- end}}
	return
}
{{- end}}
{{- if .HasBatchEdit}}

// BatchUpdate 批量编辑{{.Comment}}
func (s *s{{.ModelName}}) BatchUpdate(ctx context.Context, in *model.{{.ModelName}}BatchUpdateInput) error {
	data := do.{{.DaoName}}{}
	hasChange := false
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (.IsEnum)}}
	if in.{{.NameCamel}} != nil {
		data.{{.NameDao}} = *in.{{.NameCamel}}
		hasChange = true
	}
{{- end}}
{{- end}}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalize{{.ModelName}}IDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
{{- if .HasTenantScope}}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), normalizedIDs, dao.{{.DaoName}}.Columns().Id, dao.{{.DaoName}}.Columns().TenantId, {{if .HasMerchantID}}dao.{{.DaoName}}.Columns().MerchantId{{else}}""{{end}}, "{{.Comment}}"); err != nil {
		return err
	}
{{- end}}
{{- if or .HasCreatedBy .HasDeptID}}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.{{.DaoName}}.Ctx(ctx), normalizedIDs, dao.{{.DaoName}}.Columns().Id{{if .HasCreatedBy}}, dao.{{.DaoName}}.Columns().CreatedBy{{else}}, ""{{end}}{{if .HasDeptID}}, dao.{{.DaoName}}.Columns().DeptId{{else}}, ""{{end}}); err != nil {
		return err
	}
{{- end}}
	_, err := dao.{{.DaoName}}.Ctx(ctx).WhereIn(dao.{{.DaoName}}.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}
{{- end}}
{{- if .HasImport}}

// Import 导入{{.Comment}}
func (s *s{{.ModelName}}) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.{{.DaoName}}{
			Id: id,
{{- if .HasCreatedBy}}
			CreatedBy: middleware.GetUserID(ctx),
{{- end}}
{{- if .HasDeptID}}
			DeptId: middleware.GetDeptID(ctx),
{{- end}}
		}
		idx := 0
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsPassword) (not .IsTimeField) (or (not $.HasTenantScope) (and (ne .Name "tenant_id") (ne .Name "merchant_id")))}}
{{- if .IsMoney}}
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.{{.NameDao}} = int64(math.Round(v * 100))
			}
		}
{{- else if eq .GoType "int"}}
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.{{.NameDao}} = v
			}
		}
{{- else if eq .GoType "int64"}}
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.{{.NameDao}} = v
			}
		}
{{- else if eq .GoType "float64"}}
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.{{.NameDao}} = v
			}
		}
{{- else if eq .GoType "JsonInt64"}}
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.{{.NameDao}} = v
			}
		}
{{- else}}
		if idx < len(record) {
			data.{{.NameDao}} = strings.TrimSpace(record[idx])
		}
{{- end}}
		idx++
{{- end}}
{{- end}}
{{- if .HasTenantScope}}
		tenantID := snowflake.JsonInt64(0)
		merchantID := snowflake.JsonInt64(0)
		middleware.ApplyTenantScopeToWrite(ctx, &tenantID{{if .HasMerchantID}}, &merchantID{{else}}, nil{{end}})
		if err := middleware.EnsureTenantMerchantAccessible(ctx, tenantID{{if .HasMerchantID}}, merchantID{{else}}, 0{{end}}); err != nil {
			fail++
			continue
		}
		data.TenantId = tenantID
{{- if .HasMerchantID}}
		data.MerchantId = merchantID
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and .IsForeignKey (not .IsHidden) (ne .Name "tenant_id") (ne .Name "merchant_id")}}
		if fkVal, ok := data.{{.NameDao}}.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("{{.RefTableDB}}").Where("id", fkVal)
{{- if .RefHasDeletedAt}}
			refQuery = refQuery.Where("deleted_at", nil)
{{- end}}
{{- if $.HasTenantScope}}
{{- if .RefHasTenantID}}
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", {{if .RefHasMerchantID}}"merchant_id"{{else}}""{{end}})
{{- end}}
{{- end}}
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
{{- end}}
{{- end}}
		if _, insertErr := dao.{{.DaoName}}.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}
{{end}}
