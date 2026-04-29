
package warehouse_trade

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
	service.RegisterWarehouseTrade(New())
}

func New() *sWarehouseTrade {
	return &sWarehouseTrade{}
}

type sWarehouseTrade struct{}

func normalizeWarehouseTradeIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建仓库交易记录
func (s *sWarehouseTrade) Create(ctx context.Context, in *model.WarehouseTradeCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberWarehouseTrade.Ctx(ctx).Data(do.MemberWarehouseTrade{
		Id:        id,
		TradeNo: in.TradeNo,
		GoodsId: in.GoodsID,
		ListingId: in.ListingID,
		SellerId: in.SellerID,
		BuyerId: in.BuyerID,
		TradePrice: in.TradePrice,
		PlatformFee: in.PlatformFee,
		SellerIncome: in.SellerIncome,
		TradeStatus: in.TradeStatus,
		ConfirmedAt: in.ConfirmedAt,
		Remark: in.Remark,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新仓库交易记录
func (s *sWarehouseTrade) Update(ctx context.Context, in *model.WarehouseTradeUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberWarehouseTrade{
		TradeNo: in.TradeNo,
		GoodsId: in.GoodsID,
		ListingId: in.ListingID,
		SellerId: in.SellerID,
		BuyerId: in.BuyerID,
		TradePrice: in.TradePrice,
		PlatformFee: in.PlatformFee,
		SellerIncome: in.SellerIncome,
		TradeStatus: in.TradeStatus,
		ConfirmedAt: in.ConfirmedAt,
		Remark: in.Remark,
		Status: in.Status,
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.MemberWarehouseTrade.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseTrade.Columns().Id, in.ID).
			Where(dao.MemberWarehouseTrade.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("仓库交易记录不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx), in.ID, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().TenantId, dao.MemberWarehouseTrade.Columns().MerchantId, "仓库交易记录"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx), in.ID, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().CreatedBy, dao.MemberWarehouseTrade.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.MemberWarehouseTrade.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseTrade.Columns().Id, in.ID).
			Where(dao.MemberWarehouseTrade.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除仓库交易记录
func (s *sWarehouseTrade) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), id, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().TenantId, dao.MemberWarehouseTrade.Columns().MerchantId, "仓库交易记录"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), id, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().CreatedBy, dao.MemberWarehouseTrade.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberWarehouseTrade.Ctx(ctx).Where(dao.MemberWarehouseTrade.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除仓库交易记录
func (s *sWarehouseTrade) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeWarehouseTradeIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), normalizedIDs, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().TenantId, dao.MemberWarehouseTrade.Columns().MerchantId, "仓库交易记录"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), normalizedIDs, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().CreatedBy, dao.MemberWarehouseTrade.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberWarehouseTrade.Ctx(ctx).WhereIn(dao.MemberWarehouseTrade.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取仓库交易记录详情
func (s *sWarehouseTrade) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WarehouseTradeDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), id, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().TenantId, dao.MemberWarehouseTrade.Columns().MerchantId, "仓库交易记录"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), id, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().CreatedBy, dao.MemberWarehouseTrade.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.WarehouseTradeDetailOutput{}
	err = dao.MemberWarehouseTrade.Ctx(ctx).Where(dao.MemberWarehouseTrade.Columns().Id, id).Where(dao.MemberWarehouseTrade.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("仓库交易记录不存在或已删除")
	}
	// 查询仓库商品关联显示
	if out.GoodsID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_warehouse_goods").Where("id", out.GoodsID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("title")
		if err == nil {
			out.WarehouseGoodsTitle = val.String()
		}
	}
	// 查询挂卖记录关联显示
	if out.ListingID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_warehouse_listing").Where("id", out.ListingID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("id")
		if err == nil {
			out.WarehouseListingID = val.String()
		}
	}
	// 查询卖家关联显示
	if out.SellerID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.SellerID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("nickname")
		if err == nil {
			out.UserNickname = val.String()
		}
	}
	// 查询买家关联显示
	if out.BuyerID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.BuyerID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("nickname")
		if err == nil {
			out.BuyerNickname = val.String()
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
func (s *sWarehouseTrade) applyListFilter(ctx context.Context, in *model.WarehouseTradeListInput) *gdb.Model {
	m := dao.MemberWarehouseTrade.Ctx(ctx).Where(dao.MemberWarehouseTrade.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberWarehouseTrade.Columns().TenantId, dao.MemberWarehouseTrade.Columns().MerchantId)
	if in.TradeNo != "" {
		m = m.Where(dao.MemberWarehouseTrade.Columns().TradeNo, in.TradeNo)
	}
	if in.GoodsID != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().GoodsId, *in.GoodsID)
	}
	if in.ListingID != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().ListingId, *in.ListingID)
	}
	if in.SellerID != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().SellerId, *in.SellerID)
	}
	if in.BuyerID != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().BuyerId, *in.BuyerID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().MerchantId, *in.MerchantID)
	}
	if in.TradeStatus != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().TradeStatus, *in.TradeStatus)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberWarehouseTrade.Columns().Status, *in.Status)
	}
	if in.ConfirmedAtStart != "" {
		m = m.WhereGTE(dao.MemberWarehouseTrade.Columns().ConfirmedAt, in.ConfirmedAtStart)
	}
	if in.ConfirmedAtEnd != "" {
		m = m.WhereLTE(dao.MemberWarehouseTrade.Columns().ConfirmedAt, in.ConfirmedAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberWarehouseTrade.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberWarehouseTrade.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberWarehouseTrade.Columns().CreatedBy, dao.MemberWarehouseTrade.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sWarehouseTrade) fillRefFields(ctx context.Context, list []*model.WarehouseTradeListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.GoodsID != 0 {
				idSet[int64(item.GoodsID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_warehouse_goods").
				Fields("id", "title")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["title"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.GoodsID)]; ok {
						item.WarehouseGoodsTitle = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.ListingID != 0 {
				idSet[int64(item.ListingID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_warehouse_listing").
				Fields("id", "id")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["id"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.ListingID)]; ok {
						item.WarehouseListingID = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.SellerID != 0 {
				idSet[int64(item.SellerID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_user").
				Fields("id", "nickname")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["nickname"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.SellerID)]; ok {
						item.UserNickname = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.BuyerID != 0 {
				idSet[int64(item.BuyerID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_user").
				Fields("id", "nickname")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["nickname"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.BuyerID)]; ok {
						item.BuyerNickname = val
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

// List 获取仓库交易记录列表
func (s *sWarehouseTrade) List(ctx context.Context, in *model.WarehouseTradeListInput) (list []*model.WarehouseTradeListOutput, total int, err error) {
	if in == nil {
		in = &model.WarehouseTradeListInput{}
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
func (s *sWarehouseTrade) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberWarehouseTrade.Columns().Id:        true,
		dao.MemberWarehouseTrade.Columns().CreatedAt: true,
		dao.MemberWarehouseTrade.Columns().Status:    true,
		dao.MemberWarehouseTrade.Columns().TradeNo: true,
		dao.MemberWarehouseTrade.Columns().TradePrice: true,
		dao.MemberWarehouseTrade.Columns().PlatformFee: true,
		dao.MemberWarehouseTrade.Columns().SellerIncome: true,
		dao.MemberWarehouseTrade.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sWarehouseTrade) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.MemberWarehouseTrade.Columns().Id)
}

// Export 导出仓库交易记录（不分页）
func (s *sWarehouseTrade) Export(ctx context.Context, in *model.WarehouseTradeListInput) (list []*model.WarehouseTradeListOutput, err error) {
	if in == nil {
		in = &model.WarehouseTradeListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑仓库交易记录
func (s *sWarehouseTrade) BatchUpdate(ctx context.Context, in *model.WarehouseTradeBatchUpdateInput) error {
	data := do.MemberWarehouseTrade{}
	hasChange := false
	if in.TradeStatus != nil {
		data.TradeStatus = *in.TradeStatus
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeWarehouseTradeIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), normalizedIDs, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().TenantId, dao.MemberWarehouseTrade.Columns().MerchantId, "仓库交易记录"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberWarehouseTrade.Ctx(ctx), normalizedIDs, dao.MemberWarehouseTrade.Columns().Id, dao.MemberWarehouseTrade.Columns().CreatedBy, dao.MemberWarehouseTrade.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberWarehouseTrade.Ctx(ctx).WhereIn(dao.MemberWarehouseTrade.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入仓库交易记录
func (s *sWarehouseTrade) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.MemberWarehouseTrade{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.TradeNo = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.GoodsId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.ListingId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.SellerId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.BuyerId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.TradePrice = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.PlatformFee = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.SellerIncome = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.TradeStatus = v
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
		if fkVal, ok := data.GoodsId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_warehouse_goods").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.ListingId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_warehouse_listing").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.SellerId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.BuyerId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.MemberWarehouseTrade.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

