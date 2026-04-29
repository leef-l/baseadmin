package merchant

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
	service.RegisterMerchant(New())
}

func New() *sMerchant {
	return &sMerchant{}
}

type sMerchant struct{}

type merchantCreateData struct {
	Id           snowflake.JsonInt64 `orm:"id"`
	TenantId     snowflake.JsonInt64 `orm:"tenant_id"`
	MerchantId   snowflake.JsonInt64 `orm:"merchant_id"`
	Name         string              `orm:"name"`
	Code         string              `orm:"code"`
	ContactName  string              `orm:"contact_name"`
	ContactPhone string              `orm:"contact_phone"`
	Address      string              `orm:"address"`
	Status       int                 `orm:"status"`
	Remark       string              `orm:"remark"`
	CreatedBy    int64               `orm:"created_by"`
	DeptId       int64               `orm:"dept_id"`
}

func (s *sMerchant) Create(ctx context.Context, in *model.MerchantCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeMerchantCreateInput(in)
	s.applyActorTenant(ctx, &in.TenantID)
	if err := validateMerchantFields(in.TenantID, in.Name, in.Code, in.Status); err != nil {
		return err
	}
	if err := s.ensureTenantAccessible(ctx, in.TenantID); err != nil {
		return err
	}
	if err := s.ensureCodeUnique(ctx, 0, in.TenantID, in.Code); err != nil {
		return err
	}
	merchantID := snowflake.Generate()
	deptParentID := snowflake.JsonInt64(0)
	if in.CreateAdmin == 1 {
		parentID, err := s.lookupTenantRootDeptID(ctx, in.TenantID)
		if err != nil {
			return err
		}
		deptParentID = parentID
	}
	return dao.Merchant.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Merchant.Table()).Ctx(ctx).Data(merchantCreateData{
			Id:           merchantID,
			TenantId:     in.TenantID,
			MerchantId:   merchantID,
			Name:         in.Name,
			Code:         in.Code,
			ContactName:  in.ContactName,
			ContactPhone: in.ContactPhone,
			Address:      in.Address,
			Status:       in.Status,
			Remark:       in.Remark,
			CreatedBy:    shared.CurrentActorUserID(ctx),
			DeptId:       shared.CurrentActorDeptID(ctx),
		}).Insert(); err != nil {
			return err
		}
		if in.CreateAdmin != 1 {
			return nil
		}
		return shared.BootstrapScopedAdmin(ctx, tx, shared.AdminBootstrapInput{
			TenantID:        in.TenantID,
			MerchantID:      merchantID,
			DeptParentID:    deptParentID,
			DeptTitle:       defaultMerchantAdminDeptTitle(in.Name),
			DeptManagerName: in.ContactName,
			RoleTitle:       "商户管理员",
			AdminUsername:   in.AdminUsername,
			AdminPassword:   in.AdminPassword,
			AdminNickname:   defaultAdminNickname(in.AdminNickname, in.Name),
			AdminEmail:      in.AdminEmail,
			MenuProfile:     shared.AdminBootstrapMenuProfileMerchant,
			CreatedBy:       shared.CurrentActorUserID(ctx),
		})
	})
}

func (s *sMerchant) Update(ctx context.Context, in *model.MerchantUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeMerchantUpdateInput(in)
	if err := s.ensureAccessible(ctx, in.ID); err != nil {
		return err
	}
	s.applyActorTenant(ctx, &in.TenantID)
	if err := validateMerchantFields(in.TenantID, in.Name, in.Code, in.Status); err != nil {
		return err
	}
	if err := s.ensureTenantAccessible(ctx, in.TenantID); err != nil {
		return err
	}
	if err := s.ensureCodeUnique(ctx, in.ID, in.TenantID, in.Code); err != nil {
		return err
	}
	_, err := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().Id, in.ID).
		Where(dao.Merchant.Columns().DeletedAt, nil).
		Data(do.Merchant{
			TenantId:     in.TenantID,
			Name:         in.Name,
			Code:         in.Code,
			ContactName:  in.ContactName,
			ContactPhone: in.ContactPhone,
			Address:      in.Address,
			Status:       in.Status,
			Remark:       in.Remark,
		}).
		Update()
	return err
}

func (s *sMerchant) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureAccessible(ctx, id); err != nil {
		return err
	}
	if err := s.ensureDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().Id, id).
		Delete()
	return err
}

