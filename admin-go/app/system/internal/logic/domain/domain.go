package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	ownerTypeTenant   = 1
	ownerTypeMerchant = 2
	defaultAppCode    = "admin"
)

func init() {
	service.RegisterDomain(New())
}

func New() *sDomain {
	return &sDomain{}
}

type sDomain struct{}

func (s *sDomain) Create(ctx context.Context, in *model.DomainCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDomainCreateInput(in)
	if err := s.validateDomainWrite(ctx, 0, in.Domain, in.AppCode, in.OwnerType, &in.TenantID, &in.MerchantID, in.VerifyStatus, in.SslStatus, in.Status); err != nil {
		return err
	}
	_, err := dao.Domain.Ctx(ctx).Data(do.Domain{
		Id:           snowflake.Generate(),
		Domain:       in.Domain,
		OwnerType:    in.OwnerType,
		TenantId:     in.TenantID,
		MerchantId:   in.MerchantID,
		AppCode:      in.AppCode,
		VerifyToken:  defaultVerifyToken(),
		VerifyStatus: in.VerifyStatus,
		SslStatus:    in.SslStatus,
		NginxStatus:  0,
		Status:       in.Status,
		Remark:       in.Remark,
		CreatedBy:    shared.CurrentActorUserID(ctx),
		DeptId:       shared.CurrentActorDeptID(ctx),
	}).Insert()
	return err
}

func (s *sDomain) Update(ctx context.Context, in *model.DomainUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDomainUpdateInput(in)
	if err := s.ensureAccessible(ctx, in.ID); err != nil {
		return err
	}
	if err := s.validateDomainWrite(ctx, in.ID, in.Domain, in.AppCode, in.OwnerType, &in.TenantID, &in.MerchantID, in.VerifyStatus, in.SslStatus, in.Status); err != nil {
		return err
	}
	_, err := dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Id, in.ID).
		Where(dao.Domain.Columns().DeletedAt, nil).
		Data(do.Domain{
			Domain:       in.Domain,
			OwnerType:    in.OwnerType,
			TenantId:     in.TenantID,
			MerchantId:   in.MerchantID,
			AppCode:      in.AppCode,
			VerifyStatus: in.VerifyStatus,
			SslStatus:    in.SslStatus,
			NginxStatus:  0,
			Status:       in.Status,
			Remark:       in.Remark,
		}).
		Update()
	return err
}

func (s *sDomain) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureAccessible(ctx, id); err != nil {
		return err
	}
	_, err := dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Id, id).
		Delete()
	return err
}

func (s *sDomain) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的域名")
	}
	if err := shared.EnsureTenantScopedRowsAccessible(ctx, dao.Domain.Ctx(ctx), ids, dao.Domain.Columns().Id, dao.Domain.Columns().TenantId, dao.Domain.Columns().MerchantId, "域名"); err != nil {
		return err
	}
	_, err := dao.Domain.Ctx(ctx).
		WhereIn(dao.Domain.Columns().Id, batchutil.ToInt64s(ids)).
		Delete()
	return err
}

func (s *sDomain) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DomainDetailOutput, err error) {
	if err := s.ensureAccessible(ctx, id); err != nil {
		return nil, err
	}
	out = &model.DomainDetailOutput{}
	if err = dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Id, id).
		Where(dao.Domain.Columns().DeletedAt, nil).
		Scan(out); err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("域名不存在或已删除")
	}
	s.fillNames(ctx, []*model.DomainListOutput{out})
	return out, nil
}

