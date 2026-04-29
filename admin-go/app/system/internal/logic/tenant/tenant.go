package tenant

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
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

func init() {
	service.RegisterTenant(New())
}

func New() *sTenant {
	return &sTenant{}
}

type sTenant struct{}

type tenantCreateData struct {
	Id           snowflake.JsonInt64 `orm:"id"`
	Name         string              `orm:"name"`
	Code         string              `orm:"code"`
	ContactName  string              `orm:"contact_name"`
	ContactPhone string              `orm:"contact_phone"`
	Domain       string              `orm:"domain"`
	ExpireAt     any                 `orm:"expire_at"`
	Status       int                 `orm:"status"`
	Remark       string              `orm:"remark"`
	CreatedBy    int64               `orm:"created_by"`
	DeptId       int64               `orm:"dept_id"`
	TenantId     snowflake.JsonInt64 `orm:"tenant_id"`
	MerchantId   snowflake.JsonInt64 `orm:"merchant_id"`
}

func (s *sTenant) Create(ctx context.Context, in *model.TenantCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	if err := ensurePlatformTenantScope(ctx); err != nil {
		return err
	}
	normalizeTenantCreateInput(in)
	if err := validateTenantFields(in.Name, in.Code, in.Status); err != nil {
		return err
	}
	if err := s.ensureCodeUnique(ctx, 0, in.Code); err != nil {
		return err
	}
	tenantID := snowflake.Generate()
	return dao.Tenant.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Tenant.Table()).Ctx(ctx).Data(tenantCreateData{
			Id:           tenantID,
			Name:         in.Name,
			Code:         in.Code,
			ContactName:  in.ContactName,
			ContactPhone: in.ContactPhone,
			Domain:       in.Domain,
			ExpireAt:     in.ExpireAt,
			Status:       in.Status,
			Remark:       in.Remark,
			CreatedBy:    shared.CurrentActorUserID(ctx),
			DeptId:       shared.CurrentActorDeptID(ctx),
			TenantId:     tenantID,
			MerchantId:   0,
		}).Insert(); err != nil {
			return err
		}
		if in.CreateAdmin != 1 {
			return nil
		}
		return shared.BootstrapScopedAdmin(ctx, tx, shared.AdminBootstrapInput{
			TenantID:        tenantID,
			MerchantID:      0,
			DeptParentID:    0,
			DeptTitle:       defaultTenantAdminDeptTitle(in.Name),
			DeptManagerName: in.ContactName,
			RoleTitle:       "租户管理员",
			AdminUsername:   in.AdminUsername,
			AdminPassword:   in.AdminPassword,
			AdminNickname:   defaultAdminNickname(in.AdminNickname, in.Name),
			AdminEmail:      in.AdminEmail,
			MenuProfile:     shared.AdminBootstrapMenuProfileTenant,
			CreatedBy:       shared.CurrentActorUserID(ctx),
		})
	})
}

func (s *sTenant) Update(ctx context.Context, in *model.TenantUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	if err := ensurePlatformTenantScope(ctx); err != nil {
		return err
	}
	normalizeTenantUpdateInput(in)
	if err := s.ensureExists(ctx, in.ID); err != nil {
		return err
	}
	if err := validateTenantFields(in.Name, in.Code, in.Status); err != nil {
		return err
	}
	if err := s.ensureCodeUnique(ctx, in.ID, in.Code); err != nil {
		return err
	}
	_, err := dao.Tenant.Ctx(ctx).
		Where(dao.Tenant.Columns().Id, in.ID).
		Where(dao.Tenant.Columns().DeletedAt, nil).
		Data(do.Tenant{
			Name:         in.Name,
			Code:         in.Code,
			ContactName:  in.ContactName,
			ContactPhone: in.ContactPhone,
			Domain:       in.Domain,
			ExpireAt:     in.ExpireAt,
			Status:       in.Status,
			Remark:       in.Remark,
		}).
		Update()
	return err
}

func (s *sTenant) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := ensurePlatformTenantScope(ctx); err != nil {
		return err
	}
	if err := s.ensureExists(ctx, id); err != nil {
		return err
	}
	if err := s.ensureDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.Tenant.Ctx(ctx).
		Where(dao.Tenant.Columns().Id, id).
		Delete()
	return err
}

func (s *sTenant) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	if err := ensurePlatformTenantScope(ctx); err != nil {
		return err
	}
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的租户")
	}
	if err := s.ensureIDsExist(ctx, ids); err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.ensureDeletable(ctx, id); err != nil {
			return err
		}
	}
	_, err := dao.Tenant.Ctx(ctx).
		WhereIn(dao.Tenant.Columns().Id, batchutil.ToInt64s(ids)).
		Delete()
	return err
}

func (s *sTenant) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.TenantDetailOutput, err error) {
	if err := ensurePlatformTenantScope(ctx); err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, gerror.New("租户不存在或已删除")
	}
	out = &model.TenantDetailOutput{}
	if err = dao.Tenant.Ctx(ctx).
		Where(dao.Tenant.Columns().Id, id).
		Where(dao.Tenant.Columns().DeletedAt, nil).
		Scan(out); err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("租户不存在或已删除")
	}
	return out, nil
}