func (s *sMerchant) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的商户")
	}
	for _, id := range ids {
		if err := s.ensureAccessible(ctx, id); err != nil {
			return err
		}
		if err := s.ensureDeletable(ctx, id); err != nil {
			return err
		}
	}
	_, err := dao.Merchant.Ctx(ctx).
		WhereIn(dao.Merchant.Columns().Id, batchutil.ToInt64s(ids)).
		Delete()
	return err
}

func (s *sMerchant) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.MerchantDetailOutput, err error) {
	if err := s.ensureAccessible(ctx, id); err != nil {
		return nil, err
	}
	out = &model.MerchantDetailOutput{}
	if err = dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().Id, id).
		Where(dao.Merchant.Columns().DeletedAt, nil).
		Scan(out); err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("商户不存在或已删除")
	}
	out.TenantName = s.lookupTenantName(ctx, int64(out.TenantID))
	return out, nil
}

func (s *sMerchant) List(ctx context.Context, in *model.MerchantListInput) (list []*model.MerchantListOutput, total int, err error) {
	if in == nil {
		in = &model.MerchantListInput{}
	}
	normalizeMerchantListInput(in)
	m := dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().DeletedAt, nil)
	m = shared.ApplyTenantScopeToModel(ctx, m, dao.Merchant.Columns().TenantId, dao.Merchant.Columns().Id)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Merchant.Columns().Name, "%"+in.Keyword+"%").
			WhereOrLike(dao.Merchant.Columns().Code, "%"+in.Keyword+"%").
			WhereOrLike(dao.Merchant.Columns().ContactName, "%"+in.Keyword+"%").
			WhereOrLike(dao.Merchant.Columns().ContactPhone, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.TenantID > 0 {
		m = m.Where(dao.Merchant.Columns().TenantId, in.TenantID)
	}
	if in.Code != "" {
		m = m.WhereLike(dao.Merchant.Columns().Code, "%"+in.Code+"%")
	}
	if in.Status != nil {
		m = m.Where(dao.Merchant.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(dao.Merchant.Columns().Id).Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	s.fillTenantNames(ctx, list)
	return list, total, nil
}

func normalizeMerchantCreateInput(in *model.MerchantCreateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)
	in.ContactName = strings.TrimSpace(in.ContactName)
	in.ContactPhone = strings.TrimSpace(in.ContactPhone)
	in.Address = strings.TrimSpace(in.Address)
	in.Remark = strings.TrimSpace(in.Remark)
	in.AdminUsername = strings.TrimSpace(in.AdminUsername)
	if in.CreateAdmin == 1 && in.AdminUsername == "" && in.Code != "" {
		in.AdminUsername = in.Code + "_admin"
	}
	in.AdminPassword = strings.TrimSpace(in.AdminPassword)
	in.AdminNickname = strings.TrimSpace(in.AdminNickname)
	in.AdminEmail = strings.TrimSpace(in.AdminEmail)
}

func normalizeMerchantUpdateInput(in *model.MerchantUpdateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)
	in.ContactName = strings.TrimSpace(in.ContactName)
	in.ContactPhone = strings.TrimSpace(in.ContactPhone)
	in.Address = strings.TrimSpace(in.Address)
	in.Remark = strings.TrimSpace(in.Remark)
}

func normalizeMerchantListInput(in *model.MerchantListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Code = strings.TrimSpace(in.Code)
}