func (s *sDomain) List(ctx context.Context, in *model.DomainListInput) (list []*model.DomainListOutput, total int, err error) {
	if in == nil {
		in = &model.DomainListInput{}
	}
	normalizeDomainListInput(in)
	m := dao.Domain.Ctx(ctx).Where(dao.Domain.Columns().DeletedAt, nil)
	m = shared.ApplyTenantScopeToModel(ctx, m, dao.Domain.Columns().TenantId, dao.Domain.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Domain.Columns().Domain, "%"+in.Keyword+"%").
			WhereOrLike(dao.Domain.Columns().AppCode, "%"+in.Keyword+"%").
			WhereOrLike(dao.Domain.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Domain != "" {
		m = m.WhereLike(dao.Domain.Columns().Domain, "%"+in.Domain+"%")
	}
	if in.OwnerType > 0 {
		m = m.Where(dao.Domain.Columns().OwnerType, in.OwnerType)
	}
	if in.TenantID > 0 {
		m = m.Where(dao.Domain.Columns().TenantId, in.TenantID)
	}
	if in.MerchantID > 0 {
		m = m.Where(dao.Domain.Columns().MerchantId, in.MerchantID)
	}
	if in.AppCode != "" {
		m = m.Where(dao.Domain.Columns().AppCode, in.AppCode)
	}
	if in.Status != nil {
		m = m.Where(dao.Domain.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(dao.Domain.Columns().Id).Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	s.fillNames(ctx, list)
	return list, total, nil
}

func (s *sDomain) ApplyNginx(ctx context.Context, id snowflake.JsonInt64) (out *model.DomainApplyNginxOutput, err error) {
	if !shared.ResolveTenantAccessScope(ctx).All {
		return nil, gerror.New("仅平台账号可应用Nginx配置")
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil || row.ID == 0 {
		return nil, gerror.New("域名不存在或已删除")
	}
	if row.Status != 1 {
		return nil, gerror.New("域名已停用，不能应用Nginx配置")
	}
	if row.VerifyStatus != 1 {
		return nil, gerror.New("域名未校验，不能应用Nginx配置")
	}
	if row.AppCode != defaultAppCode {
		return nil, gerror.New("当前仅支持后台应用域名自动配置")
	}
	applyOut, err := applyNginxConfig(ctx, row)
	if err != nil {
		return nil, err
	}
	if _, err = dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Id, id).
		Data(do.Domain{
			NginxStatus: applyOut.NginxStatus,
			SslStatus:   applyOut.SslStatus,
		}).
		Update(); err != nil {
		return nil, err
	}
	return applyOut, nil
}

func (s *sDomain) ApplySSL(ctx context.Context, id snowflake.JsonInt64) (out *model.DomainApplySSLOutput, err error) {
	if !shared.ResolveTenantAccessScope(ctx).All {
		return nil, gerror.New("仅平台账号可申请SSL证书")
	}
	row, err := s.loadRow(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil || row.ID == 0 {
		return nil, gerror.New("域名不存在或已删除")
	}
	if row.Status != 1 {
		return nil, gerror.New("域名已停用，不能申请SSL证书")
	}
	if row.VerifyStatus != 1 {
		return nil, gerror.New("域名未校验，不能申请SSL证书")
	}
	if row.AppCode != defaultAppCode {
		return nil, gerror.New("当前仅支持后台应用域名自动申请SSL")
	}

	nginxOut, err := applyNginxConfig(ctx, row)
	if err != nil {
		return nil, err
	}
	certPath, err := applyBaoTaSSL(ctx, row.Domain)
	if err != nil {
		return nil, err
	}
	nginxOut, err = applyNginxConfig(ctx, row)
	if err != nil {
		return nil, err
	}
	if _, err = dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Id, id).
		Data(do.Domain{
			NginxStatus: nginxOut.NginxStatus,
			SslStatus:   1,
		}).
		Update(); err != nil {
		return nil, err
	}
	return &model.DomainApplySSLOutput{
		ConfigPath:  nginxOut.ConfigPath,
		CertPath:    certPath,
		NginxStatus: nginxOut.NginxStatus,
		SslStatus:   1,
	}, nil
}

type domainRow struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Domain       string              `json:"domain"`
	OwnerType    int                 `json:"ownerType"`
	TenantID     snowflake.JsonInt64 `json:"tenantId"`
	MerchantID   snowflake.JsonInt64 `json:"merchantId"`
	AppCode      string              `json:"appCode"`
	VerifyStatus int                 `json:"verifyStatus"`
	SslStatus    int                 `json:"sslStatus"`
	NginxStatus  int                 `json:"nginxStatus"`
	Status       int                 `json:"status"`
}

func normalizeDomainCreateInput(in *model.DomainCreateInput) {
	if in == nil {
		return
	}
	in.Domain = normalizeDomainName(in.Domain)
	in.AppCode = normalizeAppCode(in.AppCode)
	if in.AppCode == "" {
		in.AppCode = defaultAppCode
	}
	if in.OwnerType == 0 {
		in.OwnerType = ownerTypeTenant
	}
	in.Remark = strings.TrimSpace(in.Remark)
}

func normalizeDomainUpdateInput(in *model.DomainUpdateInput) {
	if in == nil {
		return
	}
	in.Domain = normalizeDomainName(in.Domain)
	in.AppCode = normalizeAppCode(in.AppCode)
	if in.AppCode == "" {
		in.AppCode = defaultAppCode
	}
	if in.OwnerType == 0 {
		in.OwnerType = ownerTypeTenant
	}
	in.Remark = strings.TrimSpace(in.Remark)
}

func normalizeDomainListInput(in *model.DomainListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Domain = normalizeDomainName(in.Domain)
	in.AppCode = normalizeAppCode(in.AppCode)
}

func (s *sDomain) validateDomainWrite(
	ctx context.Context,
	currentID snowflake.JsonInt64,
	domainName string,
	appCode string,
	ownerType int,
	tenantID *snowflake.JsonInt64,
	merchantID *snowflake.JsonInt64,
	verifyStatus int,
	sslStatus int,
	status int,
) error {
	if err := validateDomainName(domainName); err != nil {
		return err
	}
	if err := validateAppCode(appCode); err != nil {
		return err
	}
	if err := validateOwnerType(ownerType); err != nil {
		return err
	}
	if err := fieldvalid.Binary("校验状态", verifyStatus); err != nil {
		return err
	}
	if err := fieldvalid.Binary("SSL状态", sslStatus); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	if err := s.resolveOwnership(ctx, ownerType, tenantID, merchantID); err != nil {
		return err
	}
	if currentHost := shared.CurrentRequestHost(ctx); currentHost != "" && currentHost == domainName {
		return gerror.New("不能把当前平台访问域名绑定为租户或商户域名")
	}
	return s.ensureDomainUnique(ctx, currentID, domainName, appCode)
}

func (s *sDomain) resolveOwnership(ctx context.Context, ownerType int, tenantID, merchantID *snowflake.JsonInt64) error {
	if tenantID == nil || merchantID == nil {
		return gerror.New("租户或商户参数不能为空")
	}
	scope := shared.ResolveTenantAccessScope(ctx)
	if !scope.All {
		*tenantID = snowflake.JsonInt64(scope.TenantID)
		if scope.MerchantID > 0 {
			if ownerType != ownerTypeMerchant {
				return gerror.New("商户账号只能绑定商户域名")
			}
			*merchantID = snowflake.JsonInt64(scope.MerchantID)
			return shared.EnsureTenantMerchantAccessible(ctx, *tenantID, *merchantID)
		}
		if ownerType == ownerTypeTenant {
			*merchantID = 0
			return shared.EnsureTenantMerchantAccessible(ctx, *tenantID, 0)
		}
		if *merchantID <= 0 {
			return gerror.New("商户域名必须选择商户")
		}
		return shared.EnsureTenantMerchantAccessible(ctx, *tenantID, *merchantID)
	}

	if *tenantID <= 0 {
		return gerror.New("租户不能为空")
	}
	if ownerType == ownerTypeTenant {
		*merchantID = 0
		return shared.EnsureTenantMerchantAccessible(ctx, *tenantID, 0)
	}
	if *merchantID <= 0 {
		return gerror.New("商户域名必须选择商户")
	}
	return shared.EnsureTenantMerchantAccessible(ctx, *tenantID, *merchantID)
}

func validateOwnerType(value int) error {
	switch value {
	case ownerTypeTenant, ownerTypeMerchant:
		return nil
	default:
		return gerror.New("主体类型不正确")
	}
}

func validateDomainName(value string) error {
	if value == "" {
		return gerror.New("域名不能为空")
	}
	if len(value) > 255 {
		return gerror.New("域名长度不能超过255位")
	}
	if strings.Contains(value, "*") || strings.Contains(value, "_") || strings.ContainsAny(value, " \t\r\n") {
		return gerror.New("域名格式不正确")
	}
	if strings.Count(value, ".") == 0 {
		return gerror.New("域名格式不正确")
	}
	labels := strings.Split(value, ".")
	for _, label := range labels {
		if label == "" || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return gerror.New("域名格式不正确")
		}
		for _, r := range label {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				continue
			}
			return gerror.New("域名格式不正确")
		}
	}
	return nil
}