func (s *sTenant) List(ctx context.Context, in *model.TenantListInput) (list []*model.TenantListOutput, total int, err error) {
	if err := ensurePlatformTenantScope(ctx); err != nil {
		return nil, 0, err
	}
	if in == nil {
		in = &model.TenantListInput{}
	}
	normalizeTenantListInput(in)
	m := dao.Tenant.Ctx(ctx).Where(dao.Tenant.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Tenant.Columns().Name, "%"+in.Keyword+"%").
			WhereOrLike(dao.Tenant.Columns().Code, "%"+in.Keyword+"%").
			WhereOrLike(dao.Tenant.Columns().ContactName, "%"+in.Keyword+"%").
			WhereOrLike(dao.Tenant.Columns().ContactPhone, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Code != "" {
		m = m.WhereLike(dao.Tenant.Columns().Code, "%"+in.Code+"%")
	}
	if in.Status != nil {
		m = m.Where(dao.Tenant.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(dao.Tenant.Columns().Id).Scan(&list)
	return list, total, err
}

func normalizeTenantCreateInput(in *model.TenantCreateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)
	in.ContactName = strings.TrimSpace(in.ContactName)
	in.ContactPhone = strings.TrimSpace(in.ContactPhone)
	in.Domain = strings.TrimSpace(in.Domain)
	in.Remark = strings.TrimSpace(in.Remark)
	in.AdminUsername = strings.TrimSpace(in.AdminUsername)
	if in.CreateAdmin == 1 && in.AdminUsername == "" && in.Code != "" {
		in.AdminUsername = in.Code + "_admin"
	}
	in.AdminPassword = strings.TrimSpace(in.AdminPassword)
	in.AdminNickname = strings.TrimSpace(in.AdminNickname)
	in.AdminEmail = strings.TrimSpace(in.AdminEmail)
}

func normalizeTenantUpdateInput(in *model.TenantUpdateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)
	in.ContactName = strings.TrimSpace(in.ContactName)
	in.ContactPhone = strings.TrimSpace(in.ContactPhone)
	in.Domain = strings.TrimSpace(in.Domain)
	in.Remark = strings.TrimSpace(in.Remark)
}

func normalizeTenantListInput(in *model.TenantListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Code = strings.TrimSpace(in.Code)
}

func validateTenantFields(name, code string, status int) error {
	if strings.TrimSpace(name) == "" {
		return gerror.New("租户名称不能为空")
	}
	if strings.TrimSpace(code) == "" {
		return gerror.New("租户编码不能为空")
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	return nil
}

func defaultTenantAdminDeptTitle(tenantName string) string {
	tenantName = strings.TrimSpace(tenantName)
	if tenantName == "" {
		return "默认部门"
	}
	return tenantName + "默认部门"
}

func defaultAdminNickname(nickname, ownerName string) string {
	nickname = strings.TrimSpace(nickname)
	if nickname != "" {
		return nickname
	}
	ownerName = strings.TrimSpace(ownerName)
	if ownerName == "" {
		return "管理员"
	}
	return ownerName + "管理员"
}

func ensurePlatformTenantScope(ctx context.Context) error {
	if shared.ResolveTenantAccessScope(ctx).All {
		return nil
	}
	return gerror.New("仅平台账号可操作租户")
}

func (s *sTenant) ensureExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("租户不存在或已删除")
	}
	count, err := dao.Tenant.Ctx(ctx).
		Where(dao.Tenant.Columns().Id, id).
		Where(dao.Tenant.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("租户不存在或已删除")
	}
	return nil
}

func (s *sTenant) ensureIDsExist(ctx context.Context, ids []snowflake.JsonInt64) error {
	dbIDs := batchutil.ToInt64s(ids)
	count, err := dao.Tenant.Ctx(ctx).
		WhereIn(dao.Tenant.Columns().Id, dbIDs).
		Where(dao.Tenant.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count != len(dbIDs) {
		return gerror.New("包含不存在或已删除的租户")
	}
	return nil
}

func (s *sTenant) ensureCodeUnique(ctx context.Context, currentID snowflake.JsonInt64, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil
	}
	m := dao.Tenant.Ctx(ctx).
		Where(dao.Tenant.Columns().Code, code).
		Where(dao.Tenant.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Tenant.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("租户编码已存在")
	}
	return nil
}

func (s *sTenant) ensureDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	tables := []struct {
		name    string
		column  string
		message string
	}{
		{dao.Merchant.Table(), dao.Merchant.Columns().TenantId, "当前租户下存在商户，不能删除"},
		{dao.Users.Table(), dao.Users.Columns().TenantId, "当前租户下存在用户，不能删除"},
		{dao.Dept.Table(), dao.Dept.Columns().TenantId, "当前租户下存在部门，不能删除"},
		{dao.Role.Table(), dao.Role.Columns().TenantId, "当前租户下存在角色，不能删除"},
		{dao.Domain.Table(), dao.Domain.Columns().TenantId, "当前租户下存在域名绑定，不能删除"},
	}
	for _, item := range tables {
		count, err := dao.Tenant.DB().Model(item.name).Ctx(ctx).
			Where(item.column, id).
			Where("deleted_at", nil).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return gerror.New(item.message)
		}
	}
	return nil
}