func validateMerchantFields(tenantID snowflake.JsonInt64, name, code string, status int) error {
	if tenantID <= 0 {
		return gerror.New("租户不能为空")
	}
	if strings.TrimSpace(name) == "" {
		return gerror.New("商户名称不能为空")
	}
	if strings.TrimSpace(code) == "" {
		return gerror.New("商户编码不能为空")
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	return nil
}

func defaultMerchantAdminDeptTitle(merchantName string) string {
	merchantName = strings.TrimSpace(merchantName)
	if merchantName == "" {
		return "默认部门"
	}
	return merchantName + "默认部门"
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

func (s *sMerchant) applyActorTenant(ctx context.Context, tenantID *snowflake.JsonInt64) {
	scope := shared.ResolveTenantAccessScope(ctx)
	if scope.All || tenantID == nil {
		return
	}
	*tenantID = snowflake.JsonInt64(scope.TenantID)
}

func (s *sMerchant) ensureTenantAccessible(ctx context.Context, tenantID snowflake.JsonInt64) error {
	if tenantID <= 0 {
		return gerror.New("租户不存在或已删除")
	}
	if !shared.CanAccessTenantMerchant(ctx, int64(tenantID), 0) {
		return gerror.New("租户不存在或无权操作")
	}
	count, err := dao.Tenant.Ctx(ctx).
		Where(dao.Tenant.Columns().Id, tenantID).
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

func (s *sMerchant) lookupTenantRootDeptID(ctx context.Context, tenantID snowflake.JsonInt64) (snowflake.JsonInt64, error) {
	if tenantID <= 0 {
		return 0, nil
	}
	var row struct {
		Id int64 `json:"id"`
	}
	err := dao.Dept.Ctx(ctx).
		Fields(dao.Dept.Columns().Id).
		Where(dao.Dept.Columns().TenantId, tenantID).
		Where(dao.Dept.Columns().MerchantId, 0).
		Where(dao.Dept.Columns().ParentId, 0).
		Where(dao.Dept.Columns().DeletedAt, nil).
		OrderAsc(dao.Dept.Columns().Id).
		Scan(&row)
	if err != nil {
		return 0, err
	}
	return snowflake.JsonInt64(row.Id), nil
}

func (s *sMerchant) ensureAccessible(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("商户不存在或已删除")
	}
	var row struct {
		Id       int64 `json:"id"`
		TenantId int64 `json:"tenantId"`
	}
	if err := dao.Merchant.Ctx(ctx).
		Fields(dao.Merchant.Columns().Id, dao.Merchant.Columns().TenantId).
		Where(dao.Merchant.Columns().Id, id).
		Where(dao.Merchant.Columns().DeletedAt, nil).
		Scan(&row); err != nil {
		return err
	}
	if row.Id == 0 {
		return gerror.New("商户不存在或已删除")
	}
	if !shared.CanAccessTenantMerchant(ctx, row.TenantId, row.Id) {
		return gerror.New("商户不存在或无权操作")
	}
	return nil
}

func (s *sMerchant) ensureCodeUnique(ctx context.Context, currentID, tenantID snowflake.JsonInt64, code string) error {
	code = strings.TrimSpace(code)
	if code == "" || tenantID <= 0 {
		return nil
	}
	m := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().TenantId, tenantID).
		Where(dao.Merchant.Columns().Code, code).
		Where(dao.Merchant.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Merchant.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("同租户下商户编码已存在")
	}
	return nil
}

func (s *sMerchant) ensureDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	tables := []struct {
		name    string
		column  string
		message string
	}{
		{dao.Users.Table(), dao.Users.Columns().MerchantId, "当前商户下存在用户，不能删除"},
		{dao.Dept.Table(), dao.Dept.Columns().MerchantId, "当前商户下存在部门，不能删除"},
		{dao.Role.Table(), dao.Role.Columns().MerchantId, "当前商户下存在角色，不能删除"},
		{dao.Domain.Table(), dao.Domain.Columns().MerchantId, "当前商户下存在域名绑定，不能删除"},
	}
	for _, item := range tables {
		count, err := dao.Merchant.DB().Model(item.name).Ctx(ctx).
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

func (s *sMerchant) fillTenantNames(ctx context.Context, list []*model.MerchantListOutput) {
	tenantIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.TenantID > 0 {
			tenantIDs = append(tenantIDs, int64(item.TenantID))
		}
	}
	tenantMap := s.loadTenantNameMap(ctx, tenantIDs)
	for _, item := range list {
		item.TenantName = tenantMap[int64(item.TenantID)]
	}
}

func (s *sMerchant) lookupTenantName(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}
	return s.loadTenantNameMap(ctx, []int64{id})[id]
}

func (s *sMerchant) loadTenantNameMap(ctx context.Context, ids []int64) map[int64]string {
	seen := make(map[int64]struct{}, len(ids))
	normalized := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}
	if len(normalized) == 0 {
		return nil
	}
	var rows []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := dao.Tenant.Ctx(ctx).
		Fields(dao.Tenant.Columns().Id, dao.Tenant.Columns().Name).
		WhereIn(dao.Tenant.Columns().Id, normalized).
		Where(dao.Tenant.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
		return nil
	}
	out := make(map[int64]string, len(rows))
	for _, row := range rows {
		out[row.Id] = row.Name
	}
	return out
}