func validateAppCode(value string) error {
	if value == "" {
		return gerror.New("应用编码不能为空")
	}
	if len(value) > 50 {
		return gerror.New("应用编码长度不能超过50位")
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return gerror.New("应用编码格式不正确")
	}
	return nil
}

func normalizeDomainName(value string) string {
	return shared.NormalizeDomainHost(value)
}

func normalizeAppCode(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func defaultVerifyToken() string {
	return fmt.Sprintf("baseadmin-%d", snowflake.Generate())
}

func (s *sDomain) ensureDomainUnique(ctx context.Context, currentID snowflake.JsonInt64, domainName, appCode string) error {
	m := dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Domain, domainName).
		Where(dao.Domain.Columns().AppCode, appCode).
		Where(dao.Domain.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Domain.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("域名已绑定")
	}
	return nil
}

func (s *sDomain) ensureAccessible(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("域名不存在或已删除")
	}
	return shared.EnsureTenantScopedRowAccessible(ctx, dao.Domain.Ctx(ctx), id, dao.Domain.Columns().Id, dao.Domain.Columns().TenantId, dao.Domain.Columns().MerchantId, "域名")
}

func (s *sDomain) loadRow(ctx context.Context, id snowflake.JsonInt64) (*domainRow, error) {
	if id <= 0 {
		return nil, gerror.New("域名不存在或已删除")
	}
	row := &domainRow{}
	if err := dao.Domain.Ctx(ctx).
		Where(dao.Domain.Columns().Id, id).
		Where(dao.Domain.Columns().DeletedAt, nil).
		Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return row, nil
}

func (s *sDomain) fillNames(ctx context.Context, list []*model.DomainListOutput) {
	if len(list) == 0 {
		return
	}
	tenantIDs := make([]int64, 0, len(list))
	merchantIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.TenantID > 0 {
			tenantIDs = append(tenantIDs, int64(item.TenantID))
		}
		if item.MerchantID > 0 {
			merchantIDs = append(merchantIDs, int64(item.MerchantID))
		}
	}
	tenantMap := s.loadTenantNameMap(ctx, tenantIDs)
	merchantMap := s.loadMerchantNameMap(ctx, merchantIDs)
	for _, item := range list {
		item.TenantName = tenantMap[int64(item.TenantID)]
		item.MerchantName = merchantMap[int64(item.MerchantID)]
	}
}

func (s *sDomain) loadTenantNameMap(ctx context.Context, ids []int64) map[int64]string {
	ids = compactPositiveInt64s(ids)
	if len(ids) == 0 {
		return nil
	}
	var rows []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := dao.Tenant.Ctx(ctx).
		Fields(dao.Tenant.Columns().Id, dao.Tenant.Columns().Name).
		WhereIn(dao.Tenant.Columns().Id, ids).
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

func (s *sDomain) loadMerchantNameMap(ctx context.Context, ids []int64) map[int64]string {
	ids = compactPositiveInt64s(ids)
	if len(ids) == 0 {
		return nil
	}
	var rows []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := dao.Merchant.Ctx(ctx).
		Fields(dao.Merchant.Columns().Id, dao.Merchant.Columns().Name).
		WhereIn(dao.Merchant.Columns().Id, ids).
		Where(dao.Merchant.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
		return nil
	}
	out := make(map[int64]string, len(rows))
	for _, row := range rows {
		out[row.Id] = row.Name
	}
	return out
}

func compactPositiveInt64s(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	out := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
